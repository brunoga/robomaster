package main

import (
	"fmt"
	"log/slog"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/sdk2/unitybridge/support/logger"
)

func main() {
	l := logger.New(slog.LevelWarn)

	c, err := sdk2.New(l, 0)
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
