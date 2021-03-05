package internal

import (
	"net"

	"github.com/brunoga/robomaster/sdk/modules/robot"
)

type Robot struct {
	ip net.IP
	sn string
}

func NewRobot(ip net.IP, sn string) robot.Robot {
	return &Robot{
		ip,
		sn,
	}
}

func (r *Robot) IP() net.IP {
	return r.ip
}

func (r *Robot) SN() string {
	return r.sn
}
