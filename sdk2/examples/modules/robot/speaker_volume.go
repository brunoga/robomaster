package main

import (
	"fmt"
	"log/slog"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/unitybridge/support/logger"
)

func main() {
	l := logger.New(slog.LevelError)

	c, err := sdk2.New(l, 0)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer c.Stop()

	// Get the robot module.
	r := c.Robot()

	// Get current speaker volume.
	volume, err := r.SpeakerVolume()
	if err != nil {
		panic(err)
	}

	fmt.Println("Original speaker volume:", volume)

	// Reset speaker volume on exit.
	defer func(v uint8) {
		err := r.SetSpeakerVolume(v)
		if err != nil {
			panic(err)
		}

		fmt.Println("Speaker volume reset to:", v)
	}(volume)

	// Set speaker volume to 50.
	err = r.SetSpeakerVolume(volume / 2)
	if err != nil {
		panic(err)
	}

	// Get current speaker volume.
	volume, err = r.SpeakerVolume()
	if err != nil {
		panic(err)
	}

	fmt.Println("New speaker volume:", volume)
}
