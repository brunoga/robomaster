package robot

import (
	"testing"
)

func TestSpeakerVolume(t *testing.T) {
	originalVolume, err := robotModule.SpeakerVolume()
	if err != nil {
		t.Fatalf("Failed to get original volume: %v", err)
	}
	defer func() {
		err := robotModule.SetSpeakerVolume(originalVolume)
		if err != nil {
			t.Fatalf("Failed to restore original volume: %v", err)
		}
	}()

	err = robotModule.SetSpeakerVolume(0)
	if err != nil {
		t.Fatalf("Failed to set volume to 0: %v", err)
	}

	volume, err := robotModule.SpeakerVolume()
	if err != nil {
		t.Fatalf("Failed to get volume: %v", err)
	}
	if volume != 0 {
		t.Fatalf("Volume is not 0: %v", volume)
	}

	err = robotModule.SetSpeakerVolume(100)
	if err != nil {
		t.Fatalf("Failed to set volume to 100: %v", err)
	}

	volume, err = robotModule.SpeakerVolume()
	if err != nil {
		t.Fatalf("Failed to get volume: %v", err)
	}
	if volume != 100 {
		t.Fatalf("Volume is not 100: %v", volume)
	}
}
