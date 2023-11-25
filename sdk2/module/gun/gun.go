package gun

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
)

// Gun is the module that controls turret firing.
type Gun struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger
}

var _ module.Module = (*Gun)(nil)

// New creates a new Gun instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Gun, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("gun_module")

	return &Gun{
		ub: ub,
		l:  l,
	}, nil
}

// Start starts the Gun module.
func (g *Gun) Start() error {
	return nil
}

// Connected returns whether the Gun module is connected. It always returns
// true as there is no connection to handle.
func (g *Gun) Connected() bool {
	return true
}

// WaitForConnection waits for the Gun module to connect. It always returns
// immediately as there is no connection to handle.
func (g *Gun) WaitForConnection(timeout time.Duration) bool {
	return true
}

// Fire fires the gun.
//
// TODO(bga): The id parameter seems to always be 0 or 1.
func (g *Gun) Fire(typ Type, id uint8) error {
	if !typ.Valid() {
		return fmt.Errorf("invalid gun type %d", typ)
	}

	var k *key.Key
	switch typ {
	case TypeBead:
		k = key.KeyRobomasterWaterGunWaterGunFire
	case TypeInfrared:
		k = key.KeyRobomasterInfraredGunInfraredGunFire
	}

	return g.ub.DirectSendKeyValue(k, uint64(id))
}

// Stop stops the Gun module.
func (g *Gun) Stop() error {
	return nil
}

// String returns a string representation of the Gun module.
func (G *Gun) String() string {
	return "Gun"
}
