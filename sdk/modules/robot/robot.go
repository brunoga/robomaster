package robot

type Robot interface {
	GetProductVersion() (string, error)
	SetMotionMode(motionMode MotionMode) error
}
