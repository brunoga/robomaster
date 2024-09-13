package gamepad

import (
	"os"
	"testing"

	"github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/gamepad"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/support"
	"github.com/brunoga/robomaster/support/logger"
)

var gamepadModule *gamepad.GamePad
var robotModule *robot.Robot

func TestMain(m *testing.M) {
	c, err := robomaster.NewWithModules(logger.New(logger.LevelTrace), support.AnyAppID,
		module.TypeConnection|module.TypeRobot|module.TypeGamePad)
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

	gamepadModule = c.GamePad()
	robotModule = c.Robot()

	os.Exit(m.Run())
}
