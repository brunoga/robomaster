package finder

import (
	"net"

	internalfinder "github.com/brunoga/robomaster/sdk/internal/finder"
	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/support/logger"
)

const (
	udpBroadcastPort = 40927
)

// Finder is the generic implementation of the SDK Finder interface. It is used
// by both the binary and text protocols.
type Finder struct {
	*internalfinder.Finder
}

func New(l *logger.Logger) *Finder {
	return &Finder{
		internalfinder.New(udpBroadcastPort, dataParserFunc, l),
	}
}

func dataParserFunc(ip net.IP, buf []byte) (*finder.Data, error) {
	return finder.NewData(ip, append([]byte{}, buf...)), nil
}
