package mobile

import "github.com/brunoga/robomaster/module/robot"

// Robot allows reading the robot parameters.
type Robot struct {
	r *robot.Robot
}

// BatteryPowerPercent returns the current battery power percent (0 to 100).
func (r *Robot) BatteryPowerPercent() int8 {
	return int8(r.r.BatteryPowerPercent())
}
