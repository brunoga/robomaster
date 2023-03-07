package command

import (
	"encoding/binary"
	"net"
)

const (
	setSDKConnectionRequestSize = 10
)

func init() {
	Register(setSDKConnectionSet, setSDKConnectionID, NewSetSDKConnectionRequest())
}

// SetSDKConnectionRequest is the command to set the SDK connection. This must
// be sent to the proxy UDP port in the robot instead and is use to set what the
// connection mode will be. The associated TCP or UDP port will then be opened
// on the robot and we can connect the control module to it.
type SetSDKConnectionRequest struct {
	*baseRequest
}

var _ Request = (*SetSDKConnectionRequest)(nil)

// NewSetSDKConnectionRequest returns a new SetSDKConnectionRequest.
func NewSetSDKConnectionRequest() *SetSDKConnectionRequest {
	return &SetSDKConnectionRequest{
		baseRequest: newBaseRequest(
			setSDKConnectionSet,
			setSDKConnectionID,
			setSDKConnectionType,
			setSDKConnectionRequestSize,
		),
	}
}

// New implements the Command interface.
func (s *SetSDKConnectionRequest) New(data []byte) Command {
	r := NewSetSDKConnectionRequest()
	r.data = data

	return r
}

// SetControl sets the control mode.
//
// TODO(bga): Figure out what this actually is. Currently we always set it to 0.
func (s *SetSDKConnectionRequest) SetControl(control byte) {
	s.data[0] = control
}

// Control returns the control mode.
func (s *SetSDKConnectionRequest) Control() byte {
	return s.data[0]
}

// SetHost sets the host byte for this host (client).
func (s *SetSDKConnectionRequest) SetHost(host byte) {
	s.data[1] = host
}

// Host returns the currently set host byte for this host (client).
func (s *SetSDKConnectionRequest) Host() byte {
	return s.data[1]
}

// SetConnection sets the connection mode.
//
// 0 = AP mode.
// 1 = Infrastructure (router) mode.
// 2 = USB RNDIS mode.
func (s *SetSDKConnectionRequest) SetConnection(connection byte) {
	s.data[2] = connection
}

// Connection returns the currently set connection mode.
func (s *SetSDKConnectionRequest) Connection() byte {
	return s.data[2]
}

// SetProtocol sets the protocol to be used.
//
// 0 = UDP.
// 1 = TCP.
func (s *SetSDKConnectionRequest) SetProtocol(protocol byte) {
	s.data[3] = protocol
}

// Protocol returns the currently set protocol.
func (s *SetSDKConnectionRequest) Protocol() byte {
	return s.data[3]
}

// SetIP sets the local IP to be used. 0.0.0.0 means any.
func (s *SetSDKConnectionRequest) SetIP(ip net.IP) {
	copy(s.data[4:8], ip)
}

// IP returns the currently set IP.
func (s *SetSDKConnectionRequest) IP() net.IP {
	return net.IP(s.data[4:8])
}

// SetPort sets the local port to be used.
func (s *SetSDKConnectionRequest) SetPort(port uint16) {
	binary.LittleEndian.PutUint16(s.data[8:10], port)
}

// Port returns the currently set port.
func (s *SetSDKConnectionRequest) Port() uint16 {
	return binary.LittleEndian.Uint16(s.data[8:10])
}
