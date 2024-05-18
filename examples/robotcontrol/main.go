package main

import (
	"flag"
	"log/slog"

	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/examples/robotcontrol/scenes"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
)

var (
	fullscreen = flag.Bool("fullscreen", false, "Run in fullscreen mode.")
)

func main() {
	flag.Parse()

	l := logger.New(slog.LevelDebug)

	c, err := sdk2.NewWithModules(l, 0, module.TypeAllButGamePad)
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
		Fullscreen:    *fullscreen,
	}

	engo.Run(opts, &scenes.Robomaster{
		Client: c,
	})
}
