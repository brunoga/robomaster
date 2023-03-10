package command

const (
	streamCtrlResponseSize = 1
)

func init() {
	Register(streamCtrlSet, streamCtrlID, NewStreamCtrlResponse())
}

// StreamCtrlResponse is the response to the StreamCtrlRequest.
type StreamCtrlResponse struct {
	*baseResponse
}

var _ Response = (*StreamCtrlResponse)(nil)

// NewStreamCtrlResponse returns a new StreamCtrlResponse.
func NewStreamCtrlResponse() *StreamCtrlResponse {
	return &StreamCtrlResponse{
		baseResponse: newBaseResponse(
			streamCtrlSet,
			streamCtrlID,
			streamCtrlType,
			streamCtrlResponseSize,
		),
	}
}

// New implements the Command interface.
func (s *StreamCtrlResponse) New(data []byte) Command {
	r := NewStreamCtrlResponse()
	r.data = data

	return r
}
