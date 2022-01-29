package controller

import (
	"git.bug-br.org.br/bga/robomasters1/app/internal"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
)

const (
	interpolateLowerBound = float32(364)
	interpolateUpperBound = float32(1684)
)

type Controller struct {
	cc *internal.CommandController
}

func New(cc *internal.CommandController) *Controller {
	// Enable movement control.
	param := &functionEnableParameter{make([]*functionEnableInfo, 0)}
	param.setFunctionEnable(movementControl, true)
	cc.PerformAction(dji.KeyRobomasterSystemFunctionEnable, param, nil)

	return &Controller{
		cc,
	}
}

func (c *Controller) Move(chassisY, chassisX, gimbalY, gimbalX float32,
	chassisControl, gimbalControl bool, controlMode uint64) {
	// Set control mode.
	//
	// TODO(bga): Maybe we do not need to send this all the time but the DJI
	//  app does so we are doing it too.
	c.cc.DirectSendValue(dji.KeyGimbalControlMode, 0)

	// Setup movement command.
	interpolatedChassisX := interpolate(chassisX)
	interpolatedChassisY := interpolate(chassisY)
	interpolatedGimbalX := interpolate(gimbalX)
	interpolatedGimbalY := interpolate(gimbalY)

	var intChassisControl uint64 = 0
	if chassisControl {
		intChassisControl = 1
	}

	var intGimbalControl uint64 = 0
	if chassisControl {
		intGimbalControl = 1
	}

	value := uint64(interpolatedChassisY)
	value |= uint64(interpolatedChassisX) << 11
	value |= uint64(interpolatedGimbalY) << 22
	value |= uint64(interpolatedGimbalX) << 33
	value |= intChassisControl << 44
	value |= intGimbalControl << 45
	value |= controlMode << 46

	// Send movement command.
	c.cc.DirectSendValue(dji.KeyMainControllerVirtualStick, value)
}

func interpolate(value float32) float32 {
	if value < 0.0 || value > 1.0 {
		panic("Value to interpolate must be inside [0, 1].")
	}

	return (float32(1.0)-value)*interpolateLowerBound + value*
		interpolateUpperBound
}
