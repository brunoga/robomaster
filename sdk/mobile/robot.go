package mobile

import (
	"github.com/brunoga/robomaster/sdk/modules/robot"
)

type Robot struct {
	r *robot.Robot
}

func (r *Robot) GetMotionMode() (int, error) {
	mm, err := r.r.GetMotionMode()
	if err != nil {
		return int(robot.MotionModeInvalid), err
	}

	return int(mm), nil
}

func (r *Robot) SetMotionMode(motionMode int) error {
	return r.r.SetMotionMode(robot.MotionMode(motionMode))
}

func (r *Robot) GetBatteryPercentage() (int, error) {
	return r.r.GetBatteryPercentage()
}
