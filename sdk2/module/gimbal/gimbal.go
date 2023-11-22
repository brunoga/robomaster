package gimbal

import (
	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
)

type Gimbal struct {
	*internal.BaseModule
}

var _ module.Module = (*Gimbal)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Gimbal, error) {
	return &Gimbal{
		BaseModule: internal.NewBaseModule(ub, l, "Gimbal",
			key.KeyGimbalConnection, nil),
	}, nil
}
