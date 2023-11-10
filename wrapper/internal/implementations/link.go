//go:build ios && arm64

package implementations

/*
#include <stdbool.h>
#include <stdlib.h>

#include "../callback/callback.h"

extern void CreateUnityBridge(const char* name, bool debuggable, const char* logPath);
extern bool UnityBridgeInitialize();
extern void UnitySendEvent(uint64_t event_code, intptr_t data, uint64_t tag);
extern void UnitySendEventWithString(uint64_t event_code, const char* data, uint64_t tag);
extern void UnitySendEventWithNumber(uint64_t event_code, uint64_t data, uint64_t tag);
extern void UnitySetEventCallback(uint64_t event_code, EventCallback event_callback);
extern intptr_t UnityGetSecurityKeyByKeyChainIndex(int index);
extern void UnityBridgeUninitialze();
extern void DestroyUnityBridge();
*/
import "C"

import (
	"unsafe"

	"github.com/brunoga/unitybridge/wrapper/callback"

	internal_callback "github.com/brunoga/unitybridge/wrapper/internal/callback"
)

var (
	// Singleton.
	UnityBridgeImpl *linkUnityBridgeImpl = &linkUnityBridgeImpl{}
)

type linkUnityBridgeImpl struct{}

func (u *linkUnityBridgeImpl) Create(name string, debuggable bool, logPath string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cLogPath := C.CString(logPath)
	defer C.free(unsafe.Pointer(cLogPath))

	C.CreateUnityBridge(cName, C.bool(debuggable), cLogPath)
}

func (u *linkUnityBridgeImpl) Initialize() bool {
	return bool(C.UnityBridgeInitialize())
}

func (u *linkUnityBridgeImpl) SetEventCallback(eventTypeCode uint64,
	c callback.Callback) {
	var eventCallback C.EventCallback
	if c != nil {
		eventCallback = C.EventCallback(C.eventCallbackC)
	}

	C.UnitySetEventCallback(C.uint64_t(eventTypeCode), eventCallback)

	internal_callback.Set(eventTypeCode, c)
}

func (u *linkUnityBridgeImpl) SendEvent(eventCode uint64, output []byte,
	tag uint64) {
	var outputUintptr uintptr
	if len(output) > 0 {
		outputUintptr = uintptr(unsafe.Pointer(&output[0]))
	}

	C.UnitySendEvent(C.uint64_t(eventCode), C.intptr_t(outputUintptr),
		C.uint64_t(tag))
}

func (u *linkUnityBridgeImpl) SendEventWithString(eventCode uint64, data string,
	tag uint64) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.UnitySendEventWithString(C.uint64_t(eventCode), cData, C.uint64_t(tag))
}

func (u *linkUnityBridgeImpl) SendEventWithNumber(eventCode, data,
	tag uint64) {
	C.UnitySendEventWithNumber(C.uint64_t(eventCode), C.uint64_t(data),
		C.uint64_t(tag))
}

func (u *linkUnityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int) string {
	cKey := C.UnityGetSecurityKeyByKeyChainIndex(C.int(index))
	defer C.free(unsafe.Pointer(uintptr(cKey)))

	return C.GoString((*C.char)(unsafe.Pointer(uintptr(cKey))))
}

func (u *linkUnityBridgeImpl) Uninitialize() {
	C.UnityBridgeUninitialze()
}

func (u *linkUnityBridgeImpl) Destroy() {
	C.DestroyUnityBridge()
}
