package sdk2

import (
	"fmt"
	"sync"
	"time"

	"github.com/brunoga/robomaster/sdk2/module/camera"
	"github.com/brunoga/robomaster/sdk2/module/chassis"
	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/robomaster/sdk2/module/gamepad"
	"github.com/brunoga/robomaster/sdk2/module/gun"
	"github.com/brunoga/robomaster/sdk2/module/robot"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/wrapper"
)

type Client struct {
	l *logger.Logger

	ub unitybridge.UnityBridge

	cn *connection.Connection
	cm *camera.Camera
	ch *chassis.Chassis
	rb *robot.Robot
	gn *gun.Gun
	gb *gamepad.GamePad

	m       sync.RWMutex
	started bool
}

// New creates a new Client instance with the given logger and appID.
func New(l *logger.Logger, appID uint64) (*Client, error) {
	ub := unitybridge.Get(wrapper.Get(l), true, l)

	cn, err := connection.New(ub, l, appID)
	if err != nil {
		return nil, err
	}

	rb, err := robot.New(ub, l)
	if err != nil {
		return nil, err
	}

	cm, err := camera.New(ub, l)
	if err != nil {
		return nil, err
	}

	ch, err := chassis.New(rb, ub, l)
	if err != nil {
		return nil, err
	}

	gn, err := gun.New(ub, l)
	if err != nil {
		return nil, err
	}

	gb, err := gamepad.New(ub, l)
	if err != nil {
		return nil, err
	}

	return &Client{
		ub: ub,
		l:  l,
		cn: cn,
		rb: rb,
		cm: cm,
		ch: ch,
		gn: gn,
		gb: gb,
	}, nil
}

// Start starts the client and all associated modules.
func (c *Client) Start() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.started {
		return fmt.Errorf("client already started")
	}

	err := c.ub.Start()
	if err != nil {
		return err
	}

	// Start modules.

	// Connection.
	err = c.cn.Start()
	if err != nil {
		return err
	}

	if !c.cn.WaitForConnection(20 * time.Second) {
		return fmt.Errorf("network connection unexpectedly not established")
	}

	// Robot.
	err = c.rb.Start()
	if err != nil {
		return err
	}

	if !c.rb.WaitForConnection(10 * time.Second) {
		return fmt.Errorf("robot connection unexpectedly not established")
	}

	// Camera.
	err = c.cm.Start()
	if err != nil {
		return err
	}

	if !c.cm.WaitForConnection(10 * time.Second) {
		return fmt.Errorf("camera connection unexpectedly not established")
	}

	// Chassis.
	err = c.ch.Start()
	if err != nil {
		return err
	}

	if !c.ch.WaitForConnection(10 * time.Second) {
		return fmt.Errorf("chassis connection unexpectedly not established")
	}

	// Gun.
	err = c.gn.Start()
	if err != nil {
		return err
	}

	if !c.gn.WaitForConnection(10 * time.Second) {
		return fmt.Errorf("gun connection unexpectedly not established")
	}

	// GamePad. (Optional)
	err = c.gb.Start()
	if err != nil {
		return err
	}

	go func() {
		if !c.gb.WaitForConnection(2 * time.Second) {
			// GamePad is optional.
			c.l.Warn("Gamepad connection not stablished. Gamepad not available.")
			err := c.gb.Stop()
			if err != nil {
				c.l.Warn("Error stopping Gamepad module", "error", err)
			}
			c.gb = nil
		}
	}()

	c.started = true

	return nil
}

// Connection returns the Connection module.
func (c *Client) Connection() *connection.Connection {
	return c.cn
}

// Camera returns the Camera module.
func (c *Client) Camera() *camera.Camera {
	return c.cm
}

// Chassis returns the Chassis module.
func (c *Client) Chassis() *chassis.Chassis {
	return c.ch
}

// Robot returns the Robot module.
func (c *Client) Robot() *robot.Robot {
	return c.rb
}

// GamePad returns the GamePad module. The GamePad is optional and may be nil.
func (c *Client) GamePad() *gamepad.GamePad {
	return c.gb
}

// Stop stops the client and all associated modules.
func (c *Client) Stop() error {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.started {
		return fmt.Errorf("client not started")
	}

	// Stop modules.

	// Chassis.
	err := c.ch.Stop()
	if err != nil {
		return err
	}

	// Camera.
	err = c.cm.Stop()
	if err != nil {
		return err
	}

	// Robot.
	err = c.rb.Stop()
	if err != nil {
		return err
	}

	// Connection.
	err = c.cn.Stop()
	if err != nil {
		return err
	}

	// Stop Unity Bridge.
	err = c.ub.Stop()
	if err != nil {
		return err
	}

	return nil
}
