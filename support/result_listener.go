package support

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

type ResultListener struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger
	k  *key.Key

	t token.Token

	m sync.RWMutex
	r result.Result
	c chan struct{}
}

func NewResultListener(ub unitybridge.UnityBridge, l *logger.Logger,
	k *key.Key) *ResultListener {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	return &ResultListener{
		ub: ub,
		l:  l,
		k:  k,
		c:  nil,
	}
}

func (ls *ResultListener) Start(cb result.Callback) error {
	ls.m.Lock()
	defer ls.m.Unlock()

	if ls.c != nil {
		return fmt.Errorf("listener already started")
	}

	ls.c = make(chan struct{})
	ls.r = result.Result{}

	var err error

	ls.t, err = ls.ub.AddKeyListener(ls.k, func(r *result.Result) {
		ls.m.Lock()
		ls.r = *r
		close(ls.c)
		ls.c = make(chan struct{})
		if cb != nil {
			go cb(r)
		}
		ls.m.Unlock()
	}, true)

	return err
}

// WaitForNewResult blocks until a new result is available, a timeout happens
// or the listener is stopped. Callers should inspect the result error code
// and description to check if the result is valid.
func (ls *ResultListener) WaitForNewResult(timeout time.Duration) *result.Result {
	ls.m.RLock()
	if ls.c == nil {
		ls.m.RUnlock()
		r := &result.Result{}
		r.SetErrorCode(-1)
		r.SetErrorDesc("listener not started")
		return r
	}
	ls.m.RUnlock()

	if timeout != 0 {
		select {
		case <-ls.c:
			ls.m.RLock()
			defer ls.m.RUnlock()
			return &ls.r
		case <-time.After(timeout):
			r := &result.Result{}
			r.SetErrorCode(-1)
			r.SetErrorDesc("listener timeout")
			return r
		}
	}

	<-ls.c

	ls.m.RLock()
	defer ls.m.RUnlock()

	return &ls.r
}

func (ls *ResultListener) Result() *result.Result {
	ls.m.RLock()
	defer ls.m.RUnlock()

	return &ls.r
}

func (ls *ResultListener) Stop() error {
	ls.m.Lock()
	defer ls.m.Unlock()

	if ls.c == nil {
		return fmt.Errorf("listener not started")
	}

	err := ls.ub.RemoveKeyListener(ls.k, ls.t)
	if err != nil {
		return err
	}

	// Set an error in the current result and make sure waiters
	// will be notified.
	ls.r.SetErrorCode(-1)
	ls.r.SetErrorDesc("listener stopped")
	ls.r.SetValue(nil)

	close(ls.c)

	ls.c = nil

	return nil
}
