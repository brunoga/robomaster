package unity

import (
	"fmt"
)

// Event represents a Unity Event used for communication with the Robomaster S1.
type Event struct {
	typ    EventType
	subTyp uint64
}

// NewEvent creates a new Event with the given type and 0 as subtype. Returns a
// pointer to the created Event on success or nil on failure (usually if the
// specific event type is unknown).
func NewEvent(typ EventType) *Event {
	return NewEventWithSubType(typ, 0)
}

// NewEventWithSubType creates a new Event with the given type and subtype.
// Returns a pointer to the created Event on success or nil on failure (usually
// if the specific event type is unknown).
func NewEventWithSubType(typ EventType, subTyp uint64) *Event {
	if !IsValidEventType(typ) {
		return nil
	}

	return &Event{
		typ,
		subTyp,
	}
}

// NewEventFromCode creates a new Event by parsing the given code. Returns a
// pointer to the created Event on success or nil on failure (usually if the
// parsed event type is unknown).
func NewEventFromCode(code uint64) *Event {
	typ := EventType(code >> 32)
	if !IsValidEventType(typ) {
		return nil
	}

	subTyp := code & uint64(^uint32(0))

	return &Event{
		typ,
		subTyp,
	}
}

// Type returns the type associated with this event.
func (e *Event) Type() EventType {
	return e.typ
}

// SubType returns the subtype associated with this event.
func (e *Event) SubType() uint64 {
	return e.subTyp
}

// Code returns the uint64 code associated with this Event.
func (e *Event) Code() uint64 {
	return (uint64(e.typ) << 32) | e.subTyp
}

// Reset resets this Event to have the given type and subtype.
func (e *Event) Reset(typ EventType, subTyp uint64) {
	e.typ = typ
	e.subTyp = subTyp
}

/// ResetSubType resets this event to have the given subtype.
func (e *Event) ResetSubType(subTyp uint64) {
	e.subTyp = subTyp
}

// String returns this Event as a formated string. This implements the
// fmt.Stringer interface.
func (e *Event) String() string {
	return fmt.Sprintf("Event: Type=%s (%d), SubType=%d",
		EventTypeName(e.typ), e.typ, e.subTyp)
}
