package notification

import (
	"fmt"

	"github.com/brunoga/robomaster/sdk/modules/control"
)

// Event handles the robot event notification, starting/stopping monitoring
// individual events and sending them to registered EventHandlers.
type Event struct {
	*notification
}

// NewEvent returns a new Event instance. The control parameter is used to start
// stop the specific notification events.
func NewEvent(control *control.Control) (*Event, error) {
	eventConnection, err := newEventConnection(control)
	if err != nil {
		return nil, fmt.Errorf("error creating event connection: %w", err)
	}

	notification, err := newNotification(control, eventConnection)
	if err != nil {
		return nil, fmt.Errorf("error creating event notification instance: %w",
			err)
	}

	return &Event{
		notification,
	}, nil
}
