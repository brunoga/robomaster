package armor

import (
	"fmt"
	"strconv"

	"github.com/brunoga/robomaster/sdk/internal/text/modules/event"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/notification"

	"github.com/brunoga/robomaster/sdk/internal/text/modules/control"
)

// EventAttribute represents armor attributes that can be monitored through
// event notifications.
type EventAttribute uint8

// Supported armor event attributes.
const (
	EventAttributeHit EventAttribute = iota
)

type Armor struct {
	control *control.Control
	event   *event.Event
}

func New(control *control.Control, event *event.Event) *Armor {
	return &Armor{
		control,
		event,
	}
}

func (a *Armor) SetSensitivity(sensitivity int) error {
	return a.control.SendDataExpectOk(fmt.Sprintf(
		"armor sensitivity %d;", sensitivity))
}

func (a *Armor) GetSensitivity() (int, error) {
	data, err := a.control.SendAndReceiveData("armor sensitivity ?;")
	if err != nil {
		return -1, fmt.Errorf("error sending sdk command: %w", err)
	}

	sensitivity, err := strconv.Atoi(data)
	if err != nil {
		return -1, fmt.Errorf("error parsing data: %w", err)
	}

	return sensitivity, nil
}

// StartEvent starts the event notification for the given eventType and
// eventAttribute. Updates will be sent to the given handler. Returns a token
// (used to stop notifications for the given handler) and a nil  error on
// success and a non-nil error on failure.
func (a *Armor) StartEvent(eventAttribute EventAttribute,
	handler notification.Handler) (int, error) {
	var eventAttributeStr string
	switch eventAttribute {
	case EventAttributeHit:
		eventAttributeStr = "hit"
	default:
		return -1, fmt.Errorf("invalid armor event attribute")
	}

	token, err := a.event.StartListening("armor event",
		eventAttributeStr, "", handler)
	if err != nil {
		return -1, fmt.Errorf(
			"error starting to listen to armor event: %w", err)
	}

	return token, nil
}

// StopEvent stops event notifications to the handler represented by the given
// eventAttribute and token pair. Returns a nil error on success and a non-nil
// error on failure.
func (a *Armor) StopEvent(eventAttribute EventAttribute, token int) error {
	var eventAttributeStr string
	switch eventAttribute {
	case EventAttributeHit:
		eventAttributeStr = "hit"
	default:
		return fmt.Errorf("invalid armor event attribute")
	}

	err := a.event.StopListening("armor event", eventAttributeStr, token)
	if err != nil {
		return fmt.Errorf(
			"error starting to listen to armor event: %w", err)
	}

	return nil
}
