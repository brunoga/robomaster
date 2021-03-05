package internal

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/modules/robot"
)

// Finder is the generic implementation of the SDK Finder interface. It is used
// by both the binary and text protocols.
type Finder struct {
	listener   *FinderListener
	filterFunc func(FinderListenerData, finder.Filter) bool

	m      sync.Mutex
	robots []robot.Robot
}

func NewFinder(udpAddrPort string,
	filterFunc func(FinderListenerData, finder.Filter) bool) finder.Finder {
	return &Finder{
		NewFinderListener(udpAddrPort),
		filterFunc,
		sync.Mutex{},
		nil,
	}
}

func (f *Finder) Find(filter finder.Filter, timeout time.Duration) error {
	err := f.listener.Start(timeout)
	if err != nil {
		return fmt.Errorf("error starting to listen for robots: %w", err)
	}

	go f.findLoop(filter)

	return nil
}

func (f *Finder) NumRobots() int {
	f.m.Lock()
	defer f.m.Unlock()

	return len(f.robots)
}

func (f *Finder) Robot(n int) robot.Robot {
	f.m.Lock()
	defer f.m.Unlock()

	if n > len(f.robots) {
		return nil
	}

	return f.robots[n]
}

func (f *Finder) findLoop(filter finder.Filter) {
	readChannel, err := f.listener.ReadChannel()
	if err != nil {
		// This should never happen.
		panic(err)
	}

	for listenerData := range readChannel {
		if f.filterFunc(listenerData, filter) {
			var r robot.Robot
			f.m.Lock()
			f.robots = append(f.robots, r)
			f.m.Unlock()
		}
	}
}

// GetFilterParameter returns the value (as an interface{}) in the given filter
// associated with the given key. If key is not found, returns nil.
func GetFilterParameter(key string, filter finder.Filter) interface{} {
	v, ok := filter[key]
	if !ok {
		return nil
	}

	return v
}

func MatchIP(ipToMatch net.IP, filter finder.Filter) bool {
	if filter == nil {
		return true
	}

	maybeIPs := GetFilterParameter("ips", filter)
	if maybeIPs == nil {
		return true
	}

	ips, ok := maybeIPs.([]net.IP)
	if !ok {
		return true
	}

	for _, ip := range ips {
		if ipToMatch.Equal(ip) {
			return true
		}
	}

	return false
}
