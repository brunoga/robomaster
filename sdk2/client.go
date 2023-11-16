package sdk2

import (
	"fmt"
	"sync"

	"github.com/brunoga/robomaster/sdk2/module/camera"
	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/wrapper"
)

type Client struct {
	l *logger.Logger

	ub unitybridge.UnityBridge

	cn *connection.Connection
	cm *camera.Camera

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

	cm, err := camera.New(ub, l)
	if err != nil {
		return nil, err
	}

	return &Client{
		ub: ub,
		l:  l,
		cn: cn,
		cm: cm,
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

	// Camera.
	err = c.cm.Start()
	if err != nil {
		return err
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

func (c *Client) Stop() error {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.started {
		return fmt.Errorf("client not started")
	}

	// Stop modules.

	// Camera.
	err := c.cm.Stop()
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
