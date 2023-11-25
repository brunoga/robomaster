package controller

// Mode is the type of control mode for the robot. Apparently, only
// ControlModeDefault is ever used but we know for a fact that there is a
// second mode.
//
// TODO(bga): Investigate.
type Mode uint8

const (
	ModeFPV Mode = iota
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
