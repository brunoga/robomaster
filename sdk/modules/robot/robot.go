package robot

import "net"

// Robot is the interface implemented by all robot variations.
type Robot interface {
	// IP returns the IP associated with the Robot instance.
	IP() net.IP

	// SN returns the serial number associated with the Robot instance.
	SN() string
}
