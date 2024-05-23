package robot

import (
	"os"
	"testing"

	robomaster "github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/unitybridge/support"
)

var robotModule *robot.Robot

func TestMain(m *testing.M) {
	c, err := robomaster.NewWithModules(nil, support.AnyAppID,
		module.TypeConnection|module.TypeRobot)
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

	robotModule = c.Robot()

	os.Exit(m.Run())
}
