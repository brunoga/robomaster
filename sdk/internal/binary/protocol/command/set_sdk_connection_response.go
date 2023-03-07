package command

import (
	"net"
)

const (
	setSDKConnectionResponseSize = 6
)

func init() {
	Register(setSDKConnectionSet, setSDKConnectionID, NewSetSDKConnectionResponse())
}

// SetSDKConnectionResponse is the response to SetSDKConnectionRequest.
type SetSDKConnectionResponse struct {
	*baseResponse
}

var _ Response = (*SetSDKConnectionResponse)(nil)

// NewSetSDKConnectionResponse creates a new SetSDKConnectionResponse.
func NewSetSDKConnectionResponse() *SetSDKConnectionResponse {
	return &SetSDKConnectionResponse{
		baseResponse: newBaseResponse(
			setSDKConnectionSet,
			setSDKConnectionID,
			setSDKConnectionType,
			setSDKConnectionResponseSize,
		),
	}
}

// New implements the Command interface.
func (s *SetSDKConnectionResponse) New(data []byte) Command {
	r := NewSetSDKConnectionResponse()
	r.data = data

	return r
}

// SetConfigIP sets the config IP.
func (s *SetSDKConnectionResponse) SetConfigIP(ip net.IP) {
	copy(s.data[2:6], ip)
}

// ConfigIP returns the config IP.This is the IP the robot saw when it got the
// SetSDKConnectionRequest.
func (s *SetSDKConnectionResponse) ConfigIP() net.IP {
	return net.IP(s.data[2:6])
}
