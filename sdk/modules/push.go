package modules

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	pushAddrPort = ":40924"
)

// PushHandler is a handler for push notifications. The string parameter will contain
// the data for the specific event being monitored usually starting with the
// attribute name. Implementations must parse the data before using it.
type PushHandler func(string)

// Push handles robot's push notifications, starting/stopping monitoring individual
// events and sending them to registered PushHandlers.
type Push struct {
	control *Control

	m            sync.RWMutex
	quitChan     chan struct{}
	pushHandlers map[string]map[string]map[int]PushHandler
}

// NewPush returns a new Push instance. The control parameter is used to start
// stop the specific notification pushes.
func NewPush(control *Control) *Push {
	return &Push{
		control,
		sync.RWMutex{},
		nil,
		make(map[string]map[string]map[int]PushHandler),
	}
}

// StartListening starts sending push notifications of type pushType to the
// given pushHandler. If no one is listening to a specific event yet, starts
// the push notifications. Returns a token (to be used to stop receiving push
// notifications) and a nil error on success and a non-nil error on failure.
func (p *Push) StartListening(pushType, pushAttribute, pushParameters string,
	pushHandler PushHandler) (int, error) {
	p.m.Lock()
	defer p.m.Unlock()

	attributeTokenMap, ok := p.pushHandlers[pushType]
	if !ok {
		attributeTokenMap = make(map[string]map[int]PushHandler)
	}

	startTracking := false

	tokenHandlerMap, ok := attributeTokenMap[pushAttribute]
	if !ok {
		startTracking = true
		tokenHandlerMap = make(map[int]PushHandler)
	}

	token := -1
	for i := 0; i < len(tokenHandlerMap)+1; i++ {
		_, ok = tokenHandlerMap[i]
		if ok {
			continue
		}

		token = i

		break
	}

	if token == -1 {
		// Should never happen unless there is a bug.
		return -1, fmt.Errorf("can't obtain push handler token")
	}

	if startTracking {
		// This eventType/eventAttribute pair is not being tracked yet. Start
		// tracking.
		var command string
		if pushParameters == "" {
			command = fmt.Sprintf("%s %s on;", pushType, pushAttribute)
		} else {
			command = fmt.Sprintf("%s %s on %s;", pushType, pushAttribute,
				pushParameters)
		}

		err := p.control.SendDataExpectOk(command)
		if err != nil {
			return -1, fmt.Errorf("error listening for push notifications: %w", err)
		}
	}

	if len(p.pushHandlers) == 0 {
		go p.pushLoop()
	}

	tokenHandlerMap[token] = pushHandler
	attributeTokenMap[pushAttribute] = tokenHandlerMap
	p.pushHandlers[pushType] = attributeTokenMap

	return token, nil
}

// StopListening stops sending push notifications of type pushType to the
// handler represented by the given pushType and token. If all listeners of a
// specific push notification are removed, stops the push notifications.
// Returns a nil error on success and a non-nil error on failure.
func (p *Push) StopListening(pushType, pushAttribute string,
	token int) error {
	p.m.Lock()
	defer p.m.Unlock()

	attributeTokenMap, ok := p.pushHandlers[pushType]
	if !ok {
		return fmt.Errorf("no handlers for push type")
	}

	tokenHandlerMap, ok := attributeTokenMap[pushType]
	if !ok {
		return fmt.Errorf("no handlers for push attribute")
	}

	_, ok = tokenHandlerMap[token]
	if !ok {
		return fmt.Errorf("token does not match push type")
	}

	delete(tokenHandlerMap, token)

	if len(tokenHandlerMap) == 0 {
		// This pushType/pushAttribute pair is not being tracked anymore.
		// Stop tracking.
		err := p.control.SendDataExpectOk(fmt.Sprintf(
			"%s %s;", pushType, pushAttribute))
		if err != nil {
			return fmt.Errorf("error stopping push notifications: %w",
				err)
		}

		delete(attributeTokenMap, pushAttribute)
	}

	if len(attributeTokenMap) == 0 {
		delete(p.pushHandlers, pushType)
	}

	if len(p.pushHandlers) == 0 {
		close(p.quitChan)
	}

	return nil
}

func (p *Push) pushLoop() {
	p.quitChan = make(chan struct{})

	conn, err := net.ListenPacket("udp", pushAddrPort)
	if err != nil {
		// TODO(bga): Log this.
		return
	}
	defer conn.Close()

	b := make([]byte, 512)
L:
	for {
		select {
		case <-p.quitChan:
			break L
		default:
			n, addr, err := conn.ReadFrom(b)
			if err != nil {
				// TODO(bga): Log this.
				break L
			}

			robotIP, err := p.control.IP()
			if err != nil {
				// TODO(bga): Log this.
				break L
			}

			if robotIP.String() != addr.(*net.UDPAddr).IP.String() {
				// Got push notification from an unexpected ip. Ignore it.
				continue
			}

			pushType, pushAttribute, pushData, err :=
				getPushTypeAttributeAndData(b[:n])
			if err != nil {
				// TODO(bga): Log this.
				continue
			}

			p.m.RLock()

			attributeTokenMap, ok := p.pushHandlers[pushType]
			if !ok {
				// TODO(bga): Log this.
				continue
			}

			tokenPushHandlerMap, ok := attributeTokenMap[pushAttribute]
			if !ok {
				// TODO(bga): Log this.
				continue
			}

			for _, pushHandler := range tokenPushHandlerMap {
				pushHandler(pushData)
			}

			p.m.RUnlock()
		}
	}

	p.quitChan = nil
}

func getPushTypeAttributeAndData(
	receivedData []byte) (string, string, string, error) {
	fields := bytes.Fields(receivedData)
	if len(fields) < 4 {
		return "", "", "", fmt.Errorf("invalid data received")
	}

	pushType := fmt.Sprintf("%s %s", string(fields[0]),
		string(fields[1]))

	pushAttribute := string(fields[2])

	pushDataBuilder := strings.Builder{}
	for i := 3; i < len(fields); i++ {
		if i != 3 {
			pushDataBuilder.WriteByte(' ')
		}
		pushDataBuilder.Write(fields[i])
	}

	pushData := pushDataBuilder.String()

	return pushType, pushAttribute, pushData, nil
}
