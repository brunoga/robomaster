package finder

import "net"

type Data struct {
	ip     net.IP
	serial []byte
}

func NewData(ip net.IP, serial []byte) *Data {
	return &Data{
		ip:     ip,
		serial: serial,
	}
}

func (d *Data) IP() net.IP {
	return d.ip
}

func (d *Data) Serial() []byte {
	return d.serial
}
