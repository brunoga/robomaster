package listener

import (
	"encoding/json"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/event"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	wrapper_mock "github.com/brunoga/robomaster/unitybridge/wrapper/mock"
)

func TestNewResultListener(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	rl := New(ub, nil, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl)
	assert.NotNil(t, rl.l) // Asserts that a nil logger creates a internal one.

	l := logger.New(slog.LevelError)
	rl2 := New(ub, l, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl2)

	// A new logger is created nased on the given one.
	assert.NotEqual(t, l, rl2.l)

	uw.AssertExpectations(t)
}

func TestResultListenerStart_AddListenerError(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	// Key is write only.
	rl := New(ub, nil, key.KeyCameraVideoTransRate, nil)
	assert.NotNil(t, rl)

	err := rl.Start()
	assert.Error(t, err)

	uw.AssertExpectations(t)
}

func TestResultListenerStart_Success(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	rl := New(ub, nil, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl)

	ev := event.NewFromTypeAndSubType(event.TypeStartListening,
		key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0))

	ev = event.NewFromTypeAndSubType(event.TypeGetAvailableValue,
		key.KeyAirLinkConnection.SubType())
	output := make([]byte, 2048)
	uw.On("SendEvent", ev.Code(), output, uint64(0)).
		Return([]byte("invalid"))

	err := rl.Start()
	assert.NoError(t, err)

	uw.AssertExpectations(t)
}

func TestResultListenerStart_Success_Immediate(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	rl := New(ub, nil, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl)

	ev := event.NewFromTypeAndSubType(event.TypeStartListening,
		key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0))

	ev = event.NewFromTypeAndSubType(event.TypeGetAvailableValue,
		key.KeyAirLinkConnection.SubType())
	output := make([]byte, 2048)
	uw.On("SendEvent", ev.Code(), output, uint64(0)).
		Return(resultToData(result.New(
			key.KeyAirLinkConnection, 0, 0, "", &value.Bool{})))

	err := rl.Start()
	assert.NoError(t, err)

	rl.m.Lock()
	b := rl.b
	rl.m.Unlock()

	b.Wait(10 * time.Second) // Large enough wait.

	rl.m.Lock()
	r := rl.r
	rl.m.Unlock()

	assert.NotNil(t, r)

	assert.Equal(t, key.KeyAirLinkConnection, r.Key())
	assert.Equal(t, uint64(0), r.Tag())
	assert.Equal(t, int64(0), r.ErrorCode())
	assert.Equal(t, "", r.ErrorDesc())
	assert.Equal(t, &value.Bool{}, r.Value())

	uw.AssertExpectations(t)
}

func TestResultListenerStart_AlreadyStarted(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	rl := New(ub, nil, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl)

	ev := event.NewFromTypeAndSubType(event.TypeStartListening,
		key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0))

	ev = event.NewFromTypeAndSubType(event.TypeGetAvailableValue,
		key.KeyAirLinkConnection.SubType())
	output := make([]byte, 2048)
	uw.On("SendEvent", ev.Code(), output, uint64(0)).
		Return([]byte("invalid"))

	err := rl.Start()
	assert.NoError(t, err)

	err = rl.Start()
	assert.Error(t, err)

	uw.AssertExpectations(t)
}

func TestResultListenerStop_NotStarted(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	rl := New(ub, nil, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl)

	err := rl.Stop()
	assert.Error(t, err)

	uw.AssertExpectations(t)
}

func TestResultListenerStop_Success(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	rl := New(ub, nil, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl)

	ev := event.NewFromTypeAndSubType(event.TypeStartListening,
		key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0))

	ev = event.NewFromTypeAndSubType(event.TypeGetAvailableValue,
		key.KeyAirLinkConnection.SubType())
	output := make([]byte, 2048)
	uw.On("SendEvent", ev.Code(), output, uint64(0)).
		Return([]byte("invalid"))

	err := rl.Start()
	assert.NoError(t, err)

	ev = event.NewFromTypeAndSubType(event.TypeStopListening,
		key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0))

	err = rl.Stop()
	assert.NoError(t, err)

	uw.AssertExpectations(t)
}

func TestWaitForNewResult_NotStarted(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	rl := New(ub, nil, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl)

	r := rl.WaitForNewResult(0)
	assert.Nil(t, r)

	uw.AssertExpectations(t)
}

