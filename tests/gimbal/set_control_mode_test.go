package gimbal

import (
	"testing"

	"github.com/brunoga/robomaster/module/gimbal"
)

func TestSetControlMode(t *testing.T) {
	cm := gimbalModule.ControlMode()

	t.Logf("Current control mode: %d", cm)

	err := gimbalModule.SetControlMode(gimbal.ControlMode2)
	if err != nil {
		t.Errorf("Error setting control mode: %s", err)
	}

	cm = gimbalModule.ControlMode()
	if err != nil {
		t.Errorf("Error getting control mode: %s", err)
	}

	t.Logf("Current control mode: %d", cm)
}
