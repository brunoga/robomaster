package chassis

import (
	"os"
	"testing"

	robomaster "github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/chassis"
	"github.com/brunoga/robomaster/module/controller"
	"github.com/brunoga/robomaster/support"
)

var chassisModule *chassis.Chassis

func TestMain(m *testing.M) {
	c, err := robomaster.NewWithModules(nil, support.AnyAppID,
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

	// Enable robot movement.
	//err = c.Robot().EnableFunction(robot.FunctionTypeMovementControl, true)
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	err := c.Robot().EnableFunction(robot.FunctionTypeMovementControl, false)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

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
