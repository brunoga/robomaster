package chassis

import (
	"log/slog"
	"os"
	"testing"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/sdk2/module/chassis"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/logger"
)

var chassisModule *chassis.Chassis

func TestMain(m *testing.M) {
	c, err := sdk2.New(logger.New(slog.LevelDebug), support.AnyAppID)
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

	os.Exit(m.Run())
}
