package gimbal

import (
	"testing"
	"time"
)

func TestSetRotationSpeed(t *testing.T) {
	err := gimbalModule.SetRotationSpeed(10, 10)
	if err != nil {
		t.Errorf("SetRotationSpeed() failed, got: %v", err)
	}
	defer gimbalModule.StopRotation()

	time.Sleep(5 * time.Second)
}
