package push

import (
	"fmt"
	"net"
	"sync"

	"github.com/brunoga/robomaster/sdk/text/modules/control"
	"github.com/brunoga/robomaster/sdk/text/modules/internal/notification"
)

const (
	pushAddrPort = ":40924"
)

type pushConnection struct {
	control *control.Control

	m    sync.Mutex
	conn net.PacketConn
}

func newPushConnection(control *control.Control) (notification.Connection, error) {
	if control == nil {
		return nil, fmt.Errorf("control must not be nil")
	}

	return &pushConnection{
		control,
		sync.Mutex{},
		nil,
	}, nil
}

func (p *pushConnection) Open() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.conn != nil {
		return fmt.Errorf("push connection already open")
	}

	conn, err := net.ListenPacket("udp", pushAddrPort)
	if err != nil {
		return fmt.Errorf("error opening push connection: %w", err)
	}

	p.conn = conn

	return nil
}

func (p *pushConnection) Read(b []byte) (int, error) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.conn == nil {
		return 0, fmt.Errorf("push connection is closed")
	}

	for {
		n, addr, err := p.conn.ReadFrom(b)
		if err != nil {
			return 0, fmt.Errorf("error reading from push connection: %w", err)
		}

		robotIP, err := p.control.IP()
		if err != nil {
			return 0, fmt.Errorf("error obtaining robot ip: %w", err)
		}

		if robotIP.String() != addr.(*net.UDPAddr).IP.String() {
			// Got push notification from an unexpected ip. Ignore it.
			continue
		}

		return n, nil
	}
}

func (p *pushConnection) Close() error {
	p.m.Lock()
	defer p.m.Unlock()

	err := p.conn.Close()

	p.conn = nil

	if err != nil {
		return fmt.Errorf("error closing push connection: %w", err)
	}

	return nil
}
