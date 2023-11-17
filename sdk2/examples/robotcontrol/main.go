package main

import (
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/scenes"
)

func main() {
	opts := engo.RunOptions{
		Title:         "Robomaster",
		Width:         1280,
		Height:        720,
		VSync:         true,
		ScaleOnResize: true,
		FPSLimit:      60,
	}

	engo.Run(opts, &scenes.Robomaster{})
}
