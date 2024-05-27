package gimbal

import (
	"fmt"
	"testing"
)

func TestSetWorkMode(t *testing.T) {
	wm, err := gimbalModule.WorkMode()
	if err != nil {
		t.Fatalf("Error getting work mode: %s", err)
	}
	defer func(wm uint64) {
		if err := gimbalModule.SetWorkMode(wm); err != nil {
			t.Fatalf("Error setting work mode: %s", err)
		}
	}(wm)

	fmt.Printf("Work Mode: %d\n", wm)

	wm2 := uint64(0)

	if err := gimbalModule.SetWorkMode(wm2); err != nil {
		t.Fatalf("Error setting work mode: %s", err)
	}

	wm, err = gimbalModule.WorkMode()
	if err != nil {
		t.Fatalf("Error getting work mode: %s", err)
	}

	if wm != wm2 {
		t.Fatalf("Invalid work mode: %d", wm2)
	}
}
