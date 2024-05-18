package mock

import (
	"github.com/brunoga/robomaster/unitybridge/wrapper/callback"
	"github.com/stretchr/testify/mock"

	internal_callback "github.com/brunoga/robomaster/unitybridge/wrapper/internal/callback"
)

type UnityBridge struct {
	mock.Mock

	cm *internal_callback.Manager
}

func NewUnityBridgeWrapper() *UnityBridge {
	return &UnityBridge{
		cm: internal_callback.NewManager(nil),
	}
}

func (w *UnityBridge) Create(name string, debuggable bool,
	logPath string) {
	w.Called(name, debuggable, logPath)
}

func (w *UnityBridge) Initialize() bool {
	args := w.Called()
	return args.Bool(0)
}

func (w *UnityBridge) SetEventCallback(eventTypeCode uint64,
	c callback.Callback) {
	w.Called(eventTypeCode, c)

	w.cm.Set(eventTypeCode, c)
}

func (w *UnityBridge) SendEvent(eventCode uint64, output []byte,
	tag uint64) {
	args := w.Called(eventCode, output, tag)
	if len(output) > 0 && args.Get(0) != nil {
		copy(output, args.Get(0).([]byte))
	}
}

func (w *UnityBridge) SendEventWithString(eventCode uint64, data string,
	tag uint64) {
	w.Called(eventCode, data, tag)
}

func (w *UnityBridge) SendEventWithNumber(eventCode uint64, data,
	tag uint64) {
	w.Called(eventCode, data, tag)
}

func (w *UnityBridge) GetSecurityKeyByKeyChainIndex(index int) string {
	args := w.Called(index)
	return args.String(0)
}

func (w *UnityBridge) Uninitialize() {
	w.Called()
}

func (w *UnityBridge) Destroy() {
	w.Called()
}

// GenerateEvent generates an event with the given event code, data and tag.
// any callbacks registered for the same event type for the given event code
// will be called with the given data and tag.
func (w *UnityBridge) GenerateEvent(eventCode uint64, data []byte,
	tag uint64) error {
	return w.cm.Run(eventCode, data, tag)
}
