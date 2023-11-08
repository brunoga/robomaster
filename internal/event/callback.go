package event

import (
	"C"
	"fmt"
	"log"
	"sync"

	"github.com/brunoga/unitybridge/unity/event"
)

var (
	m                         sync.Mutex
	eventCodeEventCallbackMap = make(map[event.Type]event.Callback)
)

func SetEventCallback(eventType event.Type, callback event.Callback) error {
	m.Lock()
	defer m.Unlock()

	_, ok := eventCodeEventCallbackMap[eventType]
	if !ok {
		if callback == nil {
			return fmt.Errorf("no callback for event type %q found", eventType)
		}

		eventCodeEventCallbackMap[eventType] = callback
	} else {
		if callback == nil {
			delete(eventCodeEventCallbackMap, eventType)
		} else {
			return fmt.Errorf("callback for event type %q already set",
				eventType)
		}
	}

	return nil
}

func RunEventCallback(eventCode uint64, data []byte, tag uint64) error {
	m.Lock()
	defer m.Unlock()

	e := event.NewFromCode(eventCode)

	callback, ok := eventCodeEventCallbackMap[e.Type()]
	if !ok {
		return fmt.Errorf("no handlers for event type %s", e.Type())
	}

	go callback(e, data, tag)

	return nil
}

//export eventCallbackGo
func eventCallbackGo(eventCode uint64, data []byte, tag uint64) {
	err := RunEventCallback(eventCode, data, tag)
	if err != nil {
		log.Printf("error running event callback: %s\n", err)
	}
}
