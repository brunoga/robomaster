package command

const (
	setSDKModeResponseSize = 1
)

func init() {
	Register(setSDKModeSet, setSDKModeID, NewSetSDKModeResponse())
}

// SetSDKModeResponse is the response to SetSDKModeRequest.
type SetSDKModeResponse struct {
	*baseResponse
}

var _ Response = (*SetSDKModeResponse)(nil)

// NewSetSDKModeResponse creates a new SetSDKModeResponse.
func NewSetSDKModeResponse() *SetSDKModeResponse {
	return &SetSDKModeResponse{
		baseResponse: newBaseResponse(
			setSDKModeSet,
			setSDKModeID,
			setSDKModeType,
			setSDKModeResponseSize,
		),
	}
}

// New omplements the command interface.
func (s *SetSDKModeResponse) New(data []byte) Command {
	r := NewSetSDKModeResponse()
	r.data = data

	return r
}
