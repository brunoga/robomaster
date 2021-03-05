package text

import (
	"net"

	"github.com/brunoga/robomaster/sdk/internal"
	"github.com/brunoga/robomaster/sdk/modules/finder"
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

func (f *Finder) filterFunc(data internal.FinderListenerData, filter finder.Filter) bool {
	// TODO(bga): Maybe validate that the IP matches the one in data.Data?
	if filter == nil {
		return true
	}

	maybeIPs := internal.GetFilterParameter("ips", filter)
	if maybeIPs == nil {
		return true
	}

	ips, ok := maybeIPs.([]net.IP)
	if !ok {
		return true
	}

	for _, ip := range ips {
		if data.Addr.(*net.IPAddr).IP.Equal(ip) {
			return true
		}
	}

	return false
}
