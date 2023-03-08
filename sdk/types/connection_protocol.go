package types

// ConnectionProtocol is the connection protocol to use when connecting to the robot.
type ConnectionProtocol int

const (
	// ConnectionProtocolUDP uses UDP as the connection protocol.
	ConnectionProtocolUDP ConnectionProtocol = iota
	// ConnectionProtocolTCP uses TCP as the connection protocol.
	ConnectionProtocolTCP
)
