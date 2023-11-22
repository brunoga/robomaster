package chassis

import (
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
)

type Chassis struct {
	*internal.BaseModule
}

func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Chassis, error) {
	return &Chassis{
		// Apparently chassis is the main controller.
		// TODO(bga): Verify this.
		BaseModule: internal.NewBaseModule(ub, l, "Chassis",
			key.KeyMainControllerConnection, nil),
	}, nil
}
