package module

type Type uint8

const (
	TypeConnection Type = 1 << iota
	TypeRobot
	TypeController
	TypeChassis
	TypeGimbal
	TypeCamera
	TypeGun
	TypeGamePad

	TypeAll = TypeConnection | TypeRobot | TypeController | TypeChassis |
		TypeGimbal | TypeCamera | TypeGun | TypeGamePad
	TypeAllButGamePad = TypeConnection | TypeRobot | TypeController | TypeChassis |
		TypeGimbal | TypeCamera | TypeGun
)
