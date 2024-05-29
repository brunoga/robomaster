package robomaster

import (
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/camera"
	"github.com/brunoga/robomaster/module/chassis"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/controller"
	"github.com/brunoga/robomaster/module/gamepad"
	"github.com/brunoga/robomaster/module/gimbal"
	"github.com/brunoga/robomaster/module/gun"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
	"github.com/brunoga/robomaster/unitybridge/wrapper"
)

type Client struct {
	l *logger.Logger

	ub unitybridge.UnityBridge

	connectionModule *connection.Connection
	cameraModule     *camera.Camera
	chassisModule    *chassis.Chassis
	gimbalModule     *gimbal.Gimbal
	robotModule      *robot.Robot
	gunModule        *gun.Gun
	gamePadModule    *gamepad.GamePad
	controllerModule *controller.Controller

	m       sync.RWMutex
	started bool
}

// New creates a new Client instance with the given logger and appID. The appID
// parameter is used to determine which robot to connect to (i.e. it will only
// connect to robots broadcasting the given appID). If appID is 0, the client
// will connect to the first robot it finds.
//
// To get a robot to broadcast a given appID, use a QRCode to configure it (see
// https://github.com/brunoga/robomaster/unitybridge/blob/main/support/qrcode/qrcode.go).
func New(l *logger.Logger, appID uint64) (*Client, error) {
	return new(l, appID, connection.TypeRouter, module.TypeAll)
}

// NewWithModules is like New but allows selecting which mkodules to enable.
// The Connection and Robot modules are required.
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
	err = c.changeStateIfNonNil(c.connectionModule, waitTimeout, true)
	if err != nil {
		return err
	}

	// Robot.
	err = c.changeStateIfNonNil(c.robotModule, waitTimeout, true)
	if err != nil {
		return err
	}

	// Wait for devices to be available.
	if !c.robotModule.WaitForDevices(waitTimeout) {
		return fmt.Errorf("robot working devices unexpectedly not established")
	}

	// Controller.
	err = c.changeStateIfNonNil(c.controllerModule, waitTimeout, true)
	if err != nil {
		return err
	}

	// Camera.
	err = c.changeStateIfNonNil(c.cameraModule, waitTimeout, true)
	if err != nil {
		return err
	}

	// Chassis.
	err = c.changeStateIfNonNil(c.chassisModule, waitTimeout, true)
	if err != nil {
		return err
	}

	// Gimbal.
	err = c.changeStateIfNonNil(c.gimbalModule, waitTimeout, true)
	if err != nil {
		return err
	}

	// Gun.
	err = c.changeStateIfNonNil(c.gunModule, waitTimeout, true)
	if err != nil {
		return err
	}

	// GamePad.
	go func() {
		err = c.changeStateIfNonNil(c.gamePadModule, waitTimeout, true)
		if err != nil {
			if err.Error() == "GamePad connection not established" {
				// GamePad is optional so it is fine it did not connect.
				c.l.Warn("GamePad connection not established.")
			} else {
				c.l.Error("GamePad connection error: %v", err)
			}
		}
	}()

	c.started = true

	return nil
}

// Connection returns the Connection module.
func (c *Client) Connection() *connection.Connection {
	return c.connectionModule
}

// Camera returns the Camera module.
func (c *Client) Camera() *camera.Camera {
	return c.cameraModule
}

// Chassis returns the Chassis module.
func (c *Client) Chassis() *chassis.Chassis {
	return c.chassisModule
}

// Gimbal returns the Gimbal module.
func (c *Client) Gimbal() *gimbal.Gimbal {
	return c.gimbalModule
}

// Robot returns the Robot module.
func (c *Client) Robot() *robot.Robot {
	return c.robotModule
}

// Gun returns the Gun module.
func (c *Client) Gun() *gun.Gun {
	return c.gunModule
}

// GamePad returns the GamePad module. The GamePad is optional and may be nil.
func (c *Client) GamePad() *gamepad.GamePad {
	return c.gamePadModule
}

// Controller returns the Controller module.
func (c *Client) Controller() *controller.Controller {
	return c.controllerModule
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
	err := c.changeStateIfNonNil(c.gamePadModule, waitTime, false)
	if err != nil {
		return err
	}

	// Gun.
	err = c.changeStateIfNonNil(c.gunModule, waitTime, false)
	if err != nil {
		return err
	}

	// Chassis.
	err = c.changeStateIfNonNil(c.chassisModule, waitTime, false)
	if err != nil {
		return err
	}

	// Camera.
	err = c.changeStateIfNonNil(c.cameraModule, waitTime, false)
	if err != nil {
		return err
	}

	// Controller.
	err = c.changeStateIfNonNil(c.controllerModule, waitTime, false)
	if err != nil {
		return err
	}

	// Robot.
	err = c.changeStateIfNonNil(c.robotModule, waitTime, false)
	if err != nil {
		return err
	}

	// Connection.
	err = c.changeStateIfNonNil(c.connectionModule, waitTime, false)
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

	// Enable Unity Bridge debug logging if the logger level is trace.
	unityBridgeDebugEnabled := l.Level() == logger.LevelTrace

	ub := unitybridge.Get(wrapper.Get(l), unityBridgeDebugEnabled, l)

	connectionModule, err := connection.New(ub, l, appID, typ)
	if err != nil {
		return nil, err
	}

	robotModule, err := robot.New(ub, l, connectionModule)
	if err != nil {
		return nil, err
	}

	var cameraModule *camera.Camera
	if modules&module.TypeCamera != 0 {
		cameraModule, err = camera.New(ub, l, connectionModule)
		if err != nil {
			return nil, err
		}
	}

	var chassisModule *chassis.Chassis
	if modules&module.TypeChassis != 0 {
		chassisModule, err = chassis.New(ub, l, connectionModule, robotModule)
		if err != nil {
			return nil, err
		}
	}

	var gimbalModule *gimbal.Gimbal
	if modules&module.TypeGimbal != 0 {
		gimbalModule, err = gimbal.New(ub, l, connectionModule)
		if err != nil {
			return nil, err
		}
	}

	var gunModule *gun.Gun
	if modules&module.TypeGun != 0 {
		gunModule, err = gun.New(ub, l, connectionModule, robotModule)
		if err != nil {
			return nil, err
		}
	}

	var gamePadModule *gamepad.GamePad
	if modules&module.TypeGamePad != 0 {
		gamePadModule, err = gamepad.New(ub, l, connectionModule)
		if err != nil {
			return nil, err
		}
	}

	var controllerModule *controller.Controller
	if modules&module.TypeController != 0 {
		controllerModule, err = controller.New(ub, l, connectionModule)
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		ub:               ub,
		l:                l,
		connectionModule: connectionModule,
		robotModule:      robotModule,
		cameraModule:     cameraModule,
		gimbalModule:     gimbalModule,
		chassisModule:    chassisModule,
		gunModule:        gunModule,
		gamePadModule:    gamePadModule,
		controllerModule: controllerModule,
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

	if start && !m.WaitForConnection(waitTime) {
		return fmt.Errorf("%s connection not established", m)
	}

	return nil
}
