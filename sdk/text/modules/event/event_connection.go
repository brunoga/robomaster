package event

import (
	"fmt"
	"net"
	"sync"

	"github.com/brunoga/robomaster/sdk/text/modules/control"
	"github.com/brunoga/robomaster/sdk/text/modules/internal/notification"
)

const (
	eventAddrPort = ":40925"
)

type eventConnection struct {
	control *control.Control

	m    sync.Mutex
	conn net.Conn
}

func newEventConnection(control *control.Control) (notification.Connection, error) {
	if control == nil {
		return nil, fmt.Errorf("control must not be nil")
	}

	return &eventConnection{
		control,
		sync.Mutex{},
		nil,
	}, nil
}

func (e *eventConnection) Open() error {
	e.m.Lock()
	defer e.m.Unlock()

	if e.conn != nil {
		return fmt.Errorf("event connection already open")
	}

	ip, err := e.control.IP()
	if err != nil {
		return fmt.Errorf("error getting robot ip: %w", err)
	}

	eventAddr := ip.String() + eventAddrPort

	conn, err := net.Dial("tcp", eventAddr)
	if err != nil {
		return fmt.Errorf("error opening event connection: %w", err)
	}

	e.conn = conn

	return nil
}

func (e *eventConnection) Read(b []byte) (int, error) {
	e.m.Lock()
	defer e.m.Unlock()

	if e.conn == nil {
		return 0, fmt.Errorf("event connection is closed")
	}

	n, err := e.conn.Read(b)
	if err != nil {
		return 0, fmt.Errorf("error reading from event connection: %w", err)
	}

	return n, nil
}

func (e *eventConnection) Close() error {
	e.m.Lock()
	defer e.m.Unlock()

	err := e.conn.Close()

	e.conn = nil

	if err != nil {
		return fmt.Errorf("error closing event connection: %w", err)
	}

	return nil
}
