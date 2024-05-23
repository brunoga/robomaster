package chassis

import (
	"testing"
	"time"

	"github.com/brunoga/robomaster/module/chassis"
)

func TestSetPosition(t *testing.T) {
	err := chassisModule.SetPosition(chassis.ModeAngularVelocity, 1, 0, 0)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
}
