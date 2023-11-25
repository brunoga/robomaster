package gun

type Type uint8

const (
	TypeBead Type = iota
	TypeInfrared
)

func (t Type) String() string {
	switch t {
	case TypeBead:
		return "Bead"
	case TypeInfrared:
		return "Infrared"
	default:
		return "Unknown"
	}
}

func (t Type) Valid() bool {
	return t <= TypeInfrared
}
