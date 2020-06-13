package modules

import (
	"fmt"
	"strconv"
)

// RobotMotionMode represents the motion mode for a robot.
type RobotMotionMode int

// Available robot motion modes.
const (
	RobotMotionModeChassisLead RobotMotionMode = iota // gimbal follows chassis
	RobotMotionModeGimbalLead                         // chassis follow gimbal
	RobotMotionModeFree                               // chassis and gimbal move independently
	RobotMotionModeInvalid                            // invalid mode
)

// Robot handles getting/setting robot specific attributes.
type Robot struct {
	control *Control
}

// NewRobot returns a new Robot instance associated with the given control.
func NewRobot(control *Control) *Robot {
	return &Robot{
		control,
	}
}

// GetMotionMode returns the robot's current motion mode and a nil error on
// success and a non-nil error on failure.
func (r *Robot) GetMotionMode() (RobotMotionMode, error) {
	data, err := r.control.SendAndReceiveData("robot mode ?;")
	if err != nil {
		return RobotMotionModeInvalid, fmt.Errorf(
			"error sending and receiving data: %w", err)
	}

	switch data {
	case "chassis_lead":
		return RobotMotionModeChassisLead, nil
	case "gimbal_lead":
		return RobotMotionModeGimbalLead, nil
	case "free":
		return RobotMotionModeFree, nil
	default:
		return RobotMotionModeInvalid, fmt.Errorf("unknown robot mode")
	}
}

// SetMotionMode sets the robot's current motion mode. Returns a nil error
// on success and a non-nil error on failure.
func (r *Robot) SetMotionMode(motionMode RobotMotionMode) error {
	setMotionMode := "robot mode "
	switch motionMode {
	case RobotMotionModeChassisLead:
		setMotionMode += "chassis_lead;"
	case RobotMotionModeGimbalLead:
		setMotionMode += "gimbal_lead;"
	case RobotMotionModeFree:
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
