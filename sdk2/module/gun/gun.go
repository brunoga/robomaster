package gun

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/robomaster/sdk2/module/robot"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
)

// Gun is the module that controls turret firing. It supports both infrared and
// beads firing.
type Gun struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger
	rm *robot.Robot
	cm *connection.Connection
}

var _ module.Module = (*Gun)(nil)

// New creates a new Gun instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection, rm *robot.Robot) (*Gun, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("gun_module")

	return &Gun{
		ub: ub,
		l:  l,
		rm: rm,
		cm: cm,
	}, nil
}

// Start starts the Gun module.
func (g *Gun) Start() error {
	return g.rm.EnableFunction(robot.FunctionTypeGunControl, true)
}

// Connected returns whether the Gun module is connected.
func (g *Gun) Connected() bool {
	return g.rm.HasDevice(robot.DeviceTypeWaterGun) && g.cm.Connected()
}

// WaitForConnection waits for the Gun module to connect and returns the
// connected status.
func (g *Gun) WaitForConnection(timeout time.Duration) bool {
	start := time.Now()

	if !g.cm.WaitForConnection(timeout) {
		return false
	}

	timeout = timeout - time.Since(start)

	if !g.rm.WaitForDevices(timeout) {
		return false
	}

	return g.rm.HasDevice(robot.DeviceTypeWaterGun)
}

// Fire fires the Gun module with the given type.
func (g *Gun) Fire(typ Type) error {
	switch typ {
	case TypeBead:
		// TODO(bga): Maybe implement firing multiple beads and extend it to IR
		//            firing too.
		return g.fireBead(1)
	case TypeInfrared:
		return g.fireInfrared()
	}

	return fmt.Errorf("invalid gun type: %v", typ)
}

// Stop stops the Gun module.
func (g *Gun) Stop() error {
	return g.rm.EnableFunction(robot.FunctionTypeGunControl, false)
}

// String returns a string representation of the Gun module.
func (G *Gun) String() string {
	return "Gun"
}

type timesValue struct {
	Value uint64 `json:"value"`
}

func (g *Gun) fireBead(times uint64) error {
	return g.ub.PerformActionForKey(key.KeyRobomasterWaterGunWaterGunFireWithTimes,
		timesValue{times}, nil)
}

func (g *Gun) fireInfrared() error {
	go func() {
		// Disable firing after a while.
		time.Sleep(200 * time.Millisecond)

		g.ub.DirectSendKeyValue(key.KeyRobomasterWaterGunWaterGunFire, uint64(0))
	}()

	return g.ub.DirectSendKeyValue(key.KeyRobomasterWaterGunWaterGunFire, uint64(1))
}
