package event

import (
	"fmt"
	"sync"

	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/message"
)

// Manager is a simple event manager that allows to register and trigger events.
type Manager struct {
	m      sync.RWMutex
	events map[string]*Event
}

// NewManager returns a new event manager.
func NewManager() *Manager {
	return &Manager{
		events: make(map[string]*Event),
	}
}

// Register registers a new event with the given id and callback. The callback
// will be executed asynchronously when the event is triggered.
func (m *Manager) Register(id string, callback Callback) error {
	m.m.Lock()
	defer m.m.Unlock()

	_, ok := m.events[id]
	if ok {
		return fmt.Errorf("register event: event id %q already registered",
			id)
	}

	m.events[id] = NewEvent(id, Callback(callback))

	return nil
}

// Trigger triggers the event with the given id  by running the associated
// callback and passing the given data to it.
func (m *Manager) Trigger(id string, message *message.Message) error {
	m.m.Lock()
	defer m.m.Unlock()

	event, ok := m.events[id]
	if !ok {
		return fmt.Errorf("trigger event: event id %q not registered", id)
	}

	event.Trigger(message)

	delete(m.events, id)

	return nil
}

// NumPendingEvents returns the number of pending events.
func (m *Manager) NumPendingEvents() int {
	m.m.RLock()
	defer m.m.RUnlock()

	return len(m.events)
}