func TestWaitForNewResult_Success(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	rl := New(ub, nil, key.KeyAirLinkConnection, nil)
	assert.NotNil(t, rl)

	ev := event.NewFromTypeAndSubType(event.TypeStartListening,
		key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0))

	ev = event.NewFromTypeAndSubType(event.TypeGetAvailableValue,
		key.KeyAirLinkConnection.SubType())
	output := make([]byte, 2048)
	uw.On("SendEvent", ev.Code(), output, uint64(0)).
		Return(resultToData(result.New(
			key.KeyAirLinkConnection, 0, 0, "", &value.Bool{})))

	err := rl.Start()
	assert.NoError(t, err)

	r := rl.WaitForNewResult(1 * time.Millisecond)

	assert.NotNil(t, r)
	assert.Equal(t, key.KeyAirLinkConnection, r.Key())
	assert.Equal(t, uint64(0), r.Tag())
	assert.Equal(t, int64(0), r.ErrorCode())
	assert.Equal(t, "", r.ErrorDesc())
	assert.Equal(t, &value.Bool{}, r.Value())

	uw.AssertExpectations(t)
}

func TestWaitForAnyResult_AvailableResult(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	var rl *Listener
	rl = New(ub, nil, key.KeyAirLinkConnection, func(r *result.Result) {
		rl.m.Lock()
		rl.r = r
		rl.m.Unlock()
	})

	assert.NotNil(t, rl)

	ev := event.NewFromTypeAndSubType(event.TypeStartListening,
		key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev.Code(), []byte(nil), uint64(0))

	ev = event.NewFromTypeAndSubType(event.TypeGetAvailableValue,
		key.KeyAirLinkConnection.SubType())
	output := make([]byte, 2048)
	uw.On("SendEvent", ev.Code(), output, uint64(0)).
		Return(resultToData(result.New(
			key.KeyAirLinkConnection, 0, 0, "", &value.Bool{Value: true})))

	err := rl.Start()
	assert.NoError(t, err)

	r := rl.WaitForAnyResult(1 * time.Second)
	assert.NotNil(t, r)
	assert.Equal(t, key.KeyAirLinkConnection, r.Key())
	assert.Equal(t, uint64(0), r.Tag())
	assert.Equal(t, int64(0), r.ErrorCode())
	assert.Equal(t, "", r.ErrorDesc())
	assert.Equal(t, &value.Bool{Value: true}, r.Value())

	uw.AssertExpectations(t)
}

func TestWaitForAnyResult_NoAvailableResult(t *testing.T) {
	uw, ub := setupUnityBridge(t)
	defer cleanupUnityBridge(t, uw, ub)

	var rl *Listener
	rl = New(ub, nil, key.KeyAirLinkConnection, func(r *result.Result) {
		rl.m.Lock()
		rl.r = r
		rl.m.Unlock()
	})
	assert.NotNil(t, rl)

	ev1 := event.NewFromTypeAndSubType(event.TypeStartListening,
		key.KeyAirLinkConnection.SubType())
	uw.On("SendEvent", ev1.Code(), []byte(nil), uint64(0))

	ev2 := event.NewFromTypeAndSubType(event.TypeGetAvailableValue,
		key.KeyAirLinkConnection.SubType())
	output := make([]byte, 2048)
	uw.On("SendEvent", ev2.Code(), output, uint64(0)).
		Return(resultToData(nil))

	err := rl.Start()

	assert.NoError(t, err)

	uw.GenerateEvent(ev1.Code(), resultToData(result.New(key.KeyAirLinkConnection, 1, 0,
		"", &value.Bool{Value: true})), 0)

	r := rl.WaitForAnyResult(5 * time.Second)
	assert.NotNil(t, r)
	assert.Equal(t, key.KeyAirLinkConnection, r.Key())
	assert.Equal(t, uint64(1), r.Tag())
	assert.Equal(t, int64(0), r.ErrorCode())
	assert.Equal(t, "", r.ErrorDesc())
	assert.Equal(t, &value.Bool{Value: true}, r.Value())

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

func setupUnityBridge(t *testing.T) (*wrapper_mock.UnityBridge, unitybridge.UnityBridge) {
	uw := wrapper_mock.NewUnityBridgeWrapper()
	ub := unitybridge.Get(uw, false, nil)

	uw.On("Create", "Robomaster", false, "")
	uw.On("Initialize").Return(true)
	for _, typ := range event.AllTypes() {
		ev := event.NewFromType(typ)
		uw.On("SetEventCallback", ev.Code(), mock.AnythingOfType("callback.Callback"))
	}

	err := ub.Start()
	assert.NoError(t, err)

	uw.AssertExpectations(t)

	uw.ExpectedCalls = nil

	return uw, ub
}

func cleanupUnityBridge(t *testing.T, uw *wrapper_mock.UnityBridge, ub unitybridge.UnityBridge) {
	for _, typ := range event.AllTypes() {
		ev := event.NewFromType(typ)
		uw.On("SetEventCallback", ev.Code(), isNilCallback())
	}

	uw.On("Uninitialize")
	uw.On("Destroy")

	ub.Stop()

	uw.AssertExpectations(t)
}

func isNilCallback() interface{} {
	return mock.MatchedBy(func(cb interface{}) bool {
		return reflect.ValueOf(cb).IsNil()
	})
}
