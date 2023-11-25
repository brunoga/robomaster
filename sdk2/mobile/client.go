package mobile

import (
	"log/slog"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/unitybridge/support/logger"
)

// Client is the main entry point for the mobile SDK.
type Client struct {
	c *sdk2.Client
}

// NewClient creates a new Client instance.
func NewClient() (*Client, error) {
	l := logger.New(slog.LevelError)
	c, err := sdk2.New(l, 0)
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

// GamePad returns the GamePad instance for the client. The GamePad is optional
// and may be nil.
func (c *Client) GamePad() *GamePad {
	return &GamePad{
		g: c.c.GamePad(),
	}
}

// Stop stops the client.
func (c *Client) Stop() error {
	return c.c.Stop()
}
