package modules

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	ipBroadcastAddrPort = ":40926"
)

// RobotFinder provides an interface for finding a robot broadcasting its ip in
// the network.
type RobotFinder struct {
	m  sync.Mutex
	ip net.IP
}

// NewRobotFinder returns a RobotFinder instance with no associated ip.
func NewRobotFinder() *RobotFinder {
	return &RobotFinder{
		sync.Mutex{},
		nil,
	}
}

// GetOrFindIP returns the ip of a robot if it is already know or tries to
// detect a robot broadcasting its ip in the network. The search will go on
// until a robot is detected or a timeout happens. Return the robot ip and a nil
// error on success and a non-nil error on failure.
func (r *RobotFinder) GetOrFindIP(timeout time.Duration) (net.IP, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if r.ip == nil {
		ip, err := findRobotIP(timeout)
		if err != nil {
			return nil, fmt.Errorf("error finding robot ip: %w", err)
		}

		r.ip = ip
	}

	return r.ip, nil
}

func findRobotIP(timeout time.Duration) (net.IP, error) {
	packetConn, err := net.ListenPacket("udp4", ipBroadcastAddrPort)
	if err != nil {
		return nil, fmt.Errorf("error starting packet listner: %w", err)
	}
	defer packetConn.Close()

	buf := make([]byte, 1024)

	err = packetConn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, fmt.Errorf("error setting deadline: %w", err)
	}

	n, addr, err := packetConn.ReadFrom(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading packet: %w", err)
	}

	ip, err := parseAndValidateMessage(buf[:n], addr)
	if err != nil {
		return nil, fmt.Errorf("error validating message: %w", err)
	}

	return ip, nil
}

func parseAndValidateMessage(buf []byte, addr net.Addr) (net.IP, error) {
	if !bytes.HasPrefix(buf, []byte("robot ip ")) {
		return nil, fmt.Errorf("received message is invalid")
	}

	ip := net.ParseIP(string(buf[9:]))
	if ip == nil {
		return nil, fmt.Errorf("message ip is invalid")
	}

	if ip.String() != addr.(*net.UDPAddr).IP.String() {
		return nil, fmt.Errorf("message ip does not match origin")
	}

	return ip, nil
}
