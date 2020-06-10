package modules

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	controlAddrPort = ":40923"
)

// Control handles sending commands to and receiving responses from a robot
// control connection.
type Control struct {
	robotFinder *Finder

	m    sync.Mutex
	conn net.Conn
}

// NewControl returns a new Control instance with no associated ip. The given
// robotFinder will be used to detect a robot broadcasting its ip in the
// network.
func NewControl(robotFinder *Finder) *Control {
	return &Control{
		robotFinder,
		sync.Mutex{},
		nil,
	}
}

// Open tries to open the connection to the robot control port. Returns a nil
// error on success and a non-nil error on failure.
func (c *Control) Open() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.conn != nil {
		return fmt.Errorf("connection already open")
	}

	ip, err := c.robotFinder.GetOrFindIP(5 * time.Second)
	if err != nil {
		return fmt.Errorf("error obtaining ip: %w", err)
	}

	addr := ip.String() + controlAddrPort

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("error connecting to control port: %w", err)
	}

	c.conn = conn

	return nil
}

// Close tries to close the connection to the robot control port. Returns a nil
// error on success and a non-nil error on failure.
func (c *Control) Close() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.conn == nil {
		return fmt.Errorf("connection not open")
	}

	err := c.conn.Close()
	if err != nil {
		// TODO(bga): Maybe just log and continue?
		return fmt.Errorf("error closing control connection: %w", err)
	}

	c.conn = nil

	return nil
}

// SendData sends data to the control connection. The data should be the a
// fully-formed plain-text SDK command. Returns a nil error on success and a
// non-nil error on failure.
func (c *Control) SendData(data string) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.conn == nil {
		return fmt.Errorf("connection not open")
	}

	_, err := c.conn.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("error writting data to control connection: %w",
			err)
	}

	return nil
}

// ReceiveData reads any data available through the control connection. This
// is usually the result of executing a command sent with SendData().
// Returns the available data and a nil error on success and a non-nil error
// on failure
func (c *Control) ReceiveData() (string, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.conn == nil {
		return "", fmt.Errorf("connection not open")
	}

	buf := make([]byte, 512)

	_, err := c.conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("error reading data from control connection: %w",
			err)
	}

	return string(bytes.TrimSpace(buf)), nil
}

// SendAndReceiveData is a convenience method to send data and get the
// response data at once. Returns the received data and a nil error on
// success and a non-nil error on failure.
func (c *Control) SendAndReceiveData(data string) (string, error) {
	err := c.SendData(data)
	if err != nil {
		return "", fmt.Errorf("error sending data: %w", err)
	}

	rcvData, err := c.ReceiveData()
	if err != nil {
		return "", fmt.Errorf("error receiving data: %w", err)
	}

	return rcvData, nil
}

// SendDataExpectOk is a convenience method to send data and make sure we
// got an ok response back. Returns a nil error on success and a non-nil
// error on failure.
func (c *Control) SendDataExpectOk(data string) error {
	rcvData, err := c.SendAndReceiveData(data)
	if err != nil {
		fmt.Errorf("error sending and receiving data: %w", err)
	}

	if rcvData != "ok" {
		fmt.Errorf("error checking response: not ok")
	}

	return nil
}

// IP is a convenience function to get the robot ip. Returns the robot ip
// associated with this control instance and a nil error on success and a
// non-nil error on failure.
func (c *Control) IP() (net.IP, error) {
	return c.robotFinder.GetOrFindIP(5 * time.Second)
}