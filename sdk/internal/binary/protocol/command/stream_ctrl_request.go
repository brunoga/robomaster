package command

const (
	streamCtrlRequestSize = 3
)

func init() {
	Register(streamCtrlSet, streamCtrlID, NewStreamCtrlRequest())
}

// StreamCtrlRequest is the command used to control the audio and video streams.
type StreamCtrlRequest struct {
	*baseRequest
}

var _ Request = (*StreamCtrlRequest)(nil)

// NewStreamCtrlRequest returns a new StreamCtrlRequest.
func NewStreamCtrlRequest() *StreamCtrlRequest {
	return &StreamCtrlRequest{
		baseRequest: newBaseRequest(
			streamCtrlSet,
			streamCtrlID,
			streamCtrlType,
			streamCtrlRequestSize,
		),
	}
}

// New implements the Command interface.
func (s *StreamCtrlRequest) New(data []byte) Command {
	r := NewStreamCtrlRequest()
	r.data = data

	return r
}

// SetControl sets the control byte.
//
// 1 means SDK control mode.
// 2 means video control mode.
// 3 means audio and video control mode.
func (s *StreamCtrlRequest) SetControl(control byte) {
	s.data[0] = control
}

// Control returns the control byte.
func (s *StreamCtrlRequest) Control() byte {
	return s.data[0]
}

// SetConnectionType sets the connection type.
//
// 0 means WiFi.
// 1 means USB (RNDIS).
func (s *StreamCtrlRequest) SetConnectionType(connectionType byte) {
	s.data[1] = ((connectionType << 4) | (s.data[1] & 0b00000001))
}

func (s *StreamCtrlRequest) ConnectionType() byte {
	return s.data[1]
}

// SetState sets the video streaming state.
//
// 0 means stop (exit SDK mode).
// 1 means start (enter SDK mode).
func (s *StreamCtrlRequest) SetState(state byte) {
	s.data[1] = (state & 0b00000001) | (s.data[1] & 0b11110000)
}

func (s *StreamCtrlRequest) SetResolution(resolution byte) {
	s.data[2] = resolution
}

func (s *StreamCtrlRequest) Resolution() byte {
	return s.data[2]
}
