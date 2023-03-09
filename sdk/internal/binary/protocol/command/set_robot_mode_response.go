package command

const (
	setRobotModeResponseSize = 1
)

// SetRobotModeResponse is the response to set the robot mode.
type SetRobotModeResponse struct {
	*baseResponse
}

var _ Response = (*SetRobotModeResponse)(nil)

// NewSetRobotModeResponse returns a new SetRobotModeResponse.
func NewSetRobotModeResponse() *SetRobotModeResponse {
	return &SetRobotModeResponse{
		baseResponse: newBaseResponse(
			setRobotModeSet,
			setRobotModeID,
			setRobotModeType,
			setRobotModeResponseSize,
		),
	}
}

// New implements the Command interface.
func (s *SetRobotModeResponse) New(data []byte) Command {
	r := NewSetRobotModeResponse()
	r.data = data

	return r
}
