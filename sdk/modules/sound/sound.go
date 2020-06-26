package armor

import (
	"fmt"

	"github.com/brunoga/robomaster/sdk/modules"
	"github.com/brunoga/robomaster/sdk/modules/notification"
)

// EventAttribute represents sound attributes that can be monitored through
// event notifications.
type EventAttribute uint8

// Supported sound event attributes.
const (
	EventAttributeApplause EventAttribute = iota
)

type Sound struct {
	control *modules.Control
	event   *notification.Event
}

func New(control *modules.Control, event *notification.Event) *Sound {
	return &Sound{
		control,
		event,
	}
}

// StartEvent starts the event notification for the given eventType and
// eventAttribute. Updates will be sent to the given handler. Returns a token
// (used to stop notifications for the given handler) and a nil  error on
// success and a non-nil error on failure.
func (s *Sound) StartEvent(eventAttribute EventAttribute,
	handler notification.Handler) (int, error) {
	var eventAttributeStr string
	switch eventAttribute {
	case EventAttributeApplause:
		eventAttributeStr = "applause"
	default:
		return -1, fmt.Errorf("invalid sound event attribute")
	}

	token, err := s.event.StartListening("sound event",
		eventAttributeStr, "", handler)
	if err != nil {
		return -1, fmt.Errorf(
			"error starting to listen to sound event: %w", err)
	}

	return token, nil
}

// StopEvent stops event notifications to the handler represented by the given
// eventAttribute and token pair. Returns a nil error on success and a non-nil
// error on failure.
func (s *Sound) StopEvent(eventAttribute EventAttribute, token int) error {
	var eventAttributeStr string
	switch eventAttribute {
	case EventAttributeApplause:
		eventAttributeStr = "applause"
	default:
		return fmt.Errorf("invalid sound event attribute")
	}

	err := s.event.StopListening("sound event", eventAttributeStr, token)
	if err != nil {
		return fmt.Errorf(
			"error starting to listen to sound event: %w", err)
	}

	return nil
}
