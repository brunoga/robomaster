package internal

import (
	"fmt"
	"sync/atomic"

	"github.com/brunoga/robomaster/module2"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/support/token"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

type BaseModule struct {
	ub   unitybridge.UnityBridge
	l    *logger.Logger
	name string
	k    *key.Key

	started atomic.Bool

	connected   atomic.Bool
	connectedCh atomic.Pointer[chan struct{}]
}

var _ module2.Module = &BaseModule{}

func NewBaseModule(ub unitybridge.UnityBridge, l *logger.Logger, name string,
	k *key.Key) *BaseModule {

	l = l.With("connection_key", k)

	return &BaseModule{
		ub:   ub,
		l:    l,
		name: name,
		k:    k,
	}
}

// module.Module2 interface implementation.

func (bm *BaseModule) String() string {
	return bm.name
}

func (bm *BaseModule) Connected() bool {
	return bm.connected.Load()
}

func (bm *BaseModule) WaitForConnectionStatus(connected bool) {
	for !bm.connected.CompareAndSwap(connected, connected) {
		<-*bm.connectedCh.Load()
	}
}

func (bm *BaseModule) AddConnectionCallback(callback result.Callback) (token.Token, error) {
	if !bm.started.CompareAndSwap(true, true) {
		return 0, fmt.Errorf("not started")
	}

	if bm.k == nil {
		return 0, fmt.Errorf("module has no associated key")
	}

	return bm.ub.AddKeyListener(bm.k, callback, true)
}

func (bm *BaseModule) RemoveConnectionCallback(token token.Token) error {
	if !bm.started.CompareAndSwap(true, true) {
		return fmt.Errorf("not started")
	}

	if bm.k == nil {
		return fmt.Errorf("module has no associated key")
	}

	return bm.ub.RemoveKeyListener(bm.k, token)
}

func (bm *BaseModule) Start() error {
	if !bm.started.CompareAndSwap(false, true) {
		return fmt.Errorf("already started")
	}

	ch := make(chan struct{})
	bm.connectedCh.Store(&ch)

	if bm.k != nil {
		_, err := bm.AddConnectionCallback(bm.handleConnected)
		if err != nil {
			return err
		}
	} else {
		// No key so generate a connected event.
		bm.handleConnected(result.New(nil, 0, 0, "", &value.Bool{Value: true}))
	}

	return nil
}

func (bm *BaseModule) Stop() error {
	if !bm.started.CompareAndSwap(true, false) {
		return fmt.Errorf("not started")
	}

	if bm.k != nil {
		// Remove all key listeners.
		err := bm.ub.RemoveKeyListener(bm.k, 0)
		if err != nil {
			return err
		}
	}

	// When stopping, we always generate a disconnected event so any waiters can
	// be unblocked as at this point the connected/disconnected status will not
	// be handled anymore.
	bm.handleConnected(result.New(bm.k, 0, 0, "", &value.Bool{Value: false}))

	return nil
}

func (bm *BaseModule) UB() unitybridge.UnityBridge {
	return bm.ub
}

func (bm *BaseModule) L() *logger.Logger {
	return bm.l
}

// handleConnected is called whenever the connection status changes. Its
// purpose is to validate and update the module connection status accordingly
// and also log any changes.
func (bm *BaseModule) handleConnected(r *result.Result) {
	if !r.Succeeded() {
		bm.l.Error("Failed to get connection status", "error", r.ErrorDesc())
		return
	}

	value, ok := r.Value().(*value.Bool)
	if !ok {
		bm.l.Error("Invalid value type", "value", r.Value())
		return
	}

	bm.connected.Store(value.Value)

	newCh := make(chan struct{})
	oldChPtr := bm.connectedCh.Load()
	if bm.connectedCh.CompareAndSwap(oldChPtr, &newCh) {
		// oldCh was still the expected so now connectedCh has been updated and
		// we can close the old channel.
		close(*oldChPtr)
	}

	if value.Value {
		bm.l.Info("Connected")
	} else {
		bm.l.Info("Disconnected")
	}
}
