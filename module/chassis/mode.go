package chassis

import "fmt"

type Mode int8

const (
	ModeYawFollow Mode = iota
	ModeTank
	ModeFPV
	ModeAngularVelocity
	ModeWayPoint
	ModeNone
)

func (m Mode) String() string {
	switch m {
	case ModeYawFollow:
		return "YawFollow"
	case ModeTank:
		return "Tank"
	case ModeFPV:
		return "FPV"
	case ModeAngularVelocity:
		return "AngularVelocity"
	case ModeWayPoint:
		return "WayPoint"
	case ModeNone:
		return "None"
	default:
		return fmt.Sprintf("Unknown(%d)", m)
	}
}

func (m Mode) Valid() bool {
	return m >= ModeYawFollow && m <= ModeNone
}
