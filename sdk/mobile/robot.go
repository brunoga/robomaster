package mobile

import (
	"github.com/brunoga/robomaster/sdk/modules/robot"
)

type MotionMode int

const (
	MotionModeChassisLead MotionMode = iota
	MotionModeGimbalLead
	MotionModeFree
	MotionModeInvalid
)

type Robot struct {
	r *robot.Robot
}

func (r *Robot) GetMotionMode() (MotionMode, error) {
	mm, err := r.r.GetMotionMode()
	if err != nil {
		return MotionModeInvalid, err
	}

	return MotionMode(mm), nil
}

func (r *Robot) SetMotionMode(motionMode MotionMode) error {
	return r.r.SetMotionMode(robot.MotionMode(motionMode))
}

func (r *Robot) GetBatteryPercentage() (int, error) {
	return r.r.GetBatteryPercentage()
}
