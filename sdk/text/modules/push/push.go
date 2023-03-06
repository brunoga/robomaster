package push

import (
	"fmt"

	"github.com/brunoga/robomaster/sdk/text/modules/control"
	"github.com/brunoga/robomaster/sdk/text/modules/internal/notification"
)

// Push handles robot's push notification, starting/stopping monitoring individual
// events and sending them to registered PushHandlers.
type Push struct {
	*notification.Notification
}

type Handler notification.Handler

// NewPush returns a new Push instance. The control parameter is used to start
// stop the specific notification pushes.
func NewPush(control *control.Control) (*Push, error) {
	pushConnection, err := newPushConnection(control)
	if err != nil {
		return nil, fmt.Errorf("error creating push connection: %w", err)
	}

	notification, err := notification.New(control, pushConnection)
	if err != nil {
		return nil, fmt.Errorf("error creating push notification instance: %w",
			err)
	}

	return &Push{
		notification,
	}, nil
}
