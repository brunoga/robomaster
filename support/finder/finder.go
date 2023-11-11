package finder

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/brunoga/unitybridge/support/logger"
)

const (
	ipBroadcastAddrPort = ":45678"
	listenerRemotePort  = ":56789"
)

// Finder provides an interface for finding a robot broadcasting its ip in
// the network.
type Finder struct {
	appID uint64

	l *logger.Logger

	m             sync.Mutex
	listeningConn *net.UDPConn
	quit          chan struct{}
	broadcasts    map[string]*Broadcast
}

// New returns a new Finder instance. If appID is zero, consider any robots
// detected in the network regardless of their pairing status or appID. If appID
// is non-zero, returns only robots with the given appID and that are not in
// pairing mode.
func New(appID uint64, l *logger.Logger) *Finder {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	return &Finder{
		appID:      appID,
		l:          l,
		quit:       nil,
		broadcasts: nil,
	}
}

// StartFinding starts listening for Robomaster broadcast messages in the
// network. It returns a non-nil error if it is already looking for robots.
func (f *Finder) StartFinding(ch chan<- *Broadcast) error {
	f.m.Lock()
	defer f.m.Unlock()

	if f.quit != nil {
		return fmt.Errorf("already finding")
	}

	f.quit = make(chan struct{})
	f.broadcasts = make(map[string]*Broadcast)

	udpAddr, err := net.ResolveUDPAddr("udp4", ipBroadcastAddrPort)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		return err
	}

	f.listeningConn = conn

	go func() {
		f.findLoop(ch)
	}()

	return nil
}

// StopFinding stops listening for Robomaster broadcast messages in the
// network. It returns a non-nil error if it is not currently looking for
// robots.
func (f *Finder) StopFinding() error {
	f.m.Lock()
	defer f.m.Unlock()

	if f.quit == nil {
		return fmt.Errorf("not finding")
	}

	close(f.quit)
	f.listeningConn.Close()

	return nil
}

// Find waits for a robot to broadcast its IP address in the network. It
// returns a non-nil error if no robot is found in the given timeout.
func (f *Finder) Find(timeout time.Duration) (*Broadcast, error) {
	ch := make(chan *Broadcast)

	err := f.StartFinding(ch)
	if err != nil {
		return nil, err
	}
	defer f.StopFinding()

	select {
	case broadcast := <-ch:
		return broadcast, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout")
	}
}

// SendACK sends an ACK message to the given IP address. This is used to
// acknowledge a pairing request.
func (f *Finder) SendACK(ip net.IP, appID uint64) {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, appID)

	udpAddr, err := net.ResolveUDPAddr("udp4", ip.String()+listenerRemotePort)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Write(buffer)
	if err != nil {
		return
	}
}

func (f *Finder) findLoop(ch chan<- *Broadcast) {
	f.l.Debug("Starting to look for robots")
	defer f.l.Debug("Stopped looking for robots")

	buf := make([]byte, 1024)

L:
	for {
		select {
		case <-f.quit:
			break L
		default:
			f.listeningConn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, addr, err := f.listeningConn.ReadFromUDP(buf)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok {
					if opErr.Timeout() {
						continue
					}
				}

				break L
			}

			broadcast, err := parseAndValidateBroadcast(buf[:n], addr)
			if err != nil {
				f.l.Warn("error parsing broadcast message", "err", err)
				continue
			}

			f.l.Debug("Received broadcast message", "broadcast", broadcast)

			if f.appID == 0 || (broadcast.AppId() == f.appID) {
				_, ok := f.broadcasts[broadcast.SourceIp().String()]
				if !ok {
					f.broadcasts[broadcast.SourceIp().String()] = broadcast
					ch <- broadcast
				}
			}
		}
	}

	f.quit = nil
	f.listeningConn = nil
}

func parseAndValidateBroadcast(buf []byte, addr net.Addr) (*Broadcast, error) {
	broadcastMessage, err := ParseBroadcast(buf)
	if err != nil {
		return nil, fmt.Errorf("error parsing broadcast message: %w", err)
	}

	// Get IP and make sure it is IPv4
	ip := net.IP(broadcastMessage.SourceIp()).To4()
	if ip == nil {
		return nil, fmt.Errorf("not an IPv4 address")
	}

	if !ip.Equal(addr.(*net.UDPAddr).IP) {
		return nil, fmt.Errorf("broadcast message source does not match reported IP")
	}

	return broadcastMessage, nil
}
