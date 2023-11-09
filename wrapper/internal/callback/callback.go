package callback

import (
	"C"
	"fmt"
	"log"
	"sync"
)
import "github.com/brunoga/unitybridge/wrapper/callback"

var (
	m                         sync.Mutex
	eventCodeEventCallbackMap = make(map[uint64]callback.Callback)
)

type Callback func(eventCode uint64, data []byte, tag uint64)

func Set(eventTypeCode uint64, c callback.Callback) error {
	m.Lock()
	defer m.Unlock()

	// Make sure this refers to an event type (i.e unset subtype).
	eventTypeCode &= 0xFFFFFFFF00000000

	_, ok := eventCodeEventCallbackMap[eventTypeCode]
	if !ok {
		if c == nil {
			return fmt.Errorf("no callback for event type %d found", eventTypeCode)
		}

		eventCodeEventCallbackMap[eventTypeCode] = c
	} else {
		if c == nil {
			delete(eventCodeEventCallbackMap, eventTypeCode)
		} else {
			return fmt.Errorf("callback for event type code %d already set",
				eventTypeCode)
		}
	}

	return nil
}

func Run(eventCode uint64, data []byte, tag uint64) error {
	m.Lock()
	defer m.Unlock()

	eventTypeCode := eventCode & 0xFFFFFFFF00000000

	c, ok := eventCodeEventCallbackMap[eventTypeCode]
	if !ok {
		return fmt.Errorf("no handlers for event type code %d", eventTypeCode)
	}

	// TODO(bga): Maybe do this in a goroutine? If we do, we must copy data as
	// it is backed by a C array that is freed after Run() returns.
	c(eventCode, data, tag)

	return nil
}

//export eventCallbackGo
func eventCallbackGo(eventCode uint64, data []byte, tag uint64) {
	err := Run(eventCode, data, tag)
	if err != nil {
		log.Printf("error running event callback: %s\n", err)
	}
}
