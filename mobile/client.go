package mobile

import (
	"log/slog"

	sdk2 "github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
)

const (
	// All modules except for Gun (for now).
	mobileModules = module.TypeConnection | module.TypeRobot |
		module.TypeChassis | module.TypeGimbal | module.TypeCamera | module.TypeGamePad
)

// Client is the main entry point for the mobile SDK.
type Client struct {
	c *sdk2.Client
}

// NewClient creates a new Client instance. If appID is 0, the client will try
// to connect to the first available Robomaster robot. If it is non-zero, it
// will only connect to a robot that is broadcasting the given appID. The appID
// can be configured in the robot through a qrcode.
func NewClient(appID int64) (*Client, error) {
	l := logger.New(slog.LevelDebug)
	c, err := sdk2.NewWithModules(l, uint64(appID), mobileModules)
	if err != nil {
		return nil, err
	}

	return &Client{
		c: c,
	}, nil
}

// NewWifiDirectClient creates a new Client instance that will connect to a
// Robomaster robot using Wifi Direct.
func NewWifiDirectClient() (*Client, error) {
	l := logger.New(slog.LevelDebug)
	c, err := sdk2.NewWifiDirectWithModules(l, mobileModules)
	if err != nil {
		return nil, err
	}

	return &Client{
		c: c,
	}, nil
}

// Start starts the client.
func (c *Client) Start() error {
	return c.c.Start()
}

// Camera returns the Camera instance for the client.
func (c *Client) Camera() *Camera {
	return &Camera{
		c: c.c.Camera(),
	}
}

// Controller returns the Controller instance for the client.
func (c *Client) Chassis() *Chassis {
	return &Chassis{
		c: c.c.Chassis(),
	}
}

// Connnection returns the Connection instance for the client.
func (c *Client) Connection() *Connection {
	return &Connection{
		c: c.c.Connection(),
	}
}

// GamePad returns the GamePad instance for the client. The GamePad is optional
// and may be nil.
func (c *Client) GamePad() *GamePad {
	return &GamePad{
		g: c.c.GamePad(),
	}
}

// Robot returns the Robot instance for the client.
func (c *Client) Robot() *Robot {
	return &Robot{
		r: c.c.Robot(),
	}
}

// Stop stops the client.
func (c *Client) Stop() error {
	return c.c.Stop()
}
