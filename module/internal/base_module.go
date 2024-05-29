package internal

import (
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/support"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
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

	rl *support.ResultListener
	cb result.Callback
}

var _ module.Module = (*BaseModule)(nil)

// NewBaseModule creates a new BaseModule instance with the given name and that
// will listen for results with the given key. The given callback, if not nil,
// will be called whenever a new result is received.
func NewBaseModule(ub unitybridge.UnityBridge, l *logger.Logger,
	name string, k *key.Key, cb result.Callback,
	deps ...module.Module) *BaseModule {

	var rl *support.ResultListener
	if k != nil {
		// We have a key so we need a result listener.
		rl = support.NewResultListener(ub, l, k, cb)
	}

	return &BaseModule{
		ub:   ub,
		l:    l,
		name: name,
		deps: deps,
		rl:   rl,
		cb:   cb,
	}
}

// Start starts the module by starting the connection result listener.
func (g *BaseModule) Start() error {
	if g.rl == nil {
		if g.cb != nil {
			defer func() {
				g.cb(result.New(nil, 0, 0, "", &value.Bool{Value: true}))
			}()
		}

		return nil
	}

	return g.rl.Start()
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
	if g.isConnected() {
		return true
	}

	start := time.Now()

	if g.rl.WaitForAnyResult(timeout) == nil {
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
	if g.rl == nil {
		if g.cb != nil {
			defer func() {
				g.cb(result.New(nil, 0, 0, "", &value.Bool{Value: false}))
			}()
		}

		return nil
	}

	return g.rl.Stop()
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
	if g.rl != nil {
		r := g.rl.Result()

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
