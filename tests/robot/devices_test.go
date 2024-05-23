package robot

import (
	"fmt"
	"testing"
)

func TestDevices(t *testing.T) {
	devices := robotModule.Devices()
	if len(devices) == 0 {
		t.Fatal("No working devices found.")
	} else {
		fmt.Println("Found devices:", devices)
	}
}
