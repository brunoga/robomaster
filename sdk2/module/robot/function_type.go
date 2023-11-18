package robot

type FunctionType uint8

const (
	FunctionTypeBlooded FunctionType = iota + 2
	FunctionTypeMovementControl
	FunctionTypeGunControl
	FunctionTypeOffControlFunction
	FunctionTypeCount
)

func (f FunctionType) String() string {
	switch f {
	case FunctionTypeBlooded:
		return "Blooded"
	case FunctionTypeMovementControl:
		return "MovementControl"
	case FunctionTypeGunControl:
		return "GunControl"
	case FunctionTypeOffControlFunction:
		return "OffControlFunction"
	default:
		return "Unknown"
	}
}

func (f FunctionType) Valid() bool {
	return f >= FunctionTypeBlooded && f < FunctionTypeCount
}
