package robot

type Robot interface {
	GetProductVersion() (string, error)
	SetMode(motionMode Mode) error
}
