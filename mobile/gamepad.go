package mobile

import "github.com/brunoga/robomaster/module/gamepad"

// GamePad allows controlling the DJI Robomaster gamepad accessory.
type GamePad struct {
	g *gamepad.GamePad
}

// C1Pressed returns whether the C1 button is pressed.
func (g *GamePad) C1Pressed() bool {
	return g.g.C1Pressed()
}

// C2Pressed returns whether the C2 button is pressed.
func (g *GamePad) C2Pressed() bool {
	return g.g.C2Pressed()
}

// FirePressed returns whether the fire button is pressed.
func (g *GamePad) FirePressed() bool {
	return g.g.FirePressed()
}

// FnPressed returns whether the Fn button is pressed.
func (g *GamePad) FnPressed() bool {
	return g.g.FnPressed()
}
