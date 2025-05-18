package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/listener"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

// BaseModule is a base implementation of the module.Module interface. It takes
// care of handling the connection status of the module and provides default
// implementations for all interface methods.
//
// Module implementations can simply embed this or they can provide custom logic
// for each method (just be sure to always call the base implementation).
type BaseModule struct {
	ub   unitybridge.UnityBridge
	l    *logger.Logger
	name string
	deps []module.Module

	mRl sync.Mutex
	rls []*listener.Listener
	rlm map[*key.Key]struct{}

	m       sync.Mutex
	started bool
}

var _ module.Module = (*BaseModule)(nil)

// NewBaseModule creates a new BaseModule instance with the given name and that
// will listen for results with the given key. The given callback, if not nil,
// will be called whenever a new result is received.
func NewBaseModule(ub unitybridge.UnityBridge, l *logger.Logger,
	name string, k *key.Key, cb result.Callback,
	deps ...module.Module) (*BaseModule, error) {

	if k != nil && cb == nil {
		return nil, fmt.Errorf("if a key is provided, a callback must be provided")
	}

	m := &BaseModule{
		ub:   ub,
		l:    l,
		name: name,
		deps: deps,
		rls:  make([]*listener.Listener, 1),
		rlm:  make(map[*key.Key]struct{}, 1),
	}

	if k != nil {
		// We have a key so we need a result listener.
		err := m.AddKeyListener(k, cb, false)
		if err != nil {
			return nil, fmt.Errorf("failed to add key listener: %w", err)
		}
	} else {
		// No connection key.
		m.rls[0] = nil
	}

	return m, nil
}

func (g *BaseModule) AddKeyListener(k *key.Key, cb func(*result.Result), async bool) error {
	if k == nil || cb == nil {
		return fmt.Errorf("key and callback must not be nil")
	}

	g.m.Lock()
	defer g.m.Unlock()

	if g.started {
		return fmt.Errorf("module already started")
	}

	g.mRl.Lock()
	defer g.mRl.Unlock()

	if _, ok := g.rlm[k]; ok {
		return fmt.Errorf("key already registered")
	}

	g.rls = append(g.rls, listener.New(g.ub, g.l, k, cb))
	g.rlm[k] = struct{}{}

	return nil
}

// Start starts the module by starting the connection result listener.
func (g *BaseModule) Start() error {
	g.m.Lock()
	defer g.m.Unlock()

	if g.started {
		return fmt.Errorf("module already started")
	}

	g.started = true

	for i, rl := range g.rls {
		if rl != nil {
			err := rl.Start()
			if err != nil {
				// Stop all started listeners.
				for j := i - 1; j >= 0; j-- {
					if g.rls[j] != nil {
						g.rls[j].Stop()
					}
				}

				return err
			}
		}
	}

	return nil
}

// Connected returns true if the module is connected, false otherwise.
func (g *BaseModule) Connected() bool {
	return g.isConnected()
}

// WaitForConnection returns the current connection status, if one is
// available or waits for a new one for the given timeout period. It
// returns true if the module is connected, false otherwise (including
// if the timeout period is reached or an error happens).
func (g *BaseModule) WaitForConnection(timeout time.Duration) bool {
	g.m.Lock()
	if !g.started {
		g.m.Unlock()
		return false
	}
	defer g.m.Unlock()

	if g.isConnected() {
		return true
	}

	start := time.Now()

	if g.rls[0].WaitForAnyResult(timeout) == nil {
		return false
	}

	for _, dep := range g.deps {
		timeout -= time.Since(start)
		if !dep.WaitForConnection(timeout) {
			return false
		}
	}

	return g.isConnected()
}

// Stop stops the module by stopping the connection result listener.
func (g *BaseModule) Stop() error {
	g.m.Lock()
	defer g.m.Unlock()

	if !g.started {
		return fmt.Errorf("module not started")
	}

	g.started = false

	for i := len(g.rls) - 1; i >= 0; i-- {
		if g.rls[i] != nil {
			err := g.rls[i].Stop()
			if err != nil {
				g.l.Warn("Failed to stop result listener.", "error", err)
			}
		}
	}

	g.rls = make([]*listener.Listener, 1)
	g.rlm = make(map[*key.Key]struct{}, 1)

	return nil
}

// String returns the module name.
func (g *BaseModule) String() string {
	return g.name
}

// UB returns the UnityBridge instance used by the module.
func (g *BaseModule) UB() unitybridge.UnityBridge {
	return g.ub
}

// Logger returns the Logger instance used by the module.
func (g *BaseModule) Logger() *logger.Logger {
	return g.l
}

func (g *BaseModule) isConnected() bool {
	if !g.started {
		return false
	}

	if g.rls[0] != nil {
		r := g.rls[0].Result()

		if !r.Succeeded() {
			return false
		}

		connected, ok := r.Value().(*value.Bool)
		if !ok || !connected.Value {
			return false
		}
	}

	// Check if dependencies are connected. Notice that this does not wait for
	// connection currently so if any dependencies are not connected, this will
	// return false.
	//
	// TODO(bga): Maybe actually wait for connection here?
	depsConnected := true
	for _, dep := range g.deps {
		if !dep.Connected() {
			depsConnected = false
			break
		}
	}

	return depsConnected
}
