package gamepad

import (
	"log/slog"
	"sync/atomic"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
	"github.com/brunoga/robomaster/unitybridge/support/token"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
)

type GamePad struct {
	*internal.BaseModule

	c1Token   token.Token
	c2Token   token.Token
	fireToken token.Token
	fnToken   token.Token

	c1Status   atomic.Bool
	c2Status   atomic.Bool
	fireStatus atomic.Bool
	fnStatus   atomic.Bool
}

var _ module.Module = (*GamePad)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection) (*GamePad, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("gamepad_module")

	// TODO(bga): On connection we might need to activate the controller.

	return &GamePad{
		BaseModule: internal.NewBaseModule(ub, l, "GamePad",
			key.KeyRobomasterGamePadConnection, nil, cm),
	}, nil
}

func (m *GamePad) Start() error {
	var err error

	m.c1Token, err = m.UB().AddKeyListener(key.KeyRobomasterGamePadC1,
		m.onButton, false)
	if err != nil {
		return err
	}
	m.c2Token, err = m.UB().AddKeyListener(key.KeyRobomasterGamePadC2,
		m.onButton, false)
	if err != nil {
		return err
	}
	m.fireToken, err = m.UB().AddKeyListener(key.KeyRobomasterGamePadFire,
		m.onButton, false)
	if err != nil {
		return err
	}
	m.fnToken, err = m.UB().AddKeyListener(key.KeyRobomasterGamePadFn,
		m.onButton, false)
	if err != nil {
		return err
	}

	return m.BaseModule.Start()
}

func (m *GamePad) Stop() error {
	var err error

	err = m.UB().RemoveKeyListener(key.KeyRobomasterGamePadC1, m.c1Token)
	if err != nil {
		return err
	}

	err = m.UB().RemoveKeyListener(key.KeyRobomasterGamePadC2, m.c2Token)
	if err != nil {
		return err
	}

	err = m.UB().RemoveKeyListener(key.KeyRobomasterGamePadFire, m.fireToken)
	if err != nil {
		return err
	}

	err = m.UB().RemoveKeyListener(key.KeyRobomasterGamePadFn, m.fnToken)
	if err != nil {
		return err
	}

	return m.BaseModule.Stop()
}

func (m *GamePad) C1Pressed() bool {
	return m.c1Status.Load()
}

func (m *GamePad) C2Pressed() bool {
	return m.c2Status.Load()
}

func (m *GamePad) FirePressed() bool {
	return m.fireStatus.Load()
}

func (m *GamePad) FnPressed() bool {
	return m.fnStatus.Load()
}

func (m *GamePad) onButton(r *result.Result) {
	if r == nil || !r.Succeeded() {
		return
	}

	v, ok := r.Value().(bool)
	if !ok {
		return
	}

	switch r.Key() {
	case key.KeyRobomasterGamePadC1:
		m.c1Status.Store(v)
	case key.KeyRobomasterGamePadC2:
		m.c2Status.Store(v)
	case key.KeyRobomasterGamePadFire:
		m.fireStatus.Store(v)
	case key.KeyRobomasterGamePadFn:
		m.fnStatus.Store(v)
	default:
		m.Logger().Error("Received unexpected button key", "key", r.Key())
	}
}
