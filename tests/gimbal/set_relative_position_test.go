package gimbal

import (
	"testing"
	"time"
)

func TestSetRelativePosition(t *testing.T) {
	err := gimbalModule.SetRelativeAngleRotation(500, 500, 10*time.Second)
	if err != nil {
		t.Errorf("Error setting relative position: %s", err)
	}

	time.Sleep(20 * time.Second)
}
