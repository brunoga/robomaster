package command

const (
	sdkHeartBeatResponseSize = 1
)

func init() {
	Register(sdkHeartBeatSet, sdkHeartBeatID, NewSDKHeartBeatResponse())
}

// SDKHeartBeatResponse is the response to SDKHeartBeatRequest.
type SDKHeartBeatResponse struct {
	*baseResponse
}

var _ Response = (*SDKHeartBeatResponse)(nil)

// NewSDKHeartBeatResponse creates a new SDKHeartBeatResponse.
func NewSDKHeartBeatResponse() *SDKHeartBeatResponse {
	return &SDKHeartBeatResponse{
		baseResponse: newBaseResponse(
			sdkHeartBeatSet,
			sdkHeartBeatID,
			sdkHeartBeatType,
			sdkHeartBeatResponseSize,
		),
	}
}

// New implements the Command interface.
func (s *SDKHeartBeatResponse) New(data []byte) Command {
	r := NewSDKHeartBeatResponse()
	r.data = data

	return r
}
