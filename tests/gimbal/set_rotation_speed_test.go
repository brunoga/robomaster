package gimbal

import (
	"testing"
	"time"
)

func TestSetRotationSpeed(t *testing.T) {
	// Make sure we are in a sane position.
	err := gimbalModule.ResetPosition()
	if err != nil {
		t.Errorf("ResetPosition() failed, got: %v", err)
	}

	// Also make sure we return to a sane position after the test.
	defer func() {
		err := gimbalModule.ResetPosition()
		if err != nil {
			panic(err)
		}
	}()

	err = gimbalModule.SetRotationSpeed(10, 10)
	if err != nil {
		t.Errorf("SetRotationSpeed() failed, got: %v", err)
	}
	defer gimbalModule.StopRotation()

	time.Sleep(5 * time.Second)
}
