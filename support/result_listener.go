package support

import (
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

// ResultListener is a helper class to listen for event results from the
// Unity Bridge. It allows callers to wait for new results, to get the
// the last result obtained and to register a callback to be called when
// a new result is available. It is thread safe (and lock free).
type ResultListener struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger
	k  *key.Key
	cb result.Callback

	t token.Token

	r       atomic.Pointer[result.Result]
	c       atomic.Pointer[chan struct{}]
	started atomic.Bool
}

// NewResultListener creates a new ResultListener instance.
func NewResultListener(ub unitybridge.UnityBridge, l *logger.Logger,
	k *key.Key, cb result.Callback) *ResultListener {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("result_listener").With(
		slog.String("key", k.String()))

	lr := &ResultListener{
		ub: ub,
		l:  l,
		k:  k,
		cb: cb,
	}

	c := make(chan struct{})
	lr.c.Store(&c)

	return lr
}

// Start starts the listener. If cb is non nil, it will be called when a new
// result is available.
func (ls *ResultListener) Start() error {
	if ls.started.Load() {
		return fmt.Errorf("listener already started")
	}

	c := make(chan struct{})
	ls.c.Store(&c)

	ls.r.Store(nil)

	var err error

	ls.t, err = ls.ub.AddKeyListener(ls.k, func(r *result.Result) {
		ls.r.Store(r)

		ls.notifyWaiters()

		if ls.cb != nil {
			go ls.cb(r)
		}
	}, true)

	ls.started.Store(true)

	return err
}

// WaitForNewResult blocks until a new result is available, a timeout happens
// or the listener is stopped. IF result is nil, no result was available (for
// example, if the listener is closed ). If result is non nil, Callers should
// inspect the result error code and description to check if the result is
// valid.
func (ls *ResultListener) WaitForNewResult(timeout time.Duration) *result.Result {
	select {
	case <-*ls.c.Load():
		return ls.r.Load()
	case <-time.After(timeout):
		return nil
	}
}

// Result returns the current result.
func (ls *ResultListener) Result() *result.Result {
	return ls.r.Load()
}

// Stop stops the listener.
func (ls *ResultListener) Stop() error {
	if !ls.started.Load() {
		return fmt.Errorf("listener not started")
	}

	err := ls.ub.RemoveKeyListener(ls.k, ls.t)
	if err != nil {
		return err
	}

	ls.r.Store(nil)
	ls.started.Store(false)

	ls.notifyWaiters()

	return nil
}

// notifyWaiters closes the current channel and creates a new one.
// This is done in a way that is safe for concurrent access.
func (ls *ResultListener) notifyWaiters() {
	oldCPtr := ls.c.Load()
	if oldCPtr == nil {
		ls.l.Error("Old channel pointer is unexpectedly nil")
		// This should never happen but lets be safe.
		return
	}

	newC := make(chan struct{})
	if ls.c.CompareAndSwap(oldCPtr, &newC) {
		// We managed to swap the old channel pointer with the new one so we
		// can safely close the new one now (which will notify listeners). Note
		// that if CompareAndSwap() returned false, if means another thread
		// closed managed to swap (and close) the channels between our load
		// and CompareAndSwap() calls. In this case, we can safely assume our
		// old channel is already closed and do not need to retry.
		ls.l.Debug("Notifying waiters")
		close(*oldCPtr)
	}
}
