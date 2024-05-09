package module

type Type int8

const (
	TypeConnection Type = 1 << iota
	TypeRobot
	TypeChassis
	TypeGimbal
	TypeCamera
	TypeGun
	TypeGamePad

	TypeAllButGamePad = TypeConnection | TypeRobot | TypeChassis | TypeGimbal |
		TypeCamera | TypeGun
)
