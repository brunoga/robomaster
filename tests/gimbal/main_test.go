package gimbal

import (
	"os"
	"testing"

	"github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/chassis"
	"github.com/brunoga/robomaster/module/controller"
	"github.com/brunoga/robomaster/module/gimbal"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/support"
	"github.com/brunoga/robomaster/support/logger"
)

var gimbalModule *gimbal.Gimbal
var chassisModule *chassis.Chassis
var robotModule *robot.Robot
var controllerModule *controller.Controller

func TestMain(m *testing.M) {
	c, err := robomaster.NewWithModules(logger.New(logger.LevelTrace), support.AnyAppID,
		module.TypeConnection|module.TypeRobot|module.TypeController|module.TypeGimbal|module.TypeChassis)
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
	chassisModule = c.Chassis()
	robotModule = c.Robot()
	controllerModule = c.Controller()

	os.Exit(m.Run())
}
