package gimbal

// ControlMode is the gimbal control mode.
//
// Note that from an app perspective, there are 3 control modes but it appeaars
// that the robot itself only supports 2. The 3rd mode just maps to the 1st one
// in the app.
//
// TODO(bga): Figure out what each mode means.
type ControlMode uint8

const (
	ControlMode1 ControlMode = iota
	ControlMode2
	ControlMode3
)

func (cm ControlMode) String() string {
	switch cm {
	case ControlMode1:
		return "ControlMode1"
	case ControlMode2:
		return "ControlMode2"
	case ControlMode3:
		return "ControlMode3"
	default:
		return "Unknown"
	}
}

func (cm ControlMode) Valid() bool {
	return cm >= ControlMode1 && cm <= ControlMode3
}
