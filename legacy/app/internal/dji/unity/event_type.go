package unity

// EventType represents the specific type of a Unit Event.
type EventType uint64

const (
	EventTypeSetValue           = EventType(0)
	EventTypeGetValue           = EventType(1)
	EventTypeGetAvailableValue  = EventType(2)
	EventTypePerformAction      = EventType(3)
	EventTypeStartListening     = EventType(4)
	EventTypeStopListening      = EventType(5)
	EventTypeActivation         = EventType(6)
	EventTypeLocalAlbum         = EventType(7)
	EventTypeFirmwareUpgrade    = EventType(8)
	EventTypeConnection         = EventType(100)
	EventTypeSecurity           = EventType(101)
	EventTypePrintLog           = EventType(200)
	EventTypeStartVideo         = EventType(300)
	EventTypeStopVideo          = EventType(301)
	EventTypeRender             = EventType(302)
	EventTypeGetNativeTexture   = EventType(303)
	EventTypeVideoTransferSpeed = EventType(304)
	EventTypeAudioDataRecv      = EventType(305)
	EventTypeVideoDataRecv      = EventType(306)
	EventTypeNativeFunctions    = EventType(500)
)

var eventTypeNameMap = map[EventType]string{
	0:   "EventTypeSetValue",
	1:   "EventTypeGetValue",
	2:   "EventTypeGetAvailableValue",
	3:   "EventTypePerformAction",
	4:   "EventTypeStartListening",
	5:   "EventTypeStopListening",
	6:   "EventTypeActivation",
	7:   "EventTypeLocalAlbum",
	8:   "EventTypeFirmwareUpgrade",
	100: "EventTypeConnection",
	101: "EventTypeSecurity",
	200: "EventTypePrintLog",
	300: "EventTypeStartVideo",
	301: "EventTypeStopVideo",
	302: "EventTypeRender",
	303: "EventTypeGetNativeTexture",
	304: "EventTypeVideoTransferSpeed",
	305: "EventTypeAudioDataRecv",
	306: "EventTypeVideoDataRecv",
	500: "EventTypeNativeFunctions",
}

// IsValidEventType checks if the given EventType is valid. It returns true if
// it is and false oherwise.
func IsValidEventType(eventType EventType) bool {
	_, ok := eventTypeNameMap[eventType]

	return ok
}

// EventTypeName returns the name associated with the given EventType. If it
// is not known, returns an empty string.
func EventTypeName(eventType EventType) string {
	eventTypeName, ok := eventTypeNameMap[eventType]
	if !ok {
		return ""
	}

	return eventTypeName
}

func AllEventTypes() []EventType {
	eventTypes := make([]EventType, 0, len(eventTypeNameMap))
	for eventType, _ := range eventTypeNameMap {
		eventTypes = append(eventTypes, eventType)
	}

	return eventTypes
}
