package robot

// Mode represents the mode for a robot.
type Mode byte

// Available robot motion modes.
const (
	// ModeFree is the default motion mode. In this mode the chassis and
	// gimbal move independently.
	ModeFree Mode = iota
	// ModeGimbalLead is the gimbal lead motion mode. In this mode the
	// chassis follows the gimbal.
	ModeGimbalLead
	// ModeChassisLead is the chassis lead motion mode. In this mode the
	// gimbal follows the chassis.
	ModeChassisLead
)
