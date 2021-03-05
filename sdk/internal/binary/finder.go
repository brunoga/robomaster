package binary

import (
	"bytes"
	"net"

	"github.com/brunoga/robomaster/sdk/internal"
	"github.com/brunoga/robomaster/sdk/modules/finder"
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

func (f *Finder) filterFunc(data internal.FinderListenerData, filter finder.Filter) bool {
	if internal.MatchIP(data.Addr.(*net.IPAddr).IP, filter) {
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
		if bytes.Equal(data.Data, []byte(sn)) {
			return true
		}
	}

	return false
}
