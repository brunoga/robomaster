package command

const (
	sdkHeartBeatRequestSize = 0
)

func init() {
	Register(sdkHeartBeatSet, sdkHeartBeatID, NewSDKHeartBeatRequest())
}

// SDKHeartBeatRequest is the command used to keep the control connection alive.
// Must be sent at leat once every 5 seconds.
type SDKHeartBeatRequest struct {
	*baseRequest
}

var _ Request = (*SDKHeartBeatRequest)(nil)

// NewSDKHeartBeatRequest returns a new SDKHeartBeatRequest.
func NewSDKHeartBeatRequest() *SDKHeartBeatRequest {
	return &SDKHeartBeatRequest{
		baseRequest: newBaseRequest(
			sdkHeartBeatSet,
			sdkHeartBeatID,
			sdkHeartBeatType,
			sdkHeartBeatRequestSize,
		),
	}
}

// New implements the Command interface.
func (s *SDKHeartBeatRequest) New(data []byte) Command {
	r := NewSDKHeartBeatRequest()
	r.data = data

	return r
}
