package gun

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
)

// Gun is the module that controls turret firing. It supports both infrared and
// beads firing.
type Gun struct {
	*internal.BaseModule

	rm *robot.Robot
}

var _ module.Module = (*Gun)(nil)

// New creates a new Gun instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	rm *robot.Robot) (*Gun, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("gun_module")

	bm, err := internal.NewBaseModule(ub, l, "Gun", nil, nil, rm)
	if err != nil {
		return nil, fmt.Errorf("failed to create base module: %w", err)
	}

	g := &Gun{
		rm: rm,
	}

	g.BaseModule = bm

	return g, nil
}

// Start starts the Gun module.
func (g *Gun) Start() error {
	err := g.rm.EnableFunction(robot.FunctionTypeGunControl, true)
	if err != nil {
		return fmt.Errorf("failed to enable gun control: %w", err)
	}

	return g.BaseModule.Start()
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
	err := g.BaseModule.Stop()
	if err != nil {
		g.Logger().Warn("Failed to stop base module.", "error", err)
	}

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
	return g.UB().PerformActionForKey(key.KeyRobomasterWaterGunWaterGunFireWithTimes,
		timesValue{times}, nil)
}

func (g *Gun) fireInfrared() error {
	go func() {
		// Disable firing after a while.
		time.Sleep(200 * time.Millisecond)

		g.UB().DirectSendKeyValue(key.KeyRobomasterWaterGunWaterGunFire, uint64(0))
	}()

	return g.UB().DirectSendKeyValue(key.KeyRobomasterWaterGunWaterGunFire, uint64(1))
}
