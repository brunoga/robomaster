package callback

import (
	"C"
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge/wrapper/callback"
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

		l = l.WithGroup("event_callback_manager")

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
	eventTypeCode := getEventType(eventCode)

	m.m.Lock()
	defer m.m.Unlock()

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
	eventTypeCode := getEventType(eventCode)

	m.m.RLock()
	c, ok := m.eventCodeCallback[eventTypeCode]
	if !ok {
		m.m.RUnlock()
		return fmt.Errorf("no handlers for event type code %d", eventTypeCode)
	}
	m.m.RUnlock()

	// Make a copy of the data so we can:
	//
	// 1. Move the data out of the C side of things into the Go realm (so we)
	//    can benefit of our garbage collector.
	// 2. Allow us doing things in goroutines (data will not disappear under
	//	  us).
	// 3. We can also return faster to the C side of things which unblocks the
	//    Unity Bridge code and might also help speeding thing up.
	//
	// Note that for the Wine version, we end up doing 2 copies (one here and
	// another one when sending the data down the pipe from Wine to the Linux
	// side). Considering even for video frames there is not much data to copy
	// (only 6Mb or so in the worst case scenario), this is not a big deal.
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	go c(eventCode, dataCopy, tag)

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
