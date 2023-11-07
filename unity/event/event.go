package event

// Event represents a Unity Bridge event with associated type and sub-type.
type Event struct {
	typ    Type
	subTyp uint32
}

// NewFromCode creates a new Event from the given event code.
func NewFromCode(eventCode uint64) *Event {
	return &Event{
		typ:    Type(eventCode >> 32),
		subTyp: uint32(eventCode & uint64(^uint(0))),
	}
}

// NewFromType creates a new Event from the given type.
func NewFromType(typ Type) *Event {
	return &Event{
		typ: typ,
	}
}

// NewFromTypeAndSubType creates a new Event from the given type and sub-type.
func NewFromTypeAndSubType(typ Type, subTyp uint32) *Event {
	return &Event{
		typ:    typ,
		subTyp: subTyp,
	}
}

// Code returns the event code for the event.
func (e *Event) Code() uint64 {
	return uint64(e.typ)<<32 | uint64(e.subTyp)
}

// Type returns the event type for the event.
func (e *Event) Type() Type {
	return e.typ
}

// SubType returns the event sub-type for the event.
func (e *Event) SubType() uint32 {
	return e.subTyp
}

// Reset resets the event type and sub-type.
func (e *Event) Reset(typ Type, subTyp uint32) {
	e.typ = typ
	e.subTyp = subTyp
}

// ResetSubType resets the event sub-type.
func (e *Event) ResetSubType(subTyp uint32) {
	e.subTyp = subTyp
}
