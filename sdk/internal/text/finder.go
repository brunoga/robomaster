package text

import (
	"net"

	"github.com/brunoga/robomaster/sdk/internal"
	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/modules/robot"
)

const (
	udpAddrPort = ":40926"
)

// Finder is the binary mode implementation of the SDK Finder interface. It
// currently only supports filtering by ip (key:ips, value:[]net.IP).
type Finder struct {
	finder.Finder
}

func NewFinder() finder.Finder {
	f := &Finder{}
	f.Finder = internal.NewFinder(udpAddrPort, f.filterFunc)

	return f
}

func (f *Finder) filterFunc(addr net.Addr, data []byte, filter finder.Filter) robot.Robot {
	// TODO(bga): Maybe validate that the IP matches the one in data.Data?
	ip := addr.(*net.UDPAddr).IP
	if internal.MatchIP(ip, filter) {
		return internal.NewRobot(ip, "")
	}

	return nil
}
