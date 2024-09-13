package internal

import (
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/event"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
	"github.com/stretchr/testify/mock"

	wrapper_mock "github.com/brunoga/robomaster/unitybridge/wrapper/mock"
)

func TestNewBaseModule(t *testing.T) {
	uw := wrapper_mock.NewUnityBridgeWrapper()

	l := logger.New(logger.LevelTrace)
	ub := unitybridge.Get(uw, false, l)

	name := "test"
	k := key.KeyAirLinkConnection

	bm := NewBaseModule(ub, l, name, k)

	if bm == nil {
		t.Errorf("NewBaseModule() returned nil.")
	}
}

func TestBaseModule_String(t *testing.T) {
	bm, _ := newBaseModule(t)

	if bm.String() != "test" {
		t.Errorf("BaseModule.String() returned unexpected value.")
	}
}

func TestBaseModule_Start(t *testing.T) {
	bm, uw := newBaseModule(t)

	startBaseModule(t, bm, uw)
}

func TestBaseModule_Stop(t *testing.T) {
	bm, uw := newBaseModule(t)

	startBaseModule(t, bm, uw)

	stopBaseModule(t, bm, uw)
}

func TestBaseModule_Connected(t *testing.T) {
	bm, uw := newBaseModule(t)

	startBaseModule(t, bm, uw)

	if bm.Connected() {
		t.Errorf("BaseModule.Connected() returned true before connection.")
	}

	ev := event.NewFromTypeAndSubType(event.TypeStartListening, key.KeyAirLinkConnection.SubType())
	r := result.New(key.KeyAirLinkConnection, 0, 0, "", &value.Bool{Value: true})
	uw.GenerateEvent(ev.Code(), resultToData(r), 0)

	// Give some time for the callback to run.
	time.Sleep(10 * time.Millisecond)

	if !bm.Connected() {
		t.Errorf("BaseModule.Connected() returned false after connection.")
	}

	stopBaseModule(t, bm, uw)

	if bm.Connected() {
		t.Errorf("BaseModule.Connected() returned true after stopping.")
	}
}

func TestBaseModule_WaitForConnectionStatus(t *testing.T) {
	bm, uw := newBaseModule(t)

	startBaseModule(t, bm, uw)
	defer stopBaseModule(t, bm, uw)

	bm.WaitForConnectionStatus(false)

	go func() {
		ev := event.NewFromTypeAndSubType(event.TypeStartListening, key.KeyAirLinkConnection.SubType())
		r := result.New(key.KeyAirLinkConnection, 0, 0, "", &value.Bool{Value: true})
		uw.GenerateEvent(ev.Code(), resultToData(r), 0)
	}()

	bm.WaitForConnectionStatus(true)
}

func newBaseModule(t *testing.T) (*BaseModule, *wrapper_mock.UnityBridge) {
	t.Helper()

	uw := wrapper_mock.NewUnityBridgeWrapper()

	l := logger.New(slog.LevelError)
	ub := unitybridge.Get(uw, false, l)

	name := "test"
	k := key.KeyAirLinkConnection

	return NewBaseModule(ub, l, name, k), uw
}

func startBaseModule(t *testing.T, bm *BaseModule, uw *wrapper_mock.UnityBridge) {
	// Start unity bridge.
	uw.On("Create", "Robomaster", false, "")
	uw.On("Initialize").Return(true)
	for _, typ := range event.AllTypes() {
		ev := event.NewFromType(typ)
		uw.On("SetEventCallback", ev.Code(), mock.AnythingOfType("callback.Callback"))
	}

	if err := bm.ub.Start(); err != nil {
		t.Fatalf("UnityBridge.Start() failed: %v.", err)
	}

	uw.AssertExpectations(t)

	// First event: StartListening.
	ev := event.NewFromTypeAndSubType(event.TypeStartListening, key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0)).
		Return(nil)

	// Second event: GetAvailableValue.
	ev = event.NewFromTypeAndSubType(event.TypeGetAvailableValue, key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), make([]byte, 2048), uint64(0)).
		Return(nil)

	err := bm.Start()
	if err != nil {
		t.Fatalf("BaseModule.Start() failed: %v.", err)
	}

	if !uw.AssertExpectations(t) {
		t.Fatalf("Start expectations not met.")
	}
}

func stopBaseModule(t *testing.T, bm *BaseModule, uw *wrapper_mock.UnityBridge) {
	// Event: StopListening.
	ev := event.NewFromTypeAndSubType(event.TypeStopListening, key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0)).
		Return(nil)

	err := bm.Stop()
	if err != nil {
		t.Fatalf("BaseModule.Stop() failed: %v.", err)
	}

	if !uw.AssertExpectations(t) {
		t.Fatalf("Stop expectations not met.")
	}

	// Stop unity bridge.
	uw.On("Uninitialize")
	uw.On("Destroy")

	if err := bm.ub.Stop(); err != nil {
		t.Fatalf("UnityBridge.Stop() failed: %v.", err)
	}

	uw.AssertExpectations(t)
}

func resultToData(r *result.Result) []byte {
	if r == nil {
		return nil
	}

	data, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return data
}
