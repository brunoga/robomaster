package types

// SDKProtocol is the SDK protocol to use to control the robot.
type SDKProtocol byte

const (
	// SDKProtocolBinary uses binary as the SDK protocol.
	SDKProtocolBinary SDKProtocol = iota
	// SDKProtocolText uses text as the SDK protocol.
	SDKProtocolText
	// SDKProtocolInvalid is the SDK protocol used when the SDK protocol is
	// invalid.
	SDKProtocolInvalid
)
