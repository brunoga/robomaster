package command

const (
	setSDKModeRequestSize = 1
)

func init() {
	Register(setSDKModeSet, setSDKModeID, NewSetSDKModeRequest())
}

// SetSDKModeRequest is the command used to enable/disable SDK mode.
type SetSDKModeRequest struct {
	*baseRequest
}

var _ Request = (*SetSDKModeRequest)(nil)

// NewSetSDKModeRequest returns a new SetSDKModeRequest.
func NewSetSDKModeRequest() *SetSDKModeRequest {
	return &SetSDKModeRequest{
		baseRequest: newBaseRequest(
			setSDKModeSet,
			setSDKModeID,
			setSDKModeType,
			setSDKModeRequestSize,
		),
	}
}

// New implements the Command interface.
func (s *SetSDKModeRequest) New(data []byte) Command {
	r := NewSetSDKModeRequest()
	r.data = data

	return r
}

// SetEnable sets the enable flag, which turns the SDK mode on or off.
func (s *SetSDKModeRequest) SetEnable(enable bool) {
	if enable {
		s.data[0] = 1
	} else {
		s.data[0] = 0
	}
}

// Enable returns if the SDK mode is on or off.
func (s *SetSDKModeRequest) Enable() bool {
	return s.data[0] == 1
}
