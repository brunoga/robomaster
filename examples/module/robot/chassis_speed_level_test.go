package robot

import (
	"testing"
)

func TestChassisSpeedLevel(t *testing.T) {
	originalSpeedLevel, err := robotModule.ChassisSpeedLevel()
	if err != nil {
		t.Fatalf("Failed to get original speed level: %v", err)
	}
	defer func() {
		err := robotModule.SetChassisSpeedLevel(originalSpeedLevel)
		if err != nil {
			t.Fatalf("Failed to restore original speed level: %v", err)
		}
	}()

	err = robotModule.SetChassisSpeedLevel(1)
	if err != nil {
		t.Fatalf("Failed to set speed level to 1: %v", err)
	}

	speedLevel, err := robotModule.ChassisSpeedLevel()
	if err != nil {
		t.Fatalf("Failed to get speed level: %v", err)
	}
	if speedLevel != 1 {
		t.Fatalf("Speed level is not 0: %v", speedLevel)
	}

	err = robotModule.SetChassisSpeedLevel(4)
	if err != nil {
		t.Fatalf("Failed to set speed level to 4: %v", err)
	}

	speedLevel, err = robotModule.ChassisSpeedLevel()
	if err != nil {
		t.Fatalf("Failed to get speed level: %v", err)
	}

	if speedLevel != 4 {
		t.Fatalf("Speed level is not 4: %v", speedLevel)
	}
}
