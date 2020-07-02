package main

import (
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk/examples/robotcontrol/scenes"
)

func main() {
	opts := engo.RunOptions{
		Title:         "Robot Control",
		Width:         1280,
		Height:        720,
		VSync:         true,
		ScaleOnResize: true,
		FPSLimit:      30,
		NotResizable:  true,
	}

	engo.Run(opts, &scenes.RobotControl{})
}
