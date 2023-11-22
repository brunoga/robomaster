package controller

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/robot"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

type Controller struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger

	rb *robot.Robot

	connRL *support.ResultListener
}

var _ module.Module = (*Controller)(nil)

func New(rb *robot.Robot, ub unitybridge.UnityBridge,
	l *logger.Logger) (*Controller, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("controller_module")

	c := &Controller{
		ub: ub,
		l:  l,
		rb: rb,
	}

	c.connRL = support.NewResultListener(ub, l,
		key.KeyMainControllerConnection, func(r *result.Result) {
			if r.ErrorCode() != 0 {
				return
			}

			if !r.Value().(bool) {
				return
			}

			c.rb.EnableFunction(robot.FunctionTypeMovementControl, true)
		})

	return c, nil
}

func (c *Controller) Start() error {
	return c.connRL.Start()
}

func (c *Controller) Stop() error {
	return c.connRL.Stop()
}

// Connected returns true if the connection to the robot is established.
func (c *Controller) Connected() bool {
	connected, ok := c.connRL.Result().Value().(bool)
	if !ok {
		return false
	}

	return connected
}

// WaitForConnection waits for the connection to the controller to be
// established.
func (c *Controller) WaitForConnection(timeout time.Duration) bool {
	connected, ok := c.connRL.Result().Value().(bool)
	if ok && connected {
		return true
	}

	return c.connRL.WaitForNewResult(timeout).Value().(bool)
}

func (c *Controller) Move(leftStick *StickPosition, rightStick *StickPosition,
	m ControlMode) error {
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

	return c.ub.DirectSendKeyValue(key.KeyMainControllerVirtualStick, v)
}

func (c *Controller) String() string {
	return "Controller"
}
