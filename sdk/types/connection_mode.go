package types

// ConnectionMode is the connection mode used to connect to the robot.
type ConnectionMode byte

const (
	// ConnectionModeAP is the connection mode used when connecting to the
	// robot through its access point.
	ConnectionModeAP ConnectionMode = iota
	// ConnectionModeInfrastructure is the connection mode used when connecting
	// to the robot through a wireless router.
	ConnectionModeInfrastructure
	// ConnectionModeUSB is the connection mode used when connecting to the
	// robot through its USB port (RNDIS).
	ConnectionModeUSB
	// ConnectionModeInvalid is the connection mode used when the connection
	// mode is invalid.
	ConnectionModeInvalid
)
