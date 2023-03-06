package event

import (
	"fmt"

	"github.com/brunoga/robomaster/sdk/text/modules/control"
	"github.com/brunoga/robomaster/sdk/text/modules/internal/notification"
)

// Event handles the robot event notification, starting/stopping monitoring
// individual events and sending them to registered EventHandlers.
type Event struct {
	*notification.Notification
}

type Handler notification.Handler

// NewEvent returns a new Event instance. The control parameter is used to start
// stop the specific notification events.
func NewEvent(control *control.Control) (*Event, error) {
	eventConnection, err := newEventConnection(control)
	if err != nil {
		return nil, fmt.Errorf("error creating event connection: %w", err)
	}

	notification, err := notification.New(control, eventConnection)
	if err != nil {
		return nil, fmt.Errorf("error creating event notification instance: %w",
			err)
	}

	return &Event{
		notification,
	}, nil
}
