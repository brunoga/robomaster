package pairing

import (
	"net"
)

type EventType int

const (
	EventAdd    = EventType(0)
	EventRemove = EventType(1)
)

type Event struct {
	typ          EventType
	ip           net.IP
	hardwareAddr net.HardwareAddr
}

func NewEvent(typ EventType, ip net.IP,
	hardwareAddr net.HardwareAddr) *Event {
	return &Event{
		typ,
		ip,
		hardwareAddr,
	}
}

func (p *Event) Type() EventType {
	return p.typ
}

func (p *Event) IP() net.IP {
	return p.ip
}

func (p *Event) HardwareAddr() net.HardwareAddr {
	return p.hardwareAddr
}
