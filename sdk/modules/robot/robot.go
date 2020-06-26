package robot

import (
	"fmt"
	"strconv"

	"github.com/brunoga/robomaster/sdk/modules"
)

// RobotMotionMode represents the motion mode for a robot.
type MotionMode int

// Available robot motion modes.
const (
	MotionModeChassisLead MotionMode = iota // gimbal follows chassis
	MotionModeGimbalLead                    // chassis follow gimbal
	MotionModeFree                          // chassis and gimbal move independently
	MotionModeInvalid                       // invalid mode
)

// Robot handles getting/setting robot specific attributes.
type Robot struct {
	control *modules.Control
}

// New returns a new Robot instance associated with the given control.
func New(control *modules.Control) *Robot {
	return &Robot{
		control,
	}
}

// GetMotionMode returns the robot's current motion mode and a nil error on
// success and a non-nil error on failure.
func (r *Robot) GetMotionMode() (MotionMode, error) {
	data, err := r.control.SendAndReceiveData("robot mode ?;")
	if err != nil {
		return MotionModeInvalid, fmt.Errorf(
			"error sending and receiving data: %w", err)
	}

	switch data {
	case "chassis_lead":
		return MotionModeChassisLead, nil
	case "gimbal_lead":
		return MotionModeGimbalLead, nil
	case "free":
		return MotionModeFree, nil
	default:
		return MotionModeInvalid, fmt.Errorf("unknown robot mode")
	}
}

// SetMotionMode sets the robot's current motion mode. Returns a nil error
// on success and a non-nil error on failure.
func (r *Robot) SetMotionMode(motionMode MotionMode) error {
	setMotionMode := "robot mode "
	switch motionMode {
	case MotionModeChassisLead:
		setMotionMode += "chassis_lead;"
	case MotionModeGimbalLead:
		setMotionMode += "gimbal_lead;"
	case MotionModeFree:
		setMotionMode += "free;"
	default:
		return fmt.Errorf("unknown robot mode: %d", motionMode)
	}

	err := r.control.SendDataExpectOk(setMotionMode)
	if err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}

	return nil
}

// GetBatteryPercentage returns the robot's current battery percentage and a
// nil error on success and a non-nil error on failure.
func (r *Robot) GetBatteryPercentage() (int, error) {
	data, err := r.control.SendAndReceiveData("robot battery ?;")
	if err != nil {
		return -1, err
	}

	percentage, err := strconv.Atoi(data)
	if err != nil {
		return -1, fmt.Errorf("error parsing battery percentage")
	}

	return percentage, nil
}
