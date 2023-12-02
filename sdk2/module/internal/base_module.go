package internal

import (
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
	"github.com/brunoga/unitybridge/unity/result/value"
)

// BaseModule is a base implementation of the module.Module interface. It takes
// care of handling the connection status of the module and provides default
// implementations for all interface methods.
//
// Module implemnetations can simply embed this or it can provide custom logic
// for each method (just be sure to always call the base implementation).
type BaseModule struct {
	ub   unitybridge.UnityBridge
	l    *logger.Logger
	name string

	rl *support.ResultListener
}

var _ module.Module = (*BaseModule)(nil)

// NewBaseModule creates a new BaseModule instance with the given name and that
// will listen for results with the given key. The given callback, if not nil,
// will be called whenever a new result is received.
func NewBaseModule(ub unitybridge.UnityBridge, l *logger.Logger,
	name string, k *key.Key, cb result.Callback) *BaseModule {
	return &BaseModule{
		ub:   ub,
		l:    l,
		name: name,
		rl:   support.NewResultListener(ub, l, k, cb),
	}
}

// Start starte the module by starting the connection result listener.
func (g *BaseModule) Start() error {
	return g.rl.Start()
}

// Connected returns true if the module is connected, false otherwise.
func (g *BaseModule) Connected() bool {
	return g.isConnected(g.rl.Result())
}

// WaitForConnection returns the current connection status, if one is
// available or waits for a new one for the given timeout period. It
// returns true if the module is connected, false otherwise (including
// if the timeout period is reached or an error happens).
func (g *BaseModule) WaitForConnection(timeout time.Duration) bool {
	if g.isConnected(g.rl.Result()) {
		return true
	}

	r := g.rl.WaitForNewResult(timeout)

	return g.isConnected(r)
}

// Stop stops the module by stopping the connection result listener.
func (g *BaseModule) Stop() error {
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

func (g *BaseModule) isConnected(r *result.Result) bool {
	if r == nil || r.ErrorCode() != 0 {
		return false
	}

	connected, ok := r.Value().(*value.Bool)

	return ok && connected.Value
}
