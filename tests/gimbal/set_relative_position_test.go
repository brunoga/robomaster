package gimbal

import (
	"testing"
	"time"

	"github.com/brunoga/robomaster/module/gimbal"
)

func TestSetRelativePosition(t *testing.T) {
	err := gimbalModule.SetRelativeAngleRotation(-15, gimbal.AxisPitch, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting relative position: %s", err)
	}

	time.Sleep(2 * time.Second)

	err = gimbalModule.SetRelativeAngleRotation(15, gimbal.AxisPitch, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting relative position: %s", err)
	}

	time.Sleep(2 * time.Second)

	err = gimbalModule.SetRelativeAngleRotation(-15, gimbal.AxisYaw, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting relative position: %s", err)
	}

	time.Sleep(2 * time.Second)
	err = gimbalModule.SetRelativeAngleRotation(15, gimbal.AxisYaw, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting relative position: %s", err)
	}

	time.Sleep(2 * time.Second)
}
