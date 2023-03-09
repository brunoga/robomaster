package control

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/support/logger"
	"github.com/brunoga/robomaster/sdk/types"

	textfinder "github.com/brunoga/robomaster/sdk/internal/text/modules/finder"
)

const (
	controlAddrPort = ":40923"
)

// Control handles sending commands to and receiving responses from a robot
// control connection.
type Control struct {
	logger *logger.Logger

	finder *textfinder.Finder

	m             sync.Mutex
	conn          net.Conn
	receiveBuffer []byte
}

// New returns a new Control instance with no associated ip. The given
// finder will be used to detect a robot broadcasting its ip in the
// network.
func New(f *textfinder.Finder, l *logger.Logger) (*Control, error) {
	return &Control{
		l,
		f,
		sync.Mutex{},
		nil,
		make([]byte, 512),
	}, nil
}

// Open tries to open the connection to the robot control port. Returns a nil
// error on success and a non-nil error on failure.
func (c *Control) Open(connMode types.ConnectionMode,
	connProto types.ConnectionProtocol, ip net.IP) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.conn != nil {
		return fmt.Errorf("connection already open")
	}

	var (
		data []*finder.Data
		err  error
	)

	if ip == nil {
		data, err = c.finder.Find(5*time.Second, func(_ net.IP, _ []byte) (bool, bool) {
			return true, false
		})
		if err != nil {
			return fmt.Errorf("error obtaining ip: %w", err)
		}
	} else {
		data = []*finder.Data{
			finder.NewData(ip, nil),
		}
	}

	if len(data) == 0 {
		return fmt.Errorf("no robot found")
	}

	addr := data[0].IP().String() + controlAddrPort

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

	c.logger.TRACE("Control: >>> %s", data)

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

	n, err := c.conn.Read(c.receiveBuffer)
	if err != nil {
		return "", fmt.Errorf("error reading data from control connection: %w",
			err)
	}

	c.logger.TRACE("Control: <<< %s", string(c.receiveBuffer[:n]))

	return string(bytes.TrimSpace(c.receiveBuffer[:n])), nil
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
		return fmt.Errorf("error sending and receiving data: %w", err)
	}

	if rcvData != "ok;" {
		return fmt.Errorf("error checking response: not ok")
	}

	return nil
}

// SendDataExpectOkAsync is a convenience method to send data and check we
// got an ok response back. The response is received asynchronously and
// checked if it is ok. If it is not ok, the information is logged. Any errors
// when receiving the response are also logged. This should be used whenever tbe
// latency of receiving a reply would interfere with the program. Returns a nil
// error on sending success and a non-nil error on failure.
func (c *Control) SendDataExpectOkAsync(data string) error {
	err := c.SendData(data)
	if err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}

	go func() {
		rcvData, err := c.ReceiveData()
		if err != nil {
			c.logger.ERROR("error sending data: %s", err)
		} else {
			if rcvData != "ok;" {
				c.logger.ERROR("%q -> %q\n", data, rcvData)
			}
		}
	}()

	return nil
}

// IP is a convenience function to get the robot ip. Returns the robot ip
// associated with this control instance and a nil error on success and a
// non-nil error on failure.
func (c *Control) IP() (net.IP, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.conn == nil {
		return nil, fmt.Errorf("connection not open")
	}

	return c.conn.RemoteAddr().(*net.TCPAddr).IP, nil
}
