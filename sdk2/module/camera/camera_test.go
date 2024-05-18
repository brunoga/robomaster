package camera

import (
	"log/slog"
	"testing"
	"time"

	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
	"github.com/brunoga/robomaster/unitybridge/wrapper"
)

func TestVideo(t *testing.T) {
	l := logger.New(slog.LevelDebug)

	ub := unitybridge.Get(wrapper.Get(l), true, l)

	cm, err := connection.New(ub, l, 0, connection.TypeRouter)
	if err != nil {
		t.Fatal(err)
	}

	vm, err := New(ub, l, cm)
	if err != nil {
		t.Fatal(err)
	}

	err = ub.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer ub.Stop()

	err = cm.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer cm.Stop()

	err = vm.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer vm.Stop()

	time.Sleep(2 * time.Second)

	if err := vm.SetVideoQuality(6.0); err != nil {
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := vm.SetVideoFormat(VideoFormat(VideoFormat1080p_60 + 100)); err != nil {
		if err != nil {
			t.Fatal(err)
		}
	}

	if value, err := vm.VideoFormat(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(value)
	}

	time.Sleep(10 * time.Second)

	//	err = vm.StartRecordingVideoToSDCard()
	//	if err != nil {
	//		t.Fatal(err)
	//	}

	//	time.Sleep(10 * time.Second)

	// err = vm.StopRecordingVideoToSDCard()
	//
	//	if err != nil {
	//		t.Fatal(err)
	//	}
}
