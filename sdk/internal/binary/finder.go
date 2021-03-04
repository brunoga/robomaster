package binary

import (
	"time"

	"github.com/brunoga/robomaster/sdk/finder"
	"github.com/brunoga/robomaster/sdk/robot"
)

// Finder is the binary mode implementation of the SDK Finder interface. It
// currently supports filtering by ip (key:ip, value:net.IP) and by serial
// number (key:sn, value:string).
type Finder struct {
}

// NewFinder returns a new binary mode Finder instance.
func NewFinder() finder.Finder {
	// TODO(bga): Implement me.
	return &Finder{}
}

func (f *Finder) Find(filter finder.Filter, timeout time.Duration) {
	// TODO(bga): Implement me.
}

func (f *Finder) NumRobots(filter finder.Filter) int {
	// TODO(bga): Implement me.
	return 0
}

func (f *Finder) Robots(filter finder.Filter) []robot.Robot {
	// TODO(bga): Implement me.
	return nil
}
