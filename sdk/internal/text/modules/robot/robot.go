package robot

import (
	"fmt"
	"strconv"

	"github.com/brunoga/robomaster/sdk/internal/text/modules/control"
	"github.com/brunoga/robomaster/sdk/modules/robot"
)

// Robot handles getting/setting robot specific attributes.
type Robot struct {
	control *control.Control
}

var _ robot.Robot = (*Robot)(nil)

// New returns a new Robot instance associated with the given control.
func New(control *control.Control) *Robot {
	return &Robot{
		control,
	}
}

// GetMode returns the robot's current motion mode and a nil error on
// success and a non-nil error on failure.
func (r *Robot) GetMode() (robot.Mode, error) {
	data, err := r.control.SendAndReceiveData("robot mode ?;")
	if err != nil {
		return robot.ModeFree, fmt.Errorf(
			"error sending and receiving data: %w", err)
	}

	switch data {
	case "chassis_lead":
		return robot.ModeChassisLead, nil
	case "gimbal_lead":
		return robot.ModeGimbalLead, nil
	case "free":
		return robot.ModeFree, nil
	default:
		return robot.ModeFree, fmt.Errorf("unknown robot mode")
	}
}

// SetMode sets the robot's current motion mode. Returns a nil error
// on success and a non-nil error on failure.
func (r *Robot) SetMode(motionMode robot.Mode) error {
	setMotionMode := "robot mode "
	switch motionMode {
	case robot.ModeChassisLead:
		setMotionMode += "chassis_lead;"
	case robot.ModeGimbalLead:
		setMotionMode += "gimbal_lead;"
	case robot.ModeFree:
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

// GetproductVersion returns the robot's product version and a nil error on
// success and a non-nil error on failure.
func (r *Robot) GetProductVersion() (string, error) {
	data, err := r.control.SendAndReceiveData("version ?;")
	if err != nil {
		return "", err
	}

	return data[8 : len(data)-1], nil
}
