package components

import (
	"log/slog"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/unitybridge/support/logger"
)

type Robomaster struct {
	c *sdk2.Client
}

func NewRobomaster() (*Robomaster, error) {
	l := logger.New(slog.LevelDebug)

	c, err := sdk2.New(l, 0)
	if err != nil {
		return nil, err
	}

	return &Robomaster{
		c: c,
	}, nil
}

func (r *Robomaster) Client() *sdk2.Client {
	return r.c
}
