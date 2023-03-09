package control

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/brunoga/robomaster/sdk/internal/binary/modules/control/event"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/command"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/message"
	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/support/logger"
	"github.com/brunoga/robomaster/sdk/types"

	binaryfinder "github.com/brunoga/robomaster/sdk/internal/binary/modules/finder"
)

const (
	// Port used to send commands to the robot.
	sdkPort = 20020

	// Port used to select the SDK connection.
	sdkProxyPort = 30030

	localSdkPortMin = 10100
	localSdkPortMax = 10500
)

// Control is the module used to control the robot. Responsible for setting up
// the connection, sending and receiving commands and keeping the connection
// alive.
type Control struct {
	// Used to identify the sender of messages. Seems to be basically constant
	// (9 and 6).
	host  byte
	index byte
	f     *binaryfinder.Finder

	eventManager *event.Manager

	pendingData []byte

	m    sync.Mutex
	conn net.Conn
}

// New returns a new Control using the given host and index bytes as
// identifiers.
func New(host, index byte, f *binaryfinder.Finder,
	l *logger.Logger) (*Control, error) {

	return &Control{
		host:         host,
		index:        index,
		f:            f,
		eventManager: event.NewManager(),
	}, nil
}

// Open opens the control connection to the robot at the given IP address. To do
// that it must do several things:
//
// 1 - Set the SDK connection to the robot using the SDK proxy port (UDP).
// 2 - Open the connection to the rrobot port using the requested network
// protocol.
// 3 - Enable SDK mode.
// 4 - Start a receive loop to wait for incoming data.
// 5 - Start a heart beat loop to keep the connection alive.
func (c *Control) Open(connMode types.ConnectionMode,
	connProto types.ConnectionProtocol, ip net.IP) error {
	c.m.Lock()

	if c.conn != nil {
		c.m.Unlock()
		return fmt.Errorf("client open: already open")
	}

	var (
		data []*finder.Data
		err  error
	)

	if ip == nil {
		data, err = c.f.Find(5*time.Second, func(_ net.IP, _ []byte) (bool, bool) {
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

	ip = data[0].IP()

	// Set SDK connection to be used.
	localAddr, err := c.setSDKConnection(connMode, connProto, ip)
	if err != nil {
		c.m.Unlock()
		return fmt.Errorf("client open: %w", err)
	}

	// Create a Dialer so we can set the local address to use (apparently this
	// is relevant to the SDK) and also set SO_REUSEADDR on the associated
	// socket.
	d := &net.Dialer{
		LocalAddr: localAddr,
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if err != nil {
					log.Printf("client open: %v", err)
				}
			})
		},
	}

	network := "udp4"
	if connProto == types.ConnectionProtocolTCP {
		network = "tcp4"
	}

	// Connect to the robot.
	conn, err := d.Dial(network, fmt.Sprintf("%s:%d", ip, sdkPort))
	if err != nil {
		c.m.Unlock()
		return fmt.Errorf("client open: %w", err)
	}

	c.conn = conn

	// Start receive loop.
	go c.receiveLoop()

	// Unlock earlier so we can send commands using the ourselves.
	c.m.Unlock()

	// Enable SDK mode.
	err = c.setSDKMode()
	if err != nil {
		return fmt.Errorf("client open: %w", err)
	}

	// Start heart beat loop.
	go c.heartBeatLoop()

	return nil
}

// Send sends the given message to the robot asynchornously. The provided
// callback (if non-nil) will be called whenever a response is received for the
// message.
func (c *Control) Send(message *message.Message, callback event.Callback) error {
	c.m.Lock()

	if c.conn == nil {
		c.m.Unlock()
		return fmt.Errorf("control send: not open")
	}

	// Write is thread safe so we unlock earlier.
	c.m.Unlock()

	_, err := c.conn.Write(message.Data())
	if err != nil {
		return fmt.Errorf("control send: %w", err)
	}

	if callback != nil {
		c.eventManager.Register(messageEventId(message), callback)
	}

	return nil
}

