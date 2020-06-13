package modules

import (
	"fmt"
	push2 "github.com/brunoga/robomaster/sdk/modules/push"
)

// Chassis allows sending commands to control the robot's chassis.
type Chassis struct {
	control *Control
	push    *push2.Push
}

// NewChassis returns a new Chassis instance associated with the given control.
func NewChassis(control *Control, push *push2.Push) *Chassis {
	return &Chassis{
		control,
		push,
	}
}

func (c *Chassis) SetSpeed(forwardSpeedMetersPerSecond,
	lateralSpeedMetersPerSecond,
	rotationSpeedDegreesPerSecond float64) error {
	return c.control.SendDataExpectOk(fmt.Sprintf(
		"chassis speed x %f y %f z %f;", forwardSpeedMetersPerSecond,
		lateralSpeedMetersPerSecond, rotationSpeedDegreesPerSecond))
}

func (c *Chassis) SetWheelSpeed(frontRightWheelRPM, frontLeftWheelRPM,
	rearRightWheelRPM, rearLeftWheelRPM float64) error {
	return c.control.SendDataExpectOk(fmt.Sprintf(
		"chassis wheel w1 %f w2 %f w3 %f w4 %f;", frontRightWheelRPM,
		frontLeftWheelRPM, rearRightWheelRPM, rearLeftWheelRPM))
}

func (c *Chassis) MoveRelative(forwardDistanceMeters,
	lateralDistanceMeters float64, rotationDegrees int32,
	moveSpeedMetersPerSecond, rotationSpeedDegreesPerSecond float64) error {
	return c.control.SendDataExpectOk(fmt.Sprintf(
		"chassis move x %f y %f z %d vxy %f vz %f;", forwardDistanceMeters,
		lateralDistanceMeters, rotationDegrees, moveSpeedMetersPerSecond,
		rotationSpeedDegreesPerSecond))
}
