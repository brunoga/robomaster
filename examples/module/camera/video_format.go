package main

import (
	"fmt"
	"log/slog"

	robomaster "github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module/camera"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
)

func main() {
	l := logger.New(slog.LevelError)

	c, err := robomaster.New(l, 0)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer c.Stop()

	// Get the camera module.
	cm := c.Camera()

	// Get current video format.
	format, err := cm.VideoFormat()
	if err != nil {
		panic(err)
	}

	fmt.Println("Original video format:", format)

	// Reset video format on exit.
	defer func(f camera.VideoFormat) {
		err := cm.SetVideoFormat(f)
		if err != nil {
			panic(err)
		}

		fmt.Println("Video format reset to:", f)
	}(format)

	// Set video format to 1080p.
	err = cm.SetVideoFormat(camera.VideoFormat1080p_60)
	if err != nil {
		panic(err)
	}

	// Get current video format.
	format, err = cm.VideoFormat()
	if err != nil {
		panic(err)
	}

	fmt.Println("New video format:", format)
}
