package finder

import (
	"net"

	"github.com/brunoga/robomaster/sdk/modules/finder"
)

type DataParserFunc func(net.IP, []byte) (*finder.Data, error)
