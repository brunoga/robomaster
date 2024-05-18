package mobile

import (
	"github.com/brunoga/robomaster/module/chassis"
	"github.com/brunoga/robomaster/module/chassis/controller"
)

// StickPosition represents the position of a stick controller.
type StickPosition struct {
	sp *controller.StickPosition
}

// NewStickPosition creates a new StickPosition with the given X and Y
// positions.
func NewStickPosition(x, y float64) *StickPosition {
	return &StickPosition{
		sp: &controller.StickPosition{
			X: x,
			Y: y,
		},
	}
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

// Chassis allows controlling the robot chassis. It also works as the robot main
// controller interface.
type Chassis struct {
	c *chassis.Chassis
}

// Move moves the robot using the given stick positions. The left stick controls
// the robot's gimbal and chassis rotation and the right stick controls the
// robot's chassis (up/down/left/right).
func (c *Chassis) Move(ls *StickPosition, rs *StickPosition) error {
	return c.c.Move(ls.sp, rs.sp, controller.ModeFPV)
}
