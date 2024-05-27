package gimbal

import (
	"testing"
	"time"
)

func TestSetAbsoluteAngleRotation(t *testing.T) {

	//gimbalModule.ResetPosition()

	//time.Sleep(5 * time.Second)

	err := gimbalModule.SetWorkMode(2)
	if err != nil {
		t.Errorf("SetWorkMode() failed, got: %v", err)
	}

	// Set gimbal to absolute angle rotation.
	err = gimbalModule.SetAbsoluteAngleRotation(0, -300, 1000)
	if err != nil {
		t.Errorf("Error setting gimbal to absolute angle rotation: %v", err)
		return
	}

	// Wait for gimbal to reach position.
	time.Sleep(5 * time.Second)
}
