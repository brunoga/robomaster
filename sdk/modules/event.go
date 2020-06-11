package modules

import (
	"fmt"
	"net"
	"sync"
)

const (
	eventAddrPort = ":40925"
)

type EventHandler func(string)

type Event struct {
	control *Control

	m sync.Mutex
	quitChan chan struct{}
	eventHandlers map[string]map[int]EventHandler
}

func (e *Event) StartListening(eventType, eventParameters string,
	eventHandler EventHandler) (int, error) {
	e.m.Lock()
	defer e.m.Unlock()

	tokenHandlerMap, ok := e.eventHandlers[eventType]
	if !ok {
		err := e.control.SendDataExpectOk(fmt.Sprintf(
			"%s %s", eventType, eventParameters))
		if err != nil {
			return -1, fmt.Errorf("error listening for event: %w", err)
		}

		tokenHandlerMap = make(map[int]EventHandler)

		e.eventHandlers[eventType] = tokenHandlerMap
	}

	if len(e.eventHandlers[eventType]) == 1 {
		go e.eventLoop()
	}

	for i := 0; i < len(tokenHandlerMap) + 1; i++ {
		_, ok = tokenHandlerMap[i]
		if ok {
			continue
		}

		tokenHandlerMap[i] = eventHandler

		return i, nil
	}

	return -1, fmt.Errorf("event handler tokens exhausted")
}

func (e *Event) StopListening(eventType, eventParameters string,
	token int) error {
	e.m.Lock()
	defer e.m.Unlock()

	tokenHandlerMap, ok := e.eventHandlers[eventType]
	if !ok {
		return fmt.Errorf("no handlers for event type")
	}

	_, ok = tokenHandlerMap[token]
	if !ok {
		return fmt.Errorf("token does not match event type")
	}

	delete(tokenHandlerMap, token)

	if len(tokenHandlerMap) == 0 {
		delete(e.eventHandlers, eventType)
	}

	if len(e.eventHandlers[eventType]) == 0 {
		err := e.control.SendDataExpectOk(fmt.Sprintf(
			"%s %s", eventType, eventParameters))
		if err != nil {
			return fmt.Errorf("error stopping listening for event: %w",
				err)
		}
	}

	return nil
}

func (e *Event) eventLoop() {
	ip, err := e.control.IP()
	if err != nil {
		// TODO(bga): Log this.
		return
	}

	eventAddr := ip.String() + eventAddrPort

	conn, err := net.Dial("tcp", eventAddr)
	if err != nil {
		// TODO(bga): Log this.
		return
	}
	defer conn.Close()

	b := make([]byte, 512)
L:
	for {
		select {
		case <-e.quitChan:
			break L
		default:
			_, err := conn.Read(b)
			if err != nil {
				// TODO(bga): Log this.
				break L
			}

			//data := b[:n]
			e.m.Lock()
			// TODO(bga): Send event to listeners here.
			e.m.Unlock()
		}
	}

	e.quitChan = nil
}