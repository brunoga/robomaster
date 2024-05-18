package main

import (
	"fmt"
	"log/slog"

	robomaster "github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
)

func main() {
	l := logger.New(slog.LevelWarn)

	c, err := robomaster.New(l, 0)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer c.Stop()

	fmt.Println(c.Robot().Devices())
}
