package chassis

import (
	"fmt"
	"log/slog"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/chassis/controller"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/robomaster/sdk2/module/robot"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

// Chassis allows controlling the robot chassis. It also works as the robot main
// controller interface.
type Chassis struct {
	*internal.BaseModule

	rb *robot.Robot
}

var _ module.Module = (*Chassis)(nil)

// New creates a new Chassis instance.
func New(rb *robot.Robot, ub unitybridge.UnityBridge,
	l *logger.Logger) (*Chassis, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("chassis_module")

	c := &Chassis{
		rb: rb,
	}

	c.BaseModule = internal.NewBaseModule(ub, l, "Chassis",
		key.KeyMainControllerConnection, func(r *result.Result) {
			if r == nil || r.ErrorCode() != 0 {
				return
			}

			if connected, ok := r.Value().(bool); !ok || !connected {
				return
			}

			// TODO(bga): Maybe disable the function if we receive an actual
			//            false here?

			c.rb.EnableFunction(robot.FunctionTypeMovementControl, true)
		})

	return c, nil
}

// Move moves the robot using the given stick positions and control mode.
func (c *Chassis) Move(leftStick *controller.StickPosition,
	rightStick *controller.StickPosition, m controller.ControlMode) error {
	if !m.Valid() {
		return fmt.Errorf("invalid control mode: %d", m)
	}

	var leftStickEnabled uint64
	if leftStick != nil {
		leftStickEnabled = 1
	}

	var rightStickEnabled uint64
	if rightStick != nil {
		rightStickEnabled = 1
	}

	v := leftStick.InterpolatedY() |
		leftStick.InterpolatedX()<<11 |
		rightStick.InterpolatedY()<<22 |
		rightStick.InterpolatedX()<<33 |
		leftStickEnabled<<44 |
		rightStickEnabled<<45 |
		uint64(m)<<46

	return c.UB().DirectSendKeyValue(key.KeyMainControllerVirtualStick, v)
}
