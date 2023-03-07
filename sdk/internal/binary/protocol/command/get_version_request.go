package command

const (
	getVersionRequestSize = 0
)

func init() {
	Register(getVersionSet, getVersionID, NewGetVersionRequest())
}

// GetVersionRequest is a command to get the version of the SDK that exists in
// the robot.
type GetVersionRequest struct {
	*baseRequest
}

var _ Request = (*GetVersionRequest)(nil)

// NewGetVersionRequest creates a new GetVersionRequest.
func NewGetVersionRequest() *GetVersionRequest {
	return &GetVersionRequest{
		baseRequest: newBaseRequest(
			getVersionSet,
			getVersionID,
			getVersionType,
			getVersionRequestSize,
		),
	}
}

// New implements the Command interface.
func (g *GetVersionRequest) New() Command {
	return NewGetVersionRequest()
}
