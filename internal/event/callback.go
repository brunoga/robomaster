package unitybridge

import (
	"C"
	"fmt"
	"log"
	"sync"

	"github.com/brunoga/unitybridge/event"
)

var (
	m                         sync.Mutex
	eventCodeEventCallbackMap = make(map[uint64]event.Callback)
)

func SetEventCallback(eventCode uint64, callback event.Callback) error {
	m.Lock()
	defer m.Unlock()

	_, ok := eventCodeEventCallbackMap[eventCode]
	if !ok {
		if callback == nil {
			return fmt.Errorf("no callback for event code %d found", eventCode)
		}

		eventCodeEventCallbackMap[eventCode] = callback
	} else {
		if callback == nil {
			delete(eventCodeEventCallbackMap, eventCode)
		} else {
			return fmt.Errorf("callback for event code %d already set",
				eventCode)
		}
	}

	return nil
}

func RunEventCallback(eventCode uint64, data []byte, tag uint64) error {
	m.Lock()
	defer m.Unlock()

	callback, ok := eventCodeEventCallbackMap[eventCode]
	if !ok {
		return fmt.Errorf("no handlers for event code %d", eventCode)
	}

	ev := event.NewFromCode(eventCode)

	go callback(ev, data, tag)

	return nil
}

//export eventCallbackGo
func eventCallbackGo(eventCode uint64, data []byte, tag uint64) {
	err := RunEventCallback(eventCode, data, tag)
	if err != nil {
		log.Printf("error running event callback: %s\n", err)
	}
}
