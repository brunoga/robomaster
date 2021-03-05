package internal

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	readDeadline time.Duration = 500 * time.Millisecond
)

type FinderListenerData struct {
	Addr net.Addr
	Data []byte
}

type FinderListener struct {
	udpAddrPort string
	timeout     time.Duration

	packetConn net.PacketConn

	m           sync.Mutex
	controlChan chan struct{}
	readChan    chan FinderListenerData
}

func NewFinderListener(udpAddrPort string) *FinderListener {
	return &FinderListener{
		udpAddrPort,
		0,
		nil,
		sync.Mutex{},
		nil,
		nil,
	}
}
func (f *FinderListener) Start(timeout time.Duration) error {
	f.m.Lock()
	defer f.m.Unlock()

	if f.controlChan != nil {
		return fmt.Errorf("listener already started")
	}

	var err error
	f.packetConn, err = net.ListenPacket("udp4", f.udpAddrPort)
	if err != nil {
		return fmt.Errorf("error starting listner: %w", err)
	}

	err = f.packetConn.SetReadDeadline(time.Now().Add(readDeadline))
	if err != nil {
		f.packetConn.Close()

		return fmt.Errorf("error setting listener read deadline: %w", err)
	}

	f.timeout = timeout
	f.readChan = make(chan FinderListenerData)
	f.controlChan = make(chan struct{})

	go f.loop()

	return nil
}

func (f *FinderListener) ReadChannel() (<-chan FinderListenerData, error) {
	f.m.Lock()
	defer f.m.Unlock()

	if f.readChan == nil {
		return nil, fmt.Errorf("listener not started")
	}

	return f.readChan, nil
}

func (f *FinderListener) Stop() error {
	f.m.Lock()
	defer f.m.Unlock()

	if f.controlChan == nil {
		return fmt.Errorf("listener already stopped")
	}

	close(f.controlChan)
	f.controlChan = nil

	return nil
}

func (f *FinderListener) loop() {
	var timerChan <-chan time.Time

	if f.timeout > 0 {
		ticker := time.NewTicker(f.timeout)
		defer ticker.Stop()

		timerChan = ticker.C
	}

L:
	for {
		select {
		case <-f.controlChan:
			break L
		case <-timerChan:
			break L
		default:
			buf := make([]byte, 1024)
			n, addr, err := f.packetConn.ReadFrom(buf)
			if err != nil {
				break L
			}

			f.readChan <- FinderListenerData{addr, buf[:n]}
		}
	}

	close(f.readChan)

	f.m.Lock()
	close(f.readChan)
	f.readChan = nil
	f.controlChan = nil
	f.m.Unlock()
}
