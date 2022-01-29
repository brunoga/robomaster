package wrapper

/*
#cgo LDFLAGS: -static -L${SRCDIR}/ffcall/lib -lffcall -lpthread

#include <stdlib.h>
#include "event_callback_windows.h"

// Calling C.alloc_callback() directly from Go code does not work for some
// reason, so we use a wrapper function to work around that.
callback_t generate_callback(void* data) {
	return alloc_callback(&event_callback, data);
}
*/
import "C"
import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"github.com/mattn/go-pointer"
)

type Windows struct {
	dllFileName string
	dll         *syscall.DLL

	procMap map[string]*syscall.Proc

	m           sync.Mutex
	callbackMap map[uint64]struct {
		context        unsafe.Pointer
		cEventCallback C.callback_t
	}
}

func New(unityBridgeDLLPath string) (Wrapper, error) {
	// Copy the DLL to a different name. This allows loading multiple DLLs
	// in the same program (one per wrapper instance).
	dstDLLFileName, err := copyDLL(unityBridgeDLLPath)
	if err != nil {
		return nil, err
	}

	procMap, dll, err := generateProcMap(dstDLLFileName)
	if err != nil {
		return nil, err
	}

	w := &Windows{
		dstDLLFileName,
		dll,
		procMap,
		sync.Mutex{},
		make(map[uint64]struct {
			context        unsafe.Pointer
			cEventCallback C.callback_t
		}),
	}
	runtime.SetFinalizer(w, func(w *Windows) {
		w.finalize()
	})

	return w, nil
}

func (w *Windows) finalize() {
	w.m.Lock()
	defer w.m.Unlock()

	for k, v := range w.callbackMap {
		pointer.Unref(v.context)
		C.free_callback(v.cEventCallback)
		delete(w.callbackMap, k)
	}

	err := w.dll.Release()
	if err != nil {
		log.Printf("Failed releasing DLL: %s.", err)
	}

	err = os.Remove(w.dllFileName)
	if err != nil {
		log.Printf("Failed removing DLL copy: %s", err)
	}
}

func (w *Windows) CreateUnityBridge(name string, debuggable bool) {
	fmt.Printf("CreateUnityBridge(name=%s, debuggable=%t)\n", name, debuggable)
	intDebuggable := 0
	if debuggable {
		intDebuggable = 1
	}

	cName := unsafe.Pointer(C.CString(name))

	w.procMap["CreateUnityBridge"].Call(uintptr(cName),
		uintptr(intDebuggable))

	C.free(cName)
}

func (w *Windows) DestroyUnityBridge() {
	fmt.Printf("DestroyUnityBridge()\n")
	w.procMap["DestroyUnityBridge"].Call()
}

func (w *Windows) UnityBridgeInitialize() bool {
	fmt.Printf("UnityBridgeInitialize()\n")
	ok, _, _ := w.procMap["UnityBridgeInitialize"].Call()

	return ok != 0
}

func (w *Windows) UnityBridgeUninitialize() {
	fmt.Printf("UnityBridgeUninitialize()\n")
	w.procMap["UnityBridgeUninitialze"].Call()
}

func (w *Windows) UnitySendEvent(eventCode uint64, data []byte, tag uint64) {
	event := unity.NewEventFromCode(eventCode)
	if event == nil {
		log.Printf("Unknown event with code %d (Type:%d, SubType:%d).\n",
			eventCode, eventCode<<32, eventCode&0xffffffff)
		return
	}

	if event == nil {
		log.Printf("Unknown event with code %d (Type:%d, SubType:%d).\n",
			eventCode, eventCode<<32, eventCode&0xffffffff)
		return
	}

	fmt.Printf("UnitySendEvent(event=%v, data=%v, tag=%d)\n",
		event, data, tag)

	var dataPtr unsafe.Pointer = nil
	if len(data) != 0 {
		dataPtr = unsafe.Pointer(&data[0])
	}

	w.procMap["UnitySendEvent"].Call(uintptr(eventCode),
		uintptr(dataPtr), uintptr(tag))
}

