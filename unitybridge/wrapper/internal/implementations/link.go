//go:build android && (arm || arm64)

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
extern char* UnityGetSecurityKeyByKeyChainIndex(int index);
extern void UnityBridgeUninitialze();
extern void DestroyUnityBridge();
*/
import "C"

import (
	"bytes"
	"log/slog"
	"unsafe"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge/wrapper/callback"

	internal_callback "github.com/brunoga/robomaster/unitybridge/wrapper/internal/callback"
)

var (
	// Singleton.
	UnityBridgeImpl *linkUnityBridgeImpl = &linkUnityBridgeImpl{}
)

type linkUnityBridgeImpl struct {
	l *logger.Logger
	m *internal_callback.Manager
}

func Get(l *logger.Logger) *linkUnityBridgeImpl {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("unity_bridge_wrapper")

	UnityBridgeImpl.l = l
	UnityBridgeImpl.m = internal_callback.NewManager(l)

	l.Debug("Unity Bridge implementation loaded", "implememntation", "link")

	return UnityBridgeImpl
}

func (u *linkUnityBridgeImpl) Create(name string, debuggable bool, logPath string) {
	defer u.l.Trace("Create", "name", name, "debuggable", debuggable,
		"logPath", logPath)()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cLogPath := C.CString(logPath)
	defer C.free(unsafe.Pointer(cLogPath))

	C.CreateUnityBridge(cName, C.bool(debuggable), cLogPath)
}

func (u *linkUnityBridgeImpl) Initialize() (initialized bool) {
	endTrace := u.l.Trace("Initialize")
	defer func() {
		endTrace("initialized", initialized)
	}()

	return bool(C.UnityBridgeInitialize())
}

func (u *linkUnityBridgeImpl) SetEventCallback(eventTypeCode uint64,
	c callback.Callback) {
	defer u.l.Trace("SetEventCallback", "eventTypeCode", eventTypeCode,
		"callback", c)()

	var eventCallback C.EventCallback
	if c != nil {
		eventCallback = C.EventCallback(C.eventCallbackC)
	}

	C.UnitySetEventCallback(C.uint64_t(eventTypeCode), eventCallback)

	u.m.Set(eventTypeCode, c)
}

func (u *linkUnityBridgeImpl) SendEvent(eventCode uint64, output []byte,
	tag uint64) {
	endTrace := u.l.Trace("SendEvent", "eventCode", eventCode, "len(output)",
		len(output), "tag", tag)
	defer func() {
		zeroPos := bytes.Index(output, []byte{0})
		if zeroPos == -1 {
			endTrace("output", output)
		} else {
			endTrace("output", output[0:zeroPos])
		}
	}()

	var outputUintptr uintptr
	if len(output) > 0 {
		outputUintptr = uintptr(unsafe.Pointer(&output[0]))
	}

	C.UnitySendEvent(C.uint64_t(eventCode), C.intptr_t(outputUintptr),
		C.uint64_t(tag))
}

func (u *linkUnityBridgeImpl) SendEventWithString(eventCode uint64, data string,
	tag uint64) {
	defer u.l.Trace("SendEventWithString", "eventCode", eventCode,
		"data", data, "tag", tag)()

	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.UnitySendEventWithString(C.uint64_t(eventCode), cData, C.uint64_t(tag))
}

func (u *linkUnityBridgeImpl) SendEventWithNumber(eventCode, data,
	tag uint64) {
	defer u.l.Trace("SendEventWithNumber", "eventCode", eventCode, "data",
		data, "tag", tag)()

	C.UnitySendEventWithNumber(C.uint64_t(eventCode), C.uint64_t(data),
		C.uint64_t(tag))
}

func (u *linkUnityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int) (securityKey string) {
	endTrace := u.l.Trace("GetSecurityKeyByKeyChainIndex", "index", index)
	defer func() {
		endTrace("securityKey", securityKey)
	}()

	cKey := C.UnityGetSecurityKeyByKeyChainIndex(C.int(index))
	defer C.free(unsafe.Pointer(cKey))

	return C.GoString((*C.char)(unsafe.Pointer(cKey)))
}

func (u *linkUnityBridgeImpl) Uninitialize() {
	defer u.l.Trace("Uninitialize")()

	C.UnityBridgeUninitialze()
}

func (u *linkUnityBridgeImpl) Destroy() {
	defer u.l.Trace("Destroy")()

	C.DestroyUnityBridge()
}