// SendSync sends the given message to the robot and waits for a response.
func (c *Control) SendSync(m *message.Message) (*message.Message, error) {
	// Used to wait for the response.
	var wg sync.WaitGroup

	var response *message.Message

	// Increment so we can be sure the Wait() below will block until the
	// callback is called.
	wg.Add(1)
	err := c.Send(m, func(m *message.Message) error {
		response = m

		// Allow SendSync to complete.
		wg.Done()

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("control send sync: %w", err)
	}

	// Wait for response.
	wg.Wait()

	return response, nil
}

func (c *Control) Close() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.conn == nil {
		return fmt.Errorf("client close: not open")
	}

	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("client close: %w", err)
	}

	c.conn = nil

	return nil
}

func messageEventId(message *message.Message) string {
	if message.IsResponse() {
		return fmt.Sprintf("%d:%d:%d:%d", message.Sender(), message.CmdSet(),
			message.CmdID(), message.SequenceID())
	}

	return fmt.Sprintf("%d:%d:%d:%d", message.Receiver(), message.CmdSet(),
		message.CmdID(), message.SequenceID())
}

func (c *Control) setSDKConnection(connMode types.ConnectionMode,
	connProto types.ConnectionProtocol, ip net.IP) (net.Addr, error) {
	localPort := uint16(rand.Intn(localSdkPortMax-localSdkPortMin+1) + localSdkPortMin)

	cmd := command.NewSetSDKConnectionRequest()
	cmd.SetConnection(byte(connMode)) // Infrastructure mode.
	cmd.SetHost(c.HostByte())
	cmd.SetProtocol(byte(connProto)) // TCP connection.
	cmd.SetPort(localPort)
	cmd.SetIP(net.IP{0, 0, 0, 0})

	m := message.New(c.HostByte(), protocol.HostToByte(9, 0), cmd)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, sdkProxyPort))
	if err != nil {
		return nil, fmt.Errorf("client set SDK connection: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(m.Data())
	if err != nil {
		return nil, fmt.Errorf("client set SDK connection: %w", err)
	}

	b := make([]byte, 128)
	n, err := conn.Read(b)
	if err != nil {
		return nil, fmt.Errorf("client set SDK connection: %w", err)
	}

	m, _, err = message.NewFromData(b[:n])
	if err != nil {
		return nil, fmt.Errorf("client set SDK connection: %w", err)
	}

	if connProto == types.ConnectionProtocolUDP {
		return &net.UDPAddr{
			IP:   m.Command().(*command.SetSDKConnectionResponse).ConfigIP(),
			Port: int(localPort),
		}, nil
	}

	return &net.TCPAddr{
		IP:   m.Command().(*command.SetSDKConnectionResponse).ConfigIP(),
		Port: int(localPort),
	}, nil
}

func (c *Control) setSDKMode() error {
	cmd := command.NewSetSDKModeRequest()
	cmd.SetEnable(true) // SDK mode.

	resp, err := c.SendSync(message.New(c.HostByte(), protocol.HostToByte(9, 0), cmd))
	if err != nil {
		return fmt.Errorf("client set SDK mode: %w", err)
	}

	if !resp.Command().(command.Response).Ok() {
		return fmt.Errorf("client set SDK mode: not ok")
	}

	return nil
}

func (c *Control) receiveLoop() {
	b := make([]byte, 256)
	for {
		n, err := c.conn.Read(b)
		if err != nil {
			break
		}

		if n == 0 {
			continue
		}

		c.pendingData = append(c.pendingData, b[:n]...)

		m, data, err := message.NewFromData(c.pendingData)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if m == nil {
			continue
		}

		c.eventManager.Trigger(messageEventId(m), m)

		c.pendingData = data
	}
}

func (c *Control) heartBeatLoop() {
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := c.Send(message.New(c.HostByte(), protocol.HostToByte(9, 0),
				command.NewSDKHeartBeatRequest()), nil)
			if err != nil {
				log.Printf("client heart beat: %v", err)
			}
			// TODO(bga): Add a close channel to be able top exit from here
			// cleanly.
		}
	}
}

func (c *Control) Host() byte {
	return c.host
}

func (c *Control) Index() byte {
	return c.index
}

func (c *Control) HostByte() byte {
	return protocol.HostToByte(c.host, c.index)
}
