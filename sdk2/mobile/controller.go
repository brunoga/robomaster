package mobile

import "github.com/brunoga/robomaster/sdk2/module/controller"

// StickPosition represents the position of a stick controller.
type StickPosition struct {
	sp *controller.StickPosition
}

// Set sets the position of the stick controller. X and Y must be between 0.0
// and 1.0.
func (s *StickPosition) Set(x, y float64) {
	s.sp.X = x
	s.sp.Y = y
}

// InterpolatedX returns the interpolated X position of the stick controller.
func (s *StickPosition) InterpolatedX() int64 {
	return int64(s.sp.InterpolatedX())
}

// InterpolatedY returns the interpolated Y position of the stick controller.
func (s *StickPosition) InterpolatedY() int64 {
	return int64(s.sp.InterpolatedY())
}

// Controller allows controlling the robot movement by using a method similar
// to a game controller with dual sticks.
type Controller struct {
	c *controller.Controller
}

// Move moves the robot using the given stick positions. The left stick controls
// the robot's gimbal and the right stick controls the robot's chassis.
func (c *Controller) Move(ls *StickPosition, rs *StickPosition) error {
	return c.c.Move(ls.sp, rs.sp, controller.ControlModeDefault)
}
