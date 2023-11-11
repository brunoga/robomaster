//go:build (android && (arm || arm64)) || (darwin && amd64)

package implementations

/*
#include <dlfcn.h>
#include <stdlib.h>
#include <stdbool.h>

#include "../callback/callback.h"

void CreateUnityBridgeCaller(void *f, const char *name, bool debuggable,
                             const char *log_path) {
  ((void (*)(const char *, bool, const char *))f)(name, debuggable, log_path);
}

void DestroyUnityBridgeCaller(void *f) { ((void (*)())f)(); }

bool UnityBridgeInitializeCaller(void *f) { return ((bool (*)())f)(); }

void UnityBridgeUninitializeCaller(void *f) { ((void (*)())f)(); }

void UnitySendEventCaller(void *f, uint64_t event_code, intptr_t data,
	                      uint64_t tag) {
  ((void (*)(uint64_t, uintptr_t, uint64_t))f)(event_code, data, tag);
}

void UnitySendEventWithStringCaller(void *f, uint64_t event_code,
                                    const char *data, uint64_t tag) {
  ((void (*)(uint64_t, const char *, uint64_t))f)(event_code, data, tag);
}

void UnitySendEventWithNumberCaller(void *f, uint64_t event_code, uint64_t data,
                                    uint64_t tag) {
  ((void (*)(uint64_t, uint64_t, uint64_t))f)(event_code, data, tag);
}

void UnitySetEventCallbackCaller(void *f, uint64_t event_code,
                                 EventCallback event_callback) {
  ((void (*)(uint64_t, EventCallback))f)(event_code, event_callback);
}

char* UnityGetSecurityKeyByKeyChainIndexCaller(void *f, int index) {
  return (char*)((uintptr_t(*)(int))f)(index);
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/wrapper/callback"

	internal_callback "github.com/brunoga/unitybridge/wrapper/internal/callback"
)

var (
	libPaths = map[string]string{
		"android/arm":   "./lib/android/arm/libunitybridge.so",
		"android/arm64": "./lib/android/arm64/libunitybridge.so",
		"darwin/amd64":  "./lib/darwin/amd64/unitybridge.bundle/Contents/MacOS/unitybridge",
	}

	UnityBridgeImpl *dlOpenUnityBridgeImpl = &dlOpenUnityBridgeImpl{}
)

func init() {
	libPath, ok := libPaths[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !ok {
		// Should never happen.
		panic(fmt.Sprintf("Platform \"%s/%s\" not supported by Unity Bridge "+
			"library", runtime.GOOS, runtime.GOARCH))
	}

	cLibPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cLibPath))

	UnityBridgeImpl.handle = C.dlopen(cLibPath, C.RTLD_NOW)
	if UnityBridgeImpl.handle == nil {
		cError := C.dlerror()

		panic(fmt.Sprintf("Could not load Unity Bridge library at \"%s\": %s",
			libPath, C.GoString(cError)))
	}

	UnityBridgeImpl.createUnityBridge =
		UnityBridgeImpl.getSymbol("CreateUnityBridge")
	UnityBridgeImpl.destroyUnityBridge =
		UnityBridgeImpl.getSymbol("DestroyUnityBridge")
	UnityBridgeImpl.unityBridgeInitialize =
		UnityBridgeImpl.getSymbol("UnityBridgeInitialize")
	UnityBridgeImpl.unityBridgeUninitialize =
		UnityBridgeImpl.getSymbol("UnityBridgeUninitialze") // Typo in C code.
	UnityBridgeImpl.unitySendEvent =
		UnityBridgeImpl.getSymbol("UnitySendEvent")
	UnityBridgeImpl.unitySendEventWithString =
		UnityBridgeImpl.getSymbol("UnitySendEventWithString")
	UnityBridgeImpl.unitySendEventWithNumber =
		UnityBridgeImpl.getSymbol("UnitySendEventWithNumber")
	UnityBridgeImpl.unitySetEventCallback =
		UnityBridgeImpl.getSymbol("UnitySetEventCallback")
	UnityBridgeImpl.UnityGetSecurityKeyByKeyChainIndex =
		UnityBridgeImpl.getSymbol("UnityGetSecurityKeyByKeyChainIndex")
}

type dlOpenUnityBridgeImpl struct {
	handle unsafe.Pointer

	createUnityBridge                  unsafe.Pointer
	destroyUnityBridge                 unsafe.Pointer
	unityBridgeInitialize              unsafe.Pointer
	unityBridgeUninitialize            unsafe.Pointer
	unitySendEvent                     unsafe.Pointer
	unitySendEventWithString           unsafe.Pointer
	unitySendEventWithNumber           unsafe.Pointer
	unitySetEventCallback              unsafe.Pointer
	UnityGetSecurityKeyByKeyChainIndex unsafe.Pointer

	l *logger.Logger
}

func Get(l *logger.Logger) *dlOpenUnityBridgeImpl {
	UnityBridgeImpl.l = l

	return UnityBridgeImpl
}

func (d *dlOpenUnityBridgeImpl) Create(name string, debuggable bool,
	logPath string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cLogPath := C.CString(logPath)
	defer C.free(unsafe.Pointer(cLogPath))

	C.CreateUnityBridgeCaller(unsafe.Pointer(d.createUnityBridge), cName,
		C.bool(debuggable), cLogPath)
}

func (d *dlOpenUnityBridgeImpl) Initialize() bool {
	return bool(C.UnityBridgeInitializeCaller(d.unityBridgeInitialize))
}

func (d *dlOpenUnityBridgeImpl) SetEventCallback(eventTypeCode uint64,
	c callback.Callback) {
	var eventCallback C.EventCallback
	if c != nil {
		eventCallback = C.EventCallback(C.eventCallbackC)
	}

	C.UnitySetEventCallbackCaller(unsafe.Pointer(d.unitySetEventCallback),
		C.uint64_t(eventTypeCode), eventCallback)

	internal_callback.Set(eventTypeCode, c)
}

func (d *dlOpenUnityBridgeImpl) SendEvent(eventCode uint64, output []byte,
	tag uint64) {
	var outputUintptr uintptr
	if len(output) > 0 {
		outputUintptr = uintptr(unsafe.Pointer(&output[0]))
	}

	C.UnitySendEventCaller(unsafe.Pointer(d.unitySendEvent),
		C.uint64_t(eventCode), C.intptr_t(outputUintptr), C.uint64_t(tag))
}

func (d *dlOpenUnityBridgeImpl) SendEventWithString(eventCode uint64,
	data string, tag uint64) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.UnitySendEventWithStringCaller(unsafe.Pointer(d.unitySendEventWithString),
		C.uint64_t(eventCode), cData, C.uint64_t(tag))
}

func (d *dlOpenUnityBridgeImpl) SendEventWithNumber(eventCode, data,
	tag uint64) {
	C.UnitySendEventWithNumberCaller(unsafe.Pointer(d.unitySendEventWithNumber),
		C.uint64_t(eventCode), C.uint64_t(data), C.uint64_t(tag))
}

func (d *dlOpenUnityBridgeImpl) GetSecurityKeyByKeyChainIndex(
	index int) string {
	cKey := C.UnityGetSecurityKeyByKeyChainIndexCaller(
		unsafe.Pointer(d.UnityGetSecurityKeyByKeyChainIndex), C.int(index))
	defer C.free(unsafe.Pointer(cKey))

	return C.GoString(cKey)
}

func (d *dlOpenUnityBridgeImpl) Uninitialize() {
	C.UnityBridgeUninitializeCaller(unsafe.Pointer(d.unityBridgeUninitialize))
}

func (d *dlOpenUnityBridgeImpl) Destroy() {
	C.DestroyUnityBridgeCaller(unsafe.Pointer(d.destroyUnityBridge))
}

func (d *dlOpenUnityBridgeImpl) getSymbol(name string) unsafe.Pointer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	symbol := C.dlsym(d.handle, cName)
	if symbol == nil {
		cError := C.dlerror()

		panic(fmt.Sprintf("Could not load symbol \"%s\": %s",
			name, C.GoString(cError)))
	}

	return symbol
}
