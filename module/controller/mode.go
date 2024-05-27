package controller

import "fmt"

// Mode is the type of control mode for the robot.
type Mode uint8

const (
	// ModeFPV enables the dual stick controller type.
	ModeFPV Mode = iota
	// ModeSDK enables the SDK controller type.
	ModeSDK
	// modeCount is the number of modes. Intentionaly not exported.
	modeCount
)

func (m Mode) String() string {
	switch m {
	case ModeFPV:
		return "FPV"
	case ModeSDK:
		return "SDK"
	default:
		return fmt.Sprintf("Unknown(%d)", m)
	}
}

func (m Mode) Valid() bool {
	return m <= modeCount
}
