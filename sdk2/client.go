package sdk2

import (
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/camera"
	"github.com/brunoga/robomaster/sdk2/module/chassis"
	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/robomaster/sdk2/module/gamepad"
	"github.com/brunoga/robomaster/sdk2/module/gimbal"
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
	gm *gimbal.Gimbal
	rb *robot.Robot
	gn *gun.Gun
	gb *gamepad.GamePad

	m       sync.RWMutex
	started bool
}

// New creates a new Client instance with the given logger and appID. The appID
// parameter is used to determine which robot to connect to (i.e. it will only
// connect to robots broadcasting the given appID). If appID is 0, the client
// will connect to the first robot it finds.
//
// To get a robot to broadcast a given appID, use a QRCode to configure it (see
// https://github.com/brunoga/unitybridge/blob/main/support/qrcode/qrcode.go).
func New(l *logger.Logger, appID uint64) (*Client, error) {
	return new(l, appID, connection.TypeRouter, module.TypeAllButGamePad)
}

func NewWithModules(l *logger.Logger, appID uint64,
	modules module.Type) (*Client, error) {
	return new(l, appID, connection.TypeRouter, modules)
}

// NewWifiDirect creates a new Client instance with the given logger. This
// client will connect to the robot using WiFi Direct.
func NewWifiDirect(l *logger.Logger) (*Client, error) {
	return new(l, 0, connection.TypeWiFiDirect, module.TypeAllButGamePad)
}

func NewWifiDirectWithModules(l *logger.Logger,
	modules module.Type) (*Client, error) {
	return new(l, 0, connection.TypeWiFiDirect, modules)
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

	waitTimeout := 10 * time.Second

	// Connection.
	err = c.changeStateIfNonNil(c.cn, waitTimeout, true)
	if err != nil {
		return err
	}

	// Robot.
	err = c.changeStateIfNonNil(c.rb, waitTimeout, true)
	if err != nil {
		return err
	}

	// Wait for devices to be available.
	if !c.rb.WaitForDevices(waitTimeout) {
		return fmt.Errorf("robot working devices unexpectedly not established")
	}

	// Camera.
	err = c.changeStateIfNonNil(c.cm, waitTimeout, true)
	if err != nil {
		return err
	}

	// Chassis.
	err = c.changeStateIfNonNil(c.ch, waitTimeout, true)
	if err != nil {
		return err
	}

	// Gimbal.
	err = c.changeStateIfNonNil(c.gm, waitTimeout, true)
	if err != nil {
		return err
	}

	// Gun.
	err = c.changeStateIfNonNil(c.gn, waitTimeout, true)
	if err != nil {
		return err
	}

	// GamePad.
	err = c.changeStateIfNonNil(c.gb, waitTimeout, true)
	if err != nil {
		if err.Error() == "GamePad connection not established" {
			// GamePad is optional so it is fine it did not connect.
			c.l.Warn("GamePad connection not established.")
		} else {
			return err
		}
	}

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

// Gimbal returns the Gimbal module.
func (c *Client) Gimbal() *gimbal.Gimbal {
	return c.gm
}

// Robot returns the Robot module.
func (c *Client) Robot() *robot.Robot {
	return c.rb
}

// Gun returns the Gun module.
func (c *Client) Gun() *gun.Gun {
	return c.gn
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

	waitTime := 5 * time.Second

	// Gamepad.
	err := c.changeStateIfNonNil(c.gb, waitTime, false)
	if err != nil {
		return err
	}

	// Gun.
	err = c.changeStateIfNonNil(c.gn, waitTime, false)
	if err != nil {
		return err
	}

	// Chassis.
	err = c.changeStateIfNonNil(c.ch, waitTime, false)
	if err != nil {
		return err
	}

	// Camera.
	err = c.changeStateIfNonNil(c.cm, waitTime, false)
	if err != nil {
		return err
	}

	// Robot.
	err = c.changeStateIfNonNil(c.rb, waitTime, false)
	if err != nil {
		return err
	}

	// Connection.
	err = c.changeStateIfNonNil(c.cn, waitTime, false)
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

func new(l *logger.Logger, appID uint64, typ connection.Type,
	modules module.Type) (*Client, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	if modules&module.TypeConnection == 0 {
		return nil, fmt.Errorf("connection module is required")
	}

	if modules&module.TypeRobot == 0 {
		return nil, fmt.Errorf("robot module is required")
	}

	ub := unitybridge.Get(wrapper.Get(l), true, l)

	cn, err := connection.New(ub, l, appID, typ)
	if err != nil {
		return nil, err
	}

	rb, err := robot.New(ub, l, cn)
	if err != nil {
		return nil, err
	}

	var cm *camera.Camera
	if modules&module.TypeCamera != 0 {
		cm, err = camera.New(ub, l, cn)
		if err != nil {
			return nil, err
		}
	}

	var ch *chassis.Chassis
	if modules&module.TypeChassis != 0 {
		ch, err = chassis.New(ub, l, cn, rb)
		if err != nil {
			return nil, err
		}
	}

	var gm *gimbal.Gimbal
	if modules&module.TypeGimbal != 0 {
		gm, err = gimbal.New(ub, l, cn)
		if err != nil {
			return nil, err
		}
	}

	var gn *gun.Gun
	if modules&module.TypeGun != 0 {
		gn, err = gun.New(ub, l, cn, rb)
		if err != nil {
			return nil, err
		}
	}

	var gb *gamepad.GamePad
	if modules&module.TypeGamePad != 0 {
		gb, err = gamepad.New(ub, l, cn)
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		ub: ub,
		l:  l,
		cn: cn,
		rb: rb,
		cm: cm,
		gm: gm,
		ch: ch,
		gn: gn,
		gb: gb,
	}, nil
}

func (c *Client) changeStateIfNonNil(m module.Module, waitTime time.Duration,
	start bool) error {
	var err error
	if m != nil && !reflect.ValueOf(m).IsNil() {
		if start {
			err = m.Start()
		} else {
			err = m.Stop()
		}
	} else {
		return nil
	}

	if err != nil {
		return err
	}

	if !m.WaitForConnection(waitTime) {
		return fmt.Errorf("%s connection not established", m)
	}

	return nil
}
