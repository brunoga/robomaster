package listener

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/support/token"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/timedsignalwaiter"
)

// Listener is a helper class to listen for event results from the
// Unity Bridge. It allows callers to wait for new results, to get the last
// result obtained and to register a callback to be called when a new result
// is available. It is thread safe.
type Listener struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger
	k  *key.Key
	cb result.Callback

	t token.Token

	b *timedsignalwaiter.TimedSignalWaiter

	m       sync.Mutex
	r       *result.Result
	started bool
}

// New creates a new Listener instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	k *key.Key, cb result.Callback) *Listener {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("result_listener").With(
		slog.String("key", k.String()))

	return &Listener{
		ub: ub,
		l:  l,
		k:  k,
		cb: cb,
		b:  timedsignalwaiter.New(k.String()),
	}
}

// Start starts the listener. If cb is non nil, it will be called when a new
// result is available.
func (ls *Listener) Start() error {
	ls.m.Lock()
	defer ls.m.Unlock()

	if ls.started {
		return fmt.Errorf("listener already started")
	}

	ls.r = nil

	var err error

	ls.t, err = ls.ub.AddKeyListener(ls.k, func(r *result.Result) {
		ls.l.Debug("Received result.", "key", ls.k, "result", r)

		// Fierst we synchronously execute any associated callback so that any
		// required initialization might be completed before we notify any
		// waiters.
		if ls.cb != nil && r.Succeeded() {
			ls.l.Debug("Calling ResultListener callback.", "key", ls.k, "result", r)
			ls.cb(r)
		} else {
			ls.l.Debug("Not calling ResultListener callback.", "key", ls.k, "result", r, "nil_callback", ls.cb == nil)
		}

		// Now we are going to change our state, so we need to lock.
		ls.m.Lock()

		// Update our result cache.
		ls.r = r

		// And now we can notify waiters.
		ls.l.Debug("Notifying waiters.")
		ls.notifyWaitersLocked()

		ls.m.Unlock()
	}, true)

	ls.started = true

	return err
}

// WaitForNewResult blocks until a new result is available, a timeout happens
// or the listener is stopped. If result is nil, no result was available (for
// example, if the listener is closed). If result is non nil, Callers should
// inspect the result error code and description to check if the result is
// valid.
func (ls *Listener) WaitForNewResult(timeout time.Duration) *result.Result {
	if ls.b.Wait(timeout) {
		ls.m.Lock()
		defer ls.m.Unlock()
		return ls.r
	}

	return nil
}

// WaitForAnyResult returns any existing result immediatelly or blocks until a
// result is available, a timeout happens or the listener is stopped. If result
// is nil, no result was available (for example, if the listener is closed). If
// result is non nil, Callers should inspect the result error code and
// description to check if the result is valid.
func (ls *Listener) WaitForAnyResult(timeout time.Duration) *result.Result {
	// Make sure we get a correct snapshot of the current channel and result
	// state by obtaining them inside a lock. This guarantees that we either
	// have a result or that, if we do not, we are going to be listening on a
	// channel that is guaranteed to be the one existing when the value
	// was nil so either it is closed now and we do have a non-nil value or
	// it will be closed after we start waiting on it (and we will get a result
	// or a timeout.
	ls.l.Debug("Waiting for any result.", "key", ls.k)
	ls.m.Lock()
	if ls.r != nil {
		ls.l.Debug("Existing result is not nil.", "key", ls.k)
		ls.m.Unlock()
		return ls.r
	}

	ls.m.Unlock()

	ls.l.Debug("Existing result is nil.", "key", ls.k)

	ls.l.Debug("Waiting for new result.", "key", ls.k)
	return ls.WaitForNewResult(timeout)
}

// Result returns the current result.
func (ls *Listener) Result() *result.Result {
	ls.m.Lock()
	defer ls.m.Unlock()

	return ls.r
}

// AddCallback adds a callback to be called when a new result is available for
// this listener key. The returned token can be used to remove the callback
// later.
func (ls *Listener) AddCallback(cb result.Callback) (token.Token, error) {
	return ls.ub.AddKeyListener(ls.k, cb, true)
}

// RemoveCallback removes the callback associated with the likstener key and the
// given token.
func (ls *Listener) RemoveCallback(t token.Token) error {
	return ls.ub.RemoveKeyListener(ls.k, t)
}

// Stop stops the listener.
func (ls *Listener) Stop() error {
	ls.m.Lock()
	defer ls.m.Unlock()

	if !ls.started {
		return fmt.Errorf("listener not started")
	}

	err := ls.ub.RemoveKeyListener(ls.k, ls.t)
	if err != nil {
		return err
	}

	ls.r = nil
	ls.started = false

	ls.notifyWaitersLocked()

	return nil
}

// notifyWaitersLocked closes the current channel and creates a new one.
// The channel mutex must be locked when this is called.
func (ls *Listener) notifyWaitersLocked() {
	ls.l.Debug("Notifying waiters.", "key", ls.k)
	ls.b.Signal()
	ls.l.Debug("Notified waiters.", "key", ls.k)
}
