package main

import (
	"log/slog"

	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/scenes"
	"github.com/brunoga/unitybridge/support/logger"
)

func main() {
	l := logger.New(slog.LevelDebug)

	c, err := sdk2.New(l, 0)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}

	opts := engo.RunOptions{
		Title:         "Robomaster",
		Width:         1280,
		Height:        720,
		VSync:         true,
		ScaleOnResize: true,
		FPSLimit:      60,
	}

	engo.Run(opts, &scenes.Robomaster{
		Client: c,
	})
}
