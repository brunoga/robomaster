package udp

import (
	"net"
)

type Packet struct {
	ip   net.IP
	data []byte
}

func NewPacket(ip net.IP, data []byte) *Packet {
	return &Packet{
		ip,
		data,
	}
}

func (p *Packet) IP() net.IP {
	return p.ip
}

func (p *Packet) Data() []byte {
	return p.data
}
