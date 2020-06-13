package robot

import (
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules"
	"strconv"
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

// NewRobot returns a new Robot instance associated with the given control.
func NewRobot(control *modules.Control) *Robot {
	return &Robot{
		control,
	}
}

// GetMotionMode returns the robot's current motion mode and a nil error on
// success and a non-nil error on failure.
func (r *Robot) GetMotionMode() (GetMotionModeResponse, error) {
	data, err := r.control.SendAndReceiveData("robot mode ?;")
	if err != nil {
		return GetMotionModeResponse{
			MotionModeInvalid,
		}, fmt.Errorf("error sending and receiving data: %w", err)
	}

	switch data {
	case "chassis_lead":
		return GetMotionModeResponse{
			MotionModeChassisLead,
		}, nil
	case "gimbal_lead":
		return GetMotionModeResponse{
			MotionModeGimbalLead,
		}, nil
	case "free":
		return GetMotionModeResponse{
			MotionModeFree,
		}, nil
	}

	return GetMotionModeResponse{
		MotionModeInvalid,
	}, fmt.Errorf("unknown robot mode")
}

// SetMotionMode sets the robot's current motion mode. Returns a nil error
// on success and a non-nil error on failure.
func (r *Robot) SetMotionMode(req SetMotionModeRequest) error {
	setMotionMode := "robot mode "
	switch req.MotionMode {
	case MotionModeChassisLead:
		setMotionMode += "chassis_lead;"
	case MotionModeGimbalLead:
		setMotionMode += "gimbal_lead;"
	case MotionModeFree:
		setMotionMode += "free;"
	default:
		return fmt.Errorf("unknown robot mode: %d", req.MotionMode)
	}

	err := r.control.SendDataExpectOk(setMotionMode)
	if err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}

	return nil
}

// GetBatteryPercentage returns the robot's current battery percentage and a
// nil error on success and a non-nil error on failure.
func (r *Robot) GetBatteryPercentage() (GetBatteryPercentageResponse, error) {
	data, err := r.control.SendAndReceiveData("robot battery ?;")
	if err != nil {
		return GetBatteryPercentageResponse{
			-1,
		}, err
	}

	percentage, err := strconv.Atoi(data)
	if err != nil {
		return GetBatteryPercentageResponse{
			-1,
		}, fmt.Errorf("error parsing battery percentage")
	}

	return GetBatteryPercentageResponse{
		percentage,
	}, nil
}
