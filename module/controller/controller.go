package controller

import (
	"fmt"
	"log/slog"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

// Controller is the robot's main controller interface. It is also responsibe for
// movement using the dual stick interface.
type Controller struct {
	*internal.BaseModule
}

var _ module.Module = (*Controller)(nil)

// New creates a new Controller instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection) (*Controller, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("controller_module")

	c := &Controller{}

	c.BaseModule = internal.NewBaseModule(ub, l, "Chassis",
		key.KeyMainControllerConnection, func(r *result.Result) {
			if !r.Succeeded() {
				c.Logger().Error("Connection: Unsuccessfull result.", "result", r)
				return
			}

			connectedValue, ok := r.Value().(*value.Bool)
			if !ok {
				c.Logger().Error("Connection: Unexpected value.", "value", r.Value())
				return
			}

			if connectedValue.Value {
				c.Logger().Debug("Connected.")
			} else {
				c.Logger().Debug("Disconnected.")
			}
		}, cm)

	return c, nil
}

// SetMode sets the controller mode for the robot.
func (c *Controller) SetMode(m Mode) error {
	if !m.Valid() {
		return fmt.Errorf("invalid controller mode: %d", m)
	}

	return c.UB().SetKeyValueSync(key.KeyMainControllerChassisCarControlMode,
		&value.Uint64{Value: uint64(m)})
}

// Move moves the robot using the given stick positions and control mode.
func (c *Controller) Move(chassisStick *StickPosition,
	gimbalStick *StickPosition, m Mode) error {
	if !m.Valid() {
		return fmt.Errorf("invalid controller mode: %d", m)
	}

	var leftStickEnabled uint64
	if chassisStick != nil {
		leftStickEnabled = 1
	}

	var rightStickEnabled uint64
	if gimbalStick != nil {
		rightStickEnabled = 1
	}

	v := chassisStick.InterpolatedY() |
		chassisStick.InterpolatedX()<<11 |
		gimbalStick.InterpolatedY()<<22 |
		gimbalStick.InterpolatedX()<<33 |
		leftStickEnabled<<44 |
		rightStickEnabled<<45 |
		uint64(m)<<46

	return c.UB().DirectSendKeyValue(key.KeyMainControllerVirtualStick, v)
}
