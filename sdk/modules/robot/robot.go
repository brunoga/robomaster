package robot

type Robot interface {
	GetSDKVersion() (string, error)
	SetMotionMode(motionMode MotionMode) error
}
