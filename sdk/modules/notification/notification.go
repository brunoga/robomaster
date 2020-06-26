package notification

import (
	"bytes"
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules"
	"strings"
	"sync"
)

type Handler func(string)

type notification struct {
	control    *modules.Control
	connection connection

	m        sync.RWMutex
	quitChan chan struct{}
	handlers map[string]map[string]map[int]Handler
}

func newNotification(control *modules.Control, connection connection) (*notification, error) {
	if control == nil {
		return nil, fmt.Errorf("control must not be nil")
	}

	if connection == nil {
		return nil, fmt.Errorf("connection must not be nil")
	}

	return &notification{
		control,
		connection,
		sync.RWMutex{},
		nil,
		make(map[string]map[string]map[int]Handler),
	}, nil
}

func (n *notification) StartListening(notificationType, notificationAttribute,
	notificationParameters string, handler Handler) (int, error) {
	n.m.Lock()
	defer n.m.Unlock()

	attributeTokenMap, ok := n.handlers[notificationType]
	if !ok {
		attributeTokenMap = make(map[string]map[int]Handler)
	}

	startTracking := false

	tokenHandlerMap, ok := attributeTokenMap[notificationAttribute]
	if !ok {
		startTracking = true
		tokenHandlerMap = make(map[int]Handler)
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
		return -1, fmt.Errorf("error obtaining notification handler token")
	}

	if startTracking {
		// This eventType/eventAttribute pair is not being tracked yet. Start
		// tracking.
		var command string
		if notificationParameters == "" {
			command = fmt.Sprintf("%s %s on;", notificationType,
				notificationAttribute)
		} else {
			command = fmt.Sprintf("%s %s on %s;", notificationType,
				notificationAttribute, notificationParameters)
		}

		err := n.control.SendDataExpectOk(command)
		if err != nil {
			return -1, fmt.Errorf("error listening for notification: %w",
				err)
		}
	}

	if len(n.handlers) == 0 {
		go n.loop()
	}

	tokenHandlerMap[token] = handler
	attributeTokenMap[notificationAttribute] = tokenHandlerMap
	n.handlers[notificationType] = attributeTokenMap

	return token, nil
}

func (n *notification) StopListening(notificationType,
	notificationAttribute string, token int) error {
	n.m.Lock()
	defer n.m.Unlock()

	attributeTokenMap, ok := n.handlers[notificationType]
	if !ok {
		return fmt.Errorf("no handlers for notification type")
	}

	tokenHandlerMap, ok := attributeTokenMap[notificationAttribute]
	if !ok {
		return fmt.Errorf("no handlers for notification attribute")
	}

	_, ok = tokenHandlerMap[token]
	if !ok {
		return fmt.Errorf("token does not match notification type")
	}

	delete(tokenHandlerMap, token)

	if len(tokenHandlerMap) == 0 {
		// This notificationType/notificationAttribute pair is not being tracked
		// anymore. Stop listening to it.
		err := n.control.SendDataExpectOk(fmt.Sprintf(
			"%s %s off;", notificationType, notificationAttribute))
		if err != nil {
			return fmt.Errorf("error stopping notification: %w",
				err)
		}

		delete(attributeTokenMap, notificationAttribute)
	}

	if len(attributeTokenMap) == 0 {
		delete(n.handlers, notificationType)
	}

	if len(n.handlers) == 0 {
		close(n.quitChan)
	}

	return nil
}

func (n *notification) loop() {
	n.quitChan = make(chan struct{})

	err := n.connection.Open()
	if err != nil {
		// TODO(bga): Log this.
		return
	}
	defer n.connection.Close()

	b := make([]byte, 512)
L:
	for {
		select {
		case <-n.quitChan:
			break L
		default:
			nr, err := n.connection.Read(b)
			if err != nil {
				// TODO(bga): Log this.
				break L
			}

			notificationType, notificationAttribute, notificationData, err :=
				getNotificationTypeAttributeAndData(b[:nr])
			if err != nil {
				// TODO(bga): Log this.
				continue
			}

			n.m.RLock()

			attributeTokenMap, ok := n.handlers[notificationType]
			if !ok {
				// TODO(bga): Log this.
				continue
			}

			tokenPushHandlerMap, ok := attributeTokenMap[notificationAttribute]
			if !ok {
				// TODO(bga): Log this.
				continue
			}

			for _, handler := range tokenPushHandlerMap {
				handler(notificationData)
			}

			n.m.RUnlock()
		}
	}

	n.quitChan = nil
}

func getNotificationTypeAttributeAndData(
	receivedData []byte) (string, string, string, error) {
	fields := bytes.Fields(receivedData)
	if len(fields) < 4 {
		return "", "", "", fmt.Errorf("invalid data received")
	}

	notificationType := fmt.Sprintf("%s %s", string(fields[0]),
		string(fields[1]))

	notificationAttribute := string(fields[2])

	notificationDataBuilder := strings.Builder{}
	for i := 3; i < len(fields); i++ {
		if i != 3 {
			notificationDataBuilder.WriteByte(' ')
		}
		notificationDataBuilder.Write(fields[i])
	}

	notificationData := notificationDataBuilder.String()

	return notificationType, notificationAttribute, notificationData, nil
}
