package modules

// RobotMotionMode represents the motion mode for a robot.
type MotionMode int

// Available robot motion modes.
const (
	MotionModeChassisLead MotionMode = iota // gimbal follows chassis
	MotionModeGimbalLead                    // chassis follow gimbal
	MotionModeFree                          // chassis and gimbal move independently
	MotionModeInvalid                       // invalid mode
)

type Robot interface {
	GetSDKVersion() (string, error)
	SetMotionMode(motionMode MotionMode) error
}
