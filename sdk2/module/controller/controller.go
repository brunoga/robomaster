package controller

import (
	"fmt"
	"sync"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

type Controller struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger

	mccToken token.Token

	m         sync.RWMutex
	connected bool
}

var _ module.Module = (*Controller)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Controller, error) {
	return &Controller{
		ub: ub,
		l:  l,
	}, nil
}

func (c *Controller) Start() error {
	var err error
	c.mccToken, err = c.ub.AddKeyListener(key.KeyMainControllerConnection, func(r *result.Result) {
		if r.ErrorCode() != 0 {
			c.l.Error("Error getting controller connection status", "error", r.ErrorDesc())
			return
		}

		c.l.Debug("Controller connection status", "status", r.Value())
		c.m.Lock()
		c.connected = r.Value().(bool)
		c.m.Unlock()
	}, true)
	if err != nil {
		return err
	}

	//err = c.ub.SetKeyValueSync(key.KeyMainControllerVirtualStickEnabled, true)
	//if err != nil {
	//	c.ub.RemoveKeyListener(key.KeyMainControllerConnection, c.mccToken)
	//}

	return err
}

func (c *Controller) Stop() error {
	return nil
}

func (c *Controller) Connected() bool {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.connected
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