func (w *Windows) UnitySendEventWithString(eventCode uint64, data string, tag uint64) {
	event := unity.NewEventFromCode(eventCode)
	if event == nil {
		log.Printf("Unknown event with code %d (Type:%d, SubType:%d).\n",
			eventCode, eventCode<<32, eventCode&0xffffffff)
		return
	}

	if event == nil {
		log.Printf("Unknown event with code %d (Type:%d, SubType:%d).\n",
			eventCode, eventCode<<32, eventCode&0xffffffff)
		return
	}

	fmt.Printf("UnitySendEventWithString(event=%v, data=%s, tag=%d)\n",
		event, data, tag)

	cData := unsafe.Pointer(C.CString(data))

	w.procMap["UnitySendEventWithString"].Call(uintptr(eventCode),
		uintptr(cData), uintptr(tag))

	C.free(cData)
}

func (w *Windows) UnitySendEventWithNumber(eventCode uint64, data uint64, tag uint64) {
	event := unity.NewEventFromCode(eventCode)
	if event == nil {
		log.Printf("Unknown event with code %d (Type:%d, SubType:%d).\n",
			eventCode, eventCode<<32, eventCode&0xffffffff)
		return
	}

	if event == nil {
		log.Printf("Unknown event with code %d (Type:%d, SubType:%d).\n",
			eventCode, eventCode<<32, eventCode&0xffffffff)
		return
	}

	fmt.Printf("UnitySendEventWithNumber(event=%v, data=%d, tag=%d)\n",
		event, data, tag)

	w.procMap["UnitySendEventWithNumber"].Call(uintptr(eventCode),
		uintptr(data), uintptr(tag))
}

func (w *Windows) UnitySetEventCallback(eventCode uint64,
	eventCallback EventCallback) {
	fmt.Printf("UnitySetEventCallback(eventCode=%d, eventCallback=%p)\n",
		eventCode, eventCallback)
	w.m.Lock()
	defer w.m.Unlock()

	callbackData, hasCallback := w.callbackMap[eventCode]

	if eventCallback == nil {
		w.procMap["UnitySetEventCallback"].Call(uintptr(eventCode),
			uintptr(unsafe.Pointer(nil)))

		if hasCallback {
			pointer.Unref(callbackData.context)
			C.free_callback(callbackData.cEventCallback)
			delete(w.callbackMap, eventCode)
		} else {
			log.Printf("No callback set for event code %d.",
				eventCode)
		}

		return
	}

	if hasCallback {
		log.Printf("Callback already set for event code %d.",
			eventCode)

		// Set again with the existing callback just in case.
		w.procMap["UnitySetEventCallback"].Call(uintptr(eventCode),
			uintptr(unsafe.Pointer(callbackData.cEventCallback)))

		return
	}

	context := pointer.Save(eventCallback)
	cEventCallback := C.generate_callback(context)

	w.callbackMap[eventCode] = struct {
		context        unsafe.Pointer
		cEventCallback C.callback_t
	}{
		context,
		cEventCallback,
	}

	w.procMap["UnitySetEventCallback"].Call(uintptr(eventCode),
		uintptr(unsafe.Pointer(cEventCallback)))
}

func copyDLL(unityBridgeDLLPath string) (string, error) {
	fi, err := os.Stat(unityBridgeDLLPath)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		return "", fmt.Errorf("%q is a directory")
	}

	srcDLLFile, err := os.Open(unityBridgeDLLPath)
	if err != nil {
		return "", err
	}
	defer srcDLLFile.Close()

	dstDLLFile, err := ioutil.TempFile("", "unitybridge-*.dll")
	if err != nil {
		return "", err
	}
	defer dstDLLFile.Close()

	_, err = io.Copy(dstDLLFile, srcDLLFile)
	if err != nil {
		return "", err
	}

	return dstDLLFile.Name(), nil
}

func generateProcMap(dstDLLFileName string) (map[string]*syscall.Proc,
	*syscall.DLL, error) {
	dll, err := syscall.LoadDLL(dstDLLFileName)
	if err != nil {
		return nil, nil, err
	}

	procNames := []string{
		"CreateUnityBridge",
		"DestroyUnityBridge",
		"UnityBridgeInitialize",
		"UnityBridgeUninitialze",
		"UnitySendEvent",
		"UnitySendEventWithString",
		"UnitySendEventWithNumber",
		"UnitySetEventCallback",
	}

	procMap := make(map[string]*syscall.Proc)
	for _, procName := range procNames {
		proc, err := dll.FindProc(procName)
		if err != nil {
			return nil, nil, err
		}

		procMap[procName] = proc
	}

	return procMap, dll, nil
}
