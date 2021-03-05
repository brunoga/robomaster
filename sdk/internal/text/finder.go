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
	return internal.MatchIP(data.Addr.(*net.IPAddr).IP, filter)
}
