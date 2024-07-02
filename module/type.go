package module

type Type uint16

const (
	TypeConnection Type = 1 << iota
	TypeRobot
	TypeController
	TypeChassis
	TypeGimbal
	TypeCamera
	TypeSDCard
	TypeGun
	TypeGamePad

	TypeAllButGamePad = TypeConnection | TypeRobot | TypeController |
		TypeChassis | TypeGimbal | TypeCamera | TypeSDCard | TypeGun
	TypeAll = TypeAllButGamePad | TypeGamePad
)
