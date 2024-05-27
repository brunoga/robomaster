package gimbal

import (
	"log/slog"
	"os"
	"testing"

	"github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/chassis"
	"github.com/brunoga/robomaster/module/gimbal"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/unitybridge/support"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
)

var gimbalModule *gimbal.Gimbal

func TestMain(m *testing.M) {
	c, err := robomaster.NewWithModules(logger.New(slog.LevelDebug), support.AnyAppID,
		module.TypeConnection|module.TypeRobot|module.TypeChassis|module.TypeGimbal)
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

	gimbalModule.SetControlMode(gimbal.ControlMode1)

	// Enable robot movement.
	err = c.Robot().EnableFunction(robot.FunctionTypeMovementControl, true)
	if err != nil {
		panic(err)
	}

	err = c.Chassis().SetMode(chassis.ModeAngularVelocity)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
