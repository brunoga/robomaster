package controller

import (
	"log/slog"
	"os"
	"testing"

	robomaster "github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/controller"
	"github.com/brunoga/robomaster/support"
	"github.com/brunoga/robomaster/support/logger"
)

var controllerModule *controller.Controller

func TestMain(m *testing.M) {
	c, err := robomaster.NewWithModules(logger.New(slog.LevelDebug), support.AnyAppID,
		module.TypeConnection|module.TypeRobot|module.TypeController)
	if err != nil {
		panic(err)
	}

	if err := c.Start(); err != nil {
		panic(err)
	}
	defer func() {
		if err := c.Stop(); err != nil {
			panic(err)
		}
	}()

	controllerModule = c.Controller()

	os.Exit(m.Run())
}
