package push

import (
	"bytes"
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules"
	"net"
	"strings"
	"sync"
)

const (
	pushAddrPort = ":40924"
)

// Handler is a handler for push notifications. The string parameter will contain
// the data for the specific event being monitored usually starting with the
// attribute name. Implementations must parse the data before using it.
type Handler func(string)

// Push handles robot's push notifications, starting/stopping monitoring individual
// events and sending them to registered PushHandlers.
type Push struct {
	control *modules.Control

	m            sync.Mutex
	quitChan     chan struct{}
	pushHandlers map[string]map[int]Handler
}

// NewPush returns a new Push instance. The control parameter is used to start
// stop the specific notification pushes.
func NewPush(control *modules.Control) *Push {
	return &Push{
		control,
		sync.Mutex{},
		nil,
		make(map[string]map[int]Handler),
	}
}

// StartListening starts sending push notifications of type pushType to the
// given pushHandler. If no one is listening to a specific event yet, starts
// the push notifications. Returns a token (to be used to stop receiving push
// notifications) and a nil error on success and a non-nil error on failure.
func (p *Push) StartListening(
	req StartListeningRequest) (StartListeningResponse, error) {
	p.m.Lock()
	defer p.m.Unlock()

	tokenHandlerMap, ok := p.pushHandlers[req.Type]
	if !ok {
		err := p.control.SendDataExpectOk(fmt.Sprintf(
			"%s %s;", req.Type, req.Parameters))
		if err != nil {
			return StartListeningResponse{
				-1,
			}, fmt.Errorf("error listening for push notifications: %w", err)
		}

		tokenHandlerMap = make(map[int]Handler)

		p.pushHandlers[req.Type] = tokenHandlerMap
	}

	if len(p.pushHandlers[req.Type]) == 0 {
		go p.pushLoop()
	}

	for i := 0; i < len(tokenHandlerMap)+1; i++ {
		_, ok = tokenHandlerMap[i]
		if ok {
			continue
		}

		tokenHandlerMap[i] = req.Handler

		return StartListeningResponse{
			i,
		}, nil
	}

	return StartListeningResponse{
		-1,
	}, fmt.Errorf("push handler tokens exhausted")
}

// StopListening stops sending push notifications of type pushType to the
// handler represented by the given pushType and token. If all listeners of a
// specific push notification are removed, stops the push notifications.
// Returns a nil error on success and a non-nil error on failure.
func (p *Push) StopListening(req StopListeningRequest) error {
	p.m.Lock()
	defer p.m.Unlock()

	tokenHandlerMap, ok := p.pushHandlers[req.Type]
	if !ok {
		return fmt.Errorf("no handlers for push type")
	}

	_, ok = tokenHandlerMap[req.Token]
	if !ok {
		return fmt.Errorf("token does not match push type")
	}

	delete(tokenHandlerMap, req.Token)

	if len(tokenHandlerMap) == 0 {
		delete(p.pushHandlers, req.Type)

		err := p.control.SendDataExpectOk(fmt.Sprintf(
			"%s %s;", req.Type, req.Parameters))
		if err != nil {
			return fmt.Errorf("error stopping push notifications: %w",
				err)
		}
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

			pushType, pushData, err := getPushTypeAndData(b[:n])
			if err != nil {
				// TODO(bga): Log this.
				continue
			}

			p.m.Lock()

			tokenPushHandlerMap, ok := p.pushHandlers[pushType]
			if !ok {
				// TODO(bga): Log this.
				continue
			}

			for _, pushHandler := range tokenPushHandlerMap {
				pushHandler(pushData)
			}

			p.m.Unlock()
		}
	}

	p.quitChan = nil
}

func getPushTypeAndData(receivedData []byte) (string, string, error) {
	fields := bytes.Fields(receivedData)
	if len(fields) < 3 {
		return "", "", fmt.Errorf("invalid data received")
	}

	pushType := fmt.Sprintf("%s %s", string(fields[0]),
		string(fields[1]))

	pushDataBuilder := strings.Builder{}
	for i := 2; i < len(fields); i++ {
		if i != 2 {
			pushDataBuilder.WriteByte(' ')
		}
		pushDataBuilder.Write(fields[i])
	}

	pushData := pushDataBuilder.String()

	return pushType, pushData, nil
}
