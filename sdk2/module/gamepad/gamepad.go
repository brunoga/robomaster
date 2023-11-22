package gamepad

import (
	"log/slog"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
)

type GamePad struct {
	*internal.BaseModule
}

var _ module.Module = (*GamePad)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger) (*GamePad, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("gamepad_module")

	return &GamePad{
		internal.NewBaseModule(ub, l, "GamePad",
			key.KeyRobomasterGamePadConnection, nil),
	}, nil
}
