package callback

import (
	"C"
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/wrapper/callback"
)

var (
	once     sync.Once
	instance *Manager
)

// Manager is a singleton that manages callbacks for events sent from Unity.
type Manager struct {
	l *logger.Logger

	m                 sync.RWMutex
	eventCodeCallback map[uint32]callback.Callback
}

// NewManager returns a new Manager instance. It is lazely allocated singleton
// which meas that the first time it is called it will allocate a Manager with
// the given logger but subsequent calls will just return the originaly created
// instance (and the logger parameter is be ignored).
func NewManager(l *logger.Logger) *Manager {
	once.Do(func() {
		if l == nil {
			l = logger.New(slog.LevelError)
		}

		instance = &Manager{
			l:                 l,
			eventCodeCallback: make(map[uint32]callback.Callback),
		}
	})

	return instance
}

// Set sets the callback for the given event type code. If the callback is nil,
// the callback for the given event type code is removed.
func (m *Manager) Set(eventCode uint64, c callback.Callback) error {
	m.m.Lock()
	defer m.m.Unlock()

	eventTypeCode := getEventType(eventCode)

	_, ok := m.eventCodeCallback[eventTypeCode]
	if !ok {
		if c == nil {
			return fmt.Errorf("no callback for event type %d found", eventTypeCode)
		}

		m.eventCodeCallback[eventTypeCode] = c
	} else {
		if c == nil {
			delete(m.eventCodeCallback, eventTypeCode)
		} else {
			return fmt.Errorf("callback for event type code %d already set",
				eventTypeCode)
		}
	}

	return nil
}

// Run runs the callback for the given event code.
func (m *Manager) Run(eventCode uint64, data []byte, tag uint64) error {
	m.m.RLock()
	defer m.m.RUnlock()

	eventTypeCode := getEventType(eventCode)

	c, ok := m.eventCodeCallback[eventTypeCode]
	if !ok {
		return fmt.Errorf("no handlers for event type code %d", eventTypeCode)
	}

	// TODO(bga): Maybe do this in a goroutine? If we do, we must copy data as
	// it is backed by a C array that is freed after Run() returns.
	c(eventCode, data, tag)

	return nil
}

//export eventCallbackGo
func eventCallbackGo(eventCode uint64, data []byte, tag uint64) {
	err := instance.Run(eventCode, data, tag)
	if err != nil {
		log.Printf("error running event callback: %s\n", err)
	}
}

func getEventType(eventCode uint64) uint32 {
	return uint32(eventCode >> 32)
}
