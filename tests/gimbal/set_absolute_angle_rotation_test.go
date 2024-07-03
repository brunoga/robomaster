package gimbal

import (
	"testing"
	"time"

	"github.com/brunoga/robomaster/module/gimbal"
)

func TestSetAbsoluteAngleRotation(t *testing.T) {
	gimbalModule.ResetPosition()

	err := gimbalModule.SetAbsoluteAngleRotation(15, gimbal.AxisPitch, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting gimbal to absolute angle rotation: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = gimbalModule.SetAbsoluteAngleRotation(-15, gimbal.AxisPitch, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting gimbal to absolute angle rotation: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = gimbalModule.SetAbsoluteAngleRotation(-100, gimbal.AxisYaw, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting gimbal to absolute angle rotation: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = gimbalModule.SetAbsoluteAngleRotation(100, gimbal.AxisYaw, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting gimbal to absolute angle rotation: %v", err)
	}

	time.Sleep(2 * time.Second)

	err = gimbalModule.SetAbsoluteAngleRotation(0, gimbal.AxisPitch, 1*time.Second)
	if err != nil {
		t.Errorf("Error setting gimbal to absolute angle rotation: %v", err)
	}

	time.Sleep(2 * time.Second)
}
