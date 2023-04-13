package mobile

import (
	"sync"

	"github.com/brunoga/robomaster/sdk/modules/chassis"
)

type Speed struct {
	m sync.RWMutex
	x float64
	y float64
	z float64
}

type Chassis struct {
	c *chassis.Chassis
}

func (c *Chassis) SetSpeed(speed *Speed, async bool) error {
	convSpeed := chassis.NewSpeed(speed.x, speed.y, speed.z)
	return c.c.SetSpeed(convSpeed, async)
}
