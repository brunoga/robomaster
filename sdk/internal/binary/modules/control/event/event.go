package event

import (
	"log"

	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/message"
)

// Event is a simple event that allows to register a name with a callback that
// will be executed when the event triggers.
type Event struct {
	id string

	callback Callback
}

// Callback is the type of the callback function that will be executed when an
// event is triggered.
type Callback func(message *message.Message) error

// NewEvent returns a new event with the given id and callback. The callback will
// be executed asynchronously when the event is triggered.
func NewEvent(id string, callback Callback) *Event {
	return &Event{
		id:       id,
		callback: callback,
	}
}

// Trigger triggers the event by running the associated callback and passing the
// given data to it. The callback will be executed asynchronously.
func (e *Event) Trigger(message *message.Message) {
	go func() {
		err := e.callback(message)
		if err != nil {
			log.Printf("trigger event: callback for event id %q return error %q",
				e.id, err)
		}
	}()
}
