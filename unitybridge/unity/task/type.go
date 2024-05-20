package task

type Type int8

const (
	TypeUnknown Type = iota
	TypeChassisPosition
	TypeGimbalReset
	TypeGimbalAngle
	TypeCount
)
