package mobile

import "github.com/brunoga/robomaster/sdk/modules/chassis"

type Speed = chassis.Speed

type Chassis struct {
	c *chassis.Chassis
}

func (c *Chassis) SetSpeed(speed *Speed, async bool) error {
	return c.c.SetSpeed(speed, async)
}
