package binary

import (
	"net"
	"strings"

	"github.com/brunoga/robomaster/sdk/internal"
	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/modules/robot"
)

const (
	udpAddrPort = ":40927"
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

func matchSN(snToMatch string, filter finder.Filter) bool {
	if filter == nil {
		return true
	}

	maybeSNs := internal.GetFilterParameter("sns", filter)
	if maybeSNs == nil {
		return true
	}

	sns, ok := maybeSNs.([]string)
	if !ok {
		return true
	}

	for _, sn := range sns {
		if strings.Compare(snToMatch, sn) == 0 {
			return true
		}
	}

	return false
}

func (f *Finder) filterFunc(addr net.Addr, data []byte, filter finder.Filter) robot.Robot {
	ip := addr.(*net.UDPAddr).IP
	sn := string(data)
	if internal.MatchIP(ip, filter) || matchSN(sn, filter) {
		return internal.NewRobot(ip, sn)
	}

	return nil
}
