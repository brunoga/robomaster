package command

import "github.com/brunoga/robomaster/sdk/modules/robot"

const (
	setRobotModeRequestSize = 1
)

// SetRobotModeRequest is the request to set the robot mode.
type SetRobotModeRequest struct {
	*baseRequest
}

var _ Request = (*SetRobotModeRequest)(nil)

// NewSetRobotModeRequest returns a new SetRobotModeRequest.
func NewSetRobotModeRequest() *SetRobotModeRequest {
	return &SetRobotModeRequest{
		baseRequest: newBaseRequest(
			setRobotModeSet,
			setRobotModeID,
			setRobotModeType,
			setRobotModeRequestSize,
		),
	}
}

// New implements the Command interface.
func (s *SetRobotModeRequest) New(data []byte) Command {
	r := NewSetRobotModeRequest()
	r.data = data

	return r
}

// SetMode sets the robot mode.
func (s *SetRobotModeRequest) SetMode(robotMode robot.Mode) {
	s.data[0] = byte(robotMode)
}

// Mode returns the currently set robot mode.
func (s *SetRobotModeRequest) Mode() robot.Mode {
	return robot.Mode(s.data[0])
}
