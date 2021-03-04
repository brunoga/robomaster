package text

import (
	"time"

	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/modules/robot"
)

// Finder is the binary mode implementation of the SDK Finder interface. It
// currently only supports filtering by ip (key:ips, value:[]net.IP).
type Finder struct {
}

func NewFinder() finder.Finder {
	// TODO(bga): Implement me.
	return &Finder{}
}

func (f *Finder) Find(filter finder.Filter, timeout time.Duration) {
	// TODO(bga): Implement me.
}

func (f *Finder) NumRobots() int {
	// TODO(bga): Implement me.
	return 0
}

func (f *Finder) Robot(n int) robot.Robot {
	// TODO(bga): Implement me.
	return nil
}
