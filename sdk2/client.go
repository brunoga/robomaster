package sdk2

import (
	"github.com/brunoga/robomaster/sdk2/internal/manager"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/wrapper"
)

type Client struct {
	l *logger.Logger

	ub unitybridge.UnityBridge
	cm *manager.Connection
}

func NewClient(l *logger.Logger, appID uint64) (*Client, error) {
	ub := unitybridge.Get(wrapper.Get(l), true, l)

	cm, err := manager.NewConnection(ub, l, appID)
	if err != nil {
		return nil, err
	}

	return &Client{
		ub: ub,
		l:  l,
		cm: cm,
	}, nil
}

func (c *Client) Connect() error {
	err := c.ub.Start()
	if err != nil {
		return err
	}

	return c.cm.Start()
}

func (c *Client) Disconnect() error {
	err := c.cm.Stop()
	if err != nil {
		return err
	}

	return c.ub.Stop()
}
