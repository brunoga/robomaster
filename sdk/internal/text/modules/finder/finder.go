package finder

import (
	"bytes"
	"fmt"
	"net"

	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/support/logger"

	internalfinder "github.com/brunoga/robomaster/sdk/internal/finder"
)

const (
	broadcastPort = 40926
)

// Finder provides an interface for finding a robot broadcasting its ip in
// the network.
type Finder struct {
	*internalfinder.Finder
}

// New returns a Finder instance with no associated ip.
func New(l *logger.Logger) *Finder {
	return &Finder{
		internalfinder.New(broadcastPort, dataParserFunc, l),
	}
}

func dataParserFunc(ip net.IP, buf []byte) (*finder.Data, error) {
	if !bytes.HasPrefix(buf, []byte("robot ip ")) {
		return nil, fmt.Errorf(
			"finder data parser function: invalid message received")
	}

	parsedIP := net.ParseIP(string(buf[9:]))
	if parsedIP == nil {
		return nil, fmt.Errorf(
			"finder data parser function: message ip is invalid")
	}

	if parsedIP.String() != ip.String() {
		return nil, fmt.Errorf(
			"finder data parser function: message ip does not match origin")
	}

	return finder.NewData(ip, nil), nil
}
