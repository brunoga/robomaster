package types

// ConnectionProtocol is the connection protocol to use when connecting to the robot.
type ConnectionProtocol byte

const (
	// ConnectionProtocolUDP uses UDP as the connection protocol.
	ConnectionProtocolUDP ConnectionProtocol = iota
	// ConnectionProtocolTCP uses TCP as the connection protocol.
	ConnectionProtocolTCP
	// ConnectionProtocolInvalid is the connection protocol used when the connection
	// protocol is invalid.
	ConnectionProtocolInvalid
)
