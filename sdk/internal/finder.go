package internal

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/modules/robot"
)

const (
	readDeadline time.Duration = 500 * time.Millisecond
)

type FinderData struct {
	Addr net.Addr
	Data []byte
}

type FilterFunc func(net.Addr, []byte, finder.Filter) robot.Robot

// Finder is the generic implementation of the SDK Finder interface. It is used
// by both the binary and text protocols.
type Finder struct {
	udpAddrPort string
	timeout     time.Duration

	packetConn net.PacketConn

	filterFunc FilterFunc

	m       sync.Mutex
	finding bool
	robots  []robot.Robot
}

func NewFinder(udpAddrPort string, filterFunc FilterFunc) finder.Finder {
	return &Finder{
		udpAddrPort,
		0,
		nil,
		filterFunc,
		sync.Mutex{},
		false,
		nil,
	}
}

func (f *Finder) Find(filter finder.Filter, timeout time.Duration) error {
	f.m.Lock()
	defer f.m.Unlock()

	if f.finding {
		return fmt.Errorf("already looking for robots")
	}

	var err error
	f.packetConn, err = net.ListenPacket("udp4", f.udpAddrPort)
	if err != nil {
		return fmt.Errorf("error listening for packets: %w", err)
	}

	err = f.packetConn.SetReadDeadline(time.Now().Add(readDeadline))
	if err != nil {
		f.packetConn.Close()

		return fmt.Errorf("error setting read deadline: %w", err)
	}

	f.timeout = timeout
	f.finding = true

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
	var timerChan <-chan time.Time

	if f.timeout > 0 {
		ticker := time.NewTicker(f.timeout)
		defer ticker.Stop()

		timerChan = ticker.C
	}

L:
	for {
		select {
		case <-timerChan:
			break L
		default:
			buf := make([]byte, 1024)
			n, addr, err := f.packetConn.ReadFrom(buf)
			if err != nil {
				break L
			}

			if r := f.filterFunc(addr, buf[:n], filter); r != nil {
				f.m.Lock()
				f.robots = append(f.robots, r)
				f.m.Unlock()
			}
		}
	}

	f.m.Lock()
	f.finding = false
	f.m.Unlock()
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
