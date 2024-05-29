package gimbal

import (
	"testing"
	"time"
)

func TestSetAbsoluteAngleRotation(t *testing.T) {
	gimbalModule.ResetPosition()

	time.Sleep(5 * time.Second)

	// Set gimbal to absolute angle rotation.
	err := gimbalModule.SetAbsoluteAngleRotation(0, 24, 100)
	if err != nil {
		t.Errorf("Error setting gimbal to absolute angle rotation: %v", err)
		return
	}

	// Wait for gimbal to reach position.
	time.Sleep(5 * time.Second)
}
