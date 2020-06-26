package notification

import (
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules"
)

// Push handles robot's push notification, starting/stopping monitoring individual
// events and sending them to registered PushHandlers.
type Push struct {
	*notification
}

// NewPush returns a new Push instance. The control parameter is used to start
// stop the specific notification pushes.
func NewPush(control *modules.Control) (*Push, error) {
	pushConnection, err := newPushConnection(control)
	if err != nil {
		return nil, fmt.Errorf("error creating push connection: %w", err)
	}

	notification, err := newNotification(control, pushConnection)
	if err != nil {
		return nil, fmt.Errorf("error creating push notification instance: %w",
			err)
	}

	return &Push{
		notification,
	}, nil
}
