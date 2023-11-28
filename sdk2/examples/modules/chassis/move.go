package main

import (
	"time"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/sdk2/module/chassis"
)

func main() {
	c, err := sdk2.New(nil, 0)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer c.Stop()

	cs := c.Chassis()

	err = cs.SetSpeed(chassis.ModeFPV, 0.7, 0.7, 0.7)
	if err != nil {
		panic(err)
	}

	err = cs.SetPosition(chassis.ModeFPV, 0.5, 0.0, 0.0)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
}
