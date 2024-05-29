package gimbal

import (
	"log/slog"
	"os"
	"testing"

	"github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/controller"
	"github.com/brunoga/robomaster/module/gimbal"
	"github.com/brunoga/robomaster/support"
	"github.com/brunoga/robomaster/support/logger"
)

var gimbalModule *gimbal.Gimbal

func TestMain(m *testing.M) {
	c, err := robomaster.NewWithModules(logger.New(slog.LevelDebug), support.AnyAppID,
		module.TypeConnection|module.TypeRobot|module.TypeController|module.TypeGimbal)
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

	gimbalModule = c.Gimbal()

	// Set controller mode to SDK for the tests here.
	err = c.Controller().SetMode(controller.ModeSDK)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := c.Controller().SetMode(controller.ModeFPV)
		if err != nil {
			panic(err)
		}
	}()

	os.Exit(m.Run())
}
