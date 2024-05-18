package controller

// Mode is the type of control mode for the robot.
type Mode uint8

const (
	// ModeFPV enables the dual stick controller type.
	ModeFPV Mode = iota
	// ModeSDK enables the SDK controller type.
	ModeSDK
)

func (cm Mode) String() string {
	switch cm {
	case ModeFPV:
		return "FPV"
	case ModeSDK:
		return "SDK"
	default:
		return "Unknown"
	}
}

func (cm Mode) Valid() bool {
	return cm <= ModeSDK
}
