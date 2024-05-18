package chassis

import (
	"testing"
	"time"

	"github.com/brunoga/robomaster/module/chassis"
	"github.com/brunoga/robomaster/module/chassis/controller"
)

func TestSetSpeed(t *testing.T) {
	err := chassisModule.SetControllerMode(controller.ModeSDK)
	if err != nil {
		panic(err)
	}
	defer chassisModule.SetControllerMode(controller.ModeFPV)

	setSpeed(1*time.Second, 0.2, 0.0, 0.0)
	setSpeed(1*time.Second, 0.0, 0.2, 0.0)
	setSpeed(1*time.Second, -0.2, 0.0, 0.0)
	setSpeed(1*time.Second, 0.0, -0.2, 0.0)
	setSpeed(1*time.Second, 0.0, 0.0, 360.0)
}

func setSpeed(d time.Duration, x, y, z float64) {
	err := chassisModule.SetSpeed(chassis.ModeAngularVelocity, x, y, z)
	if err != nil {
		panic(err)
	}

	time.Sleep(d)

	err = chassisModule.StopMovement(chassis.ModeAngularVelocity)
	if err != nil {
		panic(err)
	}
}
