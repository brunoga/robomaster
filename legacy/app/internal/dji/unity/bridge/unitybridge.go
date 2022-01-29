package bridge

import (
	"fmt"
	"log"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/wrapper"
	"git.bug-br.org.br/bga/robomasters1/app/internal/support/callbacks"
)

// EventHandler is the required prototype for functiomns that want to process
// Unity events. Implementations are required to call wg.Done() before they
// return.
type EventHandler func(event *unity.Event, data []byte, tag uint64,
	wg *sync.WaitGroup)

// unityBridge is a frontend to DJI's Unity Bridge code. The type is not
// exported because it is currently implemented as a singleton.
type unityBridge struct {
	m     sync.Mutex
	setup bool

	cbs *callbacks.Callbacks
	w   wrapper.Wrapper
}

var (
	// Singleton instance.
	instance *unityBridge
)

func init() {
	w, err := wrapper.New("./unitybridge.dll")
	if err != nil {
		panic(err)
	}

	// Creates the singleton instance.
	instance = &unityBridge{
		sync.Mutex{},
		false,
		callbacks.New("UnityBridge", nil, nil),
		w,
	}
}

// Setup creates and initializes the underlying Unity Bridge. It returns a nil
// error on success and a non-nil error on failure.
func Setup(name string, debuggable bool) error {
	instance.m.Lock()
	defer instance.m.Unlock()

	if instance.setup {
		return fmt.Errorf("bridge already setup")
	}

	// Creates the underlying Unity Bridge.
	instance.w.CreateUnityBridge(name, debuggable)

	// Register the callback to all known events.
	instance.registerCallback()

	if !instance.w.UnityBridgeInitialize() {
		// Something went wrong so we bail out.
		instance.w.DestroyUnityBridge()
		return fmt.Errorf("bridge initialization failed")
	}

	instance.setup = true

	return nil
}

// Teardown uninitializes and destroys the underlying Unity Bridge. It returns
// a nil error on success and a non-nil error on failure.
func Teardown() error {
	instance.m.Lock()
	defer instance.m.Unlock()

	if !instance.setup {
		return fmt.Errorf("bridge not setup")
	}

	// Unregister the callbacks to all known events.
	instance.unregisterCallback()

	instance.w.UnityBridgeUninitialize()
	instance.w.DestroyUnityBridge()

	instance.setup = false

	return nil
}

// IsSetup returns true if the underlying Unity Bridge support was setup and
// false otherwise.
func IsSetup() bool {
	instance.m.Lock()
	defer instance.m.Unlock()

	return instance.setup
}

// Instance returns a pointer to the unityBridge singleton.
func Instance() *unityBridge {
	return instance
}

// AddEventHandler adds an event handler for the given event.
func (b *unityBridge) AddEventHandler(eventType unity.EventType,
	eventHandler EventHandler) (uint64, error) {
	if !unity.IsValidEventType(eventType) {
		return 0, fmt.Errorf("invalid event type")
	}

	tag, err := b.cbs.AddContinuous(callbacks.Key(eventType), eventHandler)
	if err != nil {
		return 0, err
	}

	return uint64(tag), nil
}

// RemoveEventHandler removes the event handler at the given index for the
// given event.
func (b *unityBridge) RemoveEventHandler(eventType unity.EventType,
	tag uint64) error {
	if !unity.IsValidEventType(eventType) {
		return fmt.Errorf("invalid event type")
	}

	return b.cbs.Remove(callbacks.Key(eventType), callbacks.Tag(tag))
}

// SendEvent sends a unity event through the underlying Unity Bridge. It can
// accept one, two or three parameters. The first one is the event itself and
// must be a *unity.Event. The second one is the data to send associated with
// the event and can be a []byte, a string or a uint64. The third one is the
// tag number associated with the event (which is used to disambiguate events)
// and must be a uint64.
func (b *unityBridge) SendEvent(params ...interface{}) error {
	if len(params) < 1 || len(params) > 3 {
		return fmt.Errorf("1, 2 or 3 parameters are required")
	}

	event, ok := params[0].(*unity.Event)
	if !ok {
		return fmt.Errorf("event (first) parameter must be a " +
			"*unity.Event")
	}

	dataType := 0
	var data interface{} = nil
	if len(params) > 1 {
		switch params[1].(type) {
		case []byte:
			// Do nothing.
		case string:
			dataType = 1
		case uint64:
			dataType = 2
		default:
			return fmt.Errorf("data (second) parameter must be " +
				"[]byte, string or uint64")
		}
		data = params[1]
	}

	var tag uint64 = 0
	if len(params) > 2 {
		tag, ok = params[2].(uint64)
		if !ok {
			return fmt.Errorf("tag (third) parameter must be " +
				"uint64")
		}
	}

	switch dataType {
	case 0:
		if data != nil {
			instance.w.UnitySendEvent(event.Code(), data.([]byte),
				tag)
		} else {
			instance.w.UnitySendEvent(event.Code(), nil, tag)
		}
	case 1:
		instance.w.UnitySendEventWithString(event.Code(), data.(string),
			tag)
	case 2:
		instance.w.UnitySendEventWithNumber(event.Code(), data.(uint64),
			tag)
	}

	return nil
}

func (b *unityBridge) registerCallback() {
	for _, eventType := range unity.AllEventTypes() {
		event := unity.NewEvent(eventType)
		instance.w.UnitySetEventCallback(event.Code(),
			b.unityEventCallback)
	}
}

func (b *unityBridge) unregisterCallback() {
	for _, eventType := range unity.AllEventTypes() {
		event := unity.NewEvent(eventType)
		instance.w.UnitySetEventCallback(event.Code(), nil)
	}
}

func (b *unityBridge) unityEventCallback(eventCode uint64, data []byte,
	tag uint64) {
	event := unity.NewEventFromCode(eventCode)
	if event == nil {
		log.Printf("Unknown event with code %d.\n", eventCode)
		return
	}

	eventHandlers, err := b.cbs.CallbacksForKey(callbacks.Key(event.Type()))
	if err != nil {
		log.Printf("No event handlers for %q\n",
			unity.EventTypeName(event.Type()))
		return
	}

	wg := sync.WaitGroup{}

	for _, handler := range eventHandlers {
		wg.Add(1)
		go handler.(EventHandler)(event, data, tag, &wg)
	}

	wg.Wait()
}
