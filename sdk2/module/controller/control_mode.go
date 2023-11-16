package controller

// ControlMode is the type of control mode for the robot. Apparently, only
// ControlModeDefault is ever used but we know for a fact that there is a
// second mode.
//
// TODO(bga): Investigate.
type ControlMode uint8

const (
	ControlModeDefault ControlMode = iota
	ControlModeOther
	ControlModeCount
)

func (cm ControlMode) String() string {
	switch cm {
	case ControlModeDefault:
		return "Default"
	case ControlModeOther:
		return "Other"
	default:
		return "Unknown"
	}
}

func (cm ControlMode) Valid() bool {
	return cm < ControlModeCount
}
