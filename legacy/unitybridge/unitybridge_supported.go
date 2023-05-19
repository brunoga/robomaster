//go:build (darwin && amd64) || (android && arm) || (android && arm64) || (ios && arm64) || (windows && amd64)

package unitybridge

/*
#include <stdbool.h>
#include <stdlib.h>
#include <unitybridge.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

var libPaths = map[string]string{
	"darwin/amd64":  "./lib/darwin/amd64/unitybridge.bundle/Contents/MacOS/unitybridge",
	"android/arm":   "./lib/android/arm/libunitybridge.so",
	"android/arm64": "./lib/android/arm64/libunitybridge.so",
	"ios/arm64":     "./lib/ios/arm64/unitybridge.framework/unitybridge",
	"windows/amd64": "./lib/windows/amd64/unitybridge.dll",
}

func Create(name string, debuggable bool, logPath string) {
	libPath, ok := libPaths[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !ok {
		// Should never happen.
		panic(fmt.Sprintf("Platform \"%s/%s\" not supported by Unity Bridge",
			runtime.GOOS, runtime.GOARCH))
	}

	cLibPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cLibPath))

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cLogPath *C.char
	if logPath != "" {
		cLogPath = C.CString(logPath)
		defer C.free(unsafe.Pointer(cLogPath))
	} else {
		cLogPath = nil
	}

	C.CreateUnityBridge(cLibPath, cName, C.bool(debuggable), cLogPath)
}

func Destroy() {
	C.DestroyUnityBridge()
}

func Initialize() bool {
	return bool(C.UnityBridgeInitialize())
}

func Uninitialize() {
	C.UnityBridgeUninitialize()
}

func SetEventCallback(eventCode uint64, callback C.EventCallbackFunc) {
	C.UnitySetEventCallback(C.uint64_t(eventCode), callback)
}

func SendEvent(eventCode uint64, data []byte, tag uint64) {
	var dataPtr unsafe.Pointer = nil
	if len(data) != 0 {
		dataPtr = unsafe.Pointer(&data[0])
	}

	C.UnitySendEvent(C.uint64_t(eventCode), C.uintptr_t(uintptr(dataPtr)),
		C.int(len(data)), C.uint64_t(tag))
}

func SendEventWithString(eventCode uint64, data string, tag uint64) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.UnitySendEventWithString(C.uint64_t(eventCode), cData, C.uint64_t(tag))
}

func SendEventWithNumber(eventCode uint64, data uint64, tag uint64) {
	C.UnitySendEventWithNumber(C.uint64_t(eventCode), C.uint64_t(data),
		C.uint64_t(tag))
}

func GetSecurityKeyByKeychainIndex(index int) uintptr {
	return uintptr(C.UnityGetSecurityKeyByKeyChainIndex(C.int(index)))
}
