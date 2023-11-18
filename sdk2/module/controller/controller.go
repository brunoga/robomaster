package controller

import (
	"fmt"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/robot"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

type Controller struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger

	rb *robot.Robot

	mccToken token.Token

	connRL *support.ResultListener
}

var _ module.Module = (*Controller)(nil)

func New(rb *robot.Robot, ub unitybridge.UnityBridge,
	l *logger.Logger) (*Controller, error) {
	return &Controller{
		ub: ub,
		l:  l,
		rb: rb,
		connRL: support.NewResultListener(ub, l,
			key.KeyMainControllerConnection),
	}, nil
}

func (c *Controller) Start() error {
	return c.connRL.Start(func(r *result.Result) {
		if r.ErrorCode() != 0 {
			return
		}

		if !r.Value().(bool) {
			return
		}

		c.rb.EnableFunction(robot.FunctionTypeMovementControl, true)
	})
}

func (c *Controller) Stop() error {
	return c.connRL.Stop()
}

// Connected returns true if the connection to the robot is established.
func (c *Controller) Connected() bool {
	ok, connected := c.connRL.Result().Value().(bool)
	if !ok {
		return false
	}

	return connected
}

// WaitForConnection waits for the connection to the controller to be
// established.
func (c *Controller) WaitForConnection() bool {
	ok, connected := c.connRL.Result().Value().(bool)
	if ok && connected {
		return true
	}

	return c.connRL.WaitForNewResult(5 * time.Second).Value().(bool)
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
