package robot

import (
	"log/slog"
	"os"
	"testing"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/robot"
	"github.com/brunoga/robomaster/sdk2/unitybridge/support"
	"github.com/brunoga/robomaster/sdk2/unitybridge/support/logger"
)

var robotModule *robot.Robot

func TestMain(m *testing.M) {
	c, err := sdk2.NewWithModules(logger.New(slog.LevelDebug), support.AnyAppID,
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
