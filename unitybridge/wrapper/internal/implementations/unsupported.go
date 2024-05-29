//go:build !(windows && amd64) && !(ios && arm64) && !(android && (arm || arm64)) && !(darwin && amd64) && !(linux && amd64)

package implementations

import (
	"fmt"
	"runtime"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge/wrapper/callback"
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

func (u *unsupportedUnityBridgeImpl) SendEvent(eventCode uint64, output []byte,
	tag uint64) {
}

func (u *unsupportedUnityBridgeImpl) SendEventWithString(eventCode uint64,
	data string, tag uint64) {
}

func (u *unsupportedUnityBridgeImpl) SendEventWithNumber(eventCode, data,
	tag uint64) {
}

func (u *unsupportedUnityBridgeImpl) GetSecurityKeyByKeyChainIndex(
	index int) string {
	return ""
}

func (u *unsupportedUnityBridgeImpl) Uninitialize() {}

func (u *unsupportedUnityBridgeImpl) Destroy() {}
