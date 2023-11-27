//go:build !(windows && amd64) && !(ios && arm64) && !(android && (arm || arm64)) && !(darwin && amd64) && !(linux && amd64)

package implementations

import (
	"fmt"
	"runtime"

	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/wrapper/callback"
)

func init() {
	panic(fmt.Sprintf("The UnityBridge library is not available for platform "+
		"%s_%s", runtime.GOOS, runtime.GOARCH))
}

var (
	UnityBridgeImpl *unsupportedUnityBridgeImpl = &unsupportedUnityBridgeImpl{}
)

type unsupportedUnityBridgeImpl struct{}

func Get(l *logger.Logger) *unsupportedUnityBridgeImpl {
	return nil
}

func (u *unsupportedUnityBridgeImpl) Create(name string, debuggable bool,
	logPath string) {
}

func (u *unsupportedUnityBridgeImpl) Initialize() bool { return false }

func (u *unsupportedUnityBridgeImpl) SetEventCallback(eventTypeCode uint64,
	cb callback.Callback) {
}

func (u *unsupportedUnityBridgeImpl) SendEvent(e *event.Event, data uintptr,
	tag uint64) {
}

func (u *unsupportedUnityBridgeImpl) SendEventWithString(e *event.Event,
	data string, tag uint64) {
}

func (u *unsupportedUnityBridgeImpl) SendEventWithNumber(e *event.Event, data,
	tag uint64) {
}

func (u *unsupportedUnityBridgeImpl) GetSecurityKeyByKeyChainIndex(
	index int) string {
	return ""
}

func (u *unsupportedUnityBridgeImpl) Uninitialize() {}

func (u *unsupportedUnityBridgeImpl) Destroy() {}
