package chassis

import (
	"os"
	"testing"

	robomaster "github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/chassis"
	"github.com/brunoga/robomaster/module/controller"
	"github.com/brunoga/robomaster/support"
	"github.com/brunoga/robomaster/support/logger"
)

var chassisModule *chassis.Chassis

func TestMain(m *testing.M) {
	l := logger.New(logger.LevelTrace, "unity_bridge", "wrapper")

	c, err := robomaster.NewWithModules(l, support.AnyAppID,
		module.TypeConnection|module.TypeRobot|module.TypeController|module.TypeChassis)
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

	chassisModule = c.Chassis()

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
