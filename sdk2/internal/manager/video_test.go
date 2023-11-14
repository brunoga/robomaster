package manager

import (
	"log/slog"
	"testing"
	"time"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/wrapper"
)

func TestVideo(t *testing.T) {
	l := logger.New(slog.LevelDebug)

	ub := unitybridge.Get(wrapper.Get(l), true, l)

	cm, err := NewConnection(ub, l, 0)
	if err != nil {
		t.Fatal(err)
	}

	vm, err := NewVideo(ub, l)
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

	err = vm.StartRecordingToSDCard()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	err = vm.StopRecordingToSDCard()
	if err != nil {
		t.Fatal(err)
	}
}
