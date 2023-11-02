package event

type Event struct {
	typ    Type
	subTyp uint32
}

func NewFromCode(eventCode uint64) *Event {
	return &Event{
		typ:    Type(eventCode >> 32),
		subTyp: uint32(eventCode & uint64(^uint(0))),
	}
}

func NewFromType(typ Type) *Event {
	return &Event{
		typ: typ,
	}
}

func NewFromTypeAndSubType(typ Type, subTyp uint32) *Event {
	return &Event{
		typ:    typ,
		subTyp: subTyp,
	}
}

func (e *Event) Code() uint64 {
	return uint64(e.typ)<<32 | uint64(e.subTyp)
}

func (e *Event) Type() Type {
	return e.typ
}

func (e *Event) SubType() uint32 {
	return e.subTyp
}

func (e *Event) Reset(typ Type, subTyp uint32) {
	e.typ = typ
	e.subTyp = subTyp
}

func (e *Event) ResetSubType(subTyp uint32) {
	e.subTyp = subTyp
}
