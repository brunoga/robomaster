package udp

import (
	"fmt"
	"net"
	"sync"
)

type PortPair struct {
	localPort     int
	remotePort    int
	maxBufferSize int

	m          sync.Mutex
	packetConn net.PacketConn
	packetChan chan *Packet
	quitChan   chan struct{}

	wg sync.WaitGroup
}

func NewPortPair(localPort, remotePort, maxBufferSize int) *PortPair {
	return &PortPair{
		localPort,
		remotePort,
		maxBufferSize,
		sync.Mutex{},
		nil,
		nil,
		nil,
		sync.WaitGroup{},
	}
}

func (p *PortPair) Start() (<-chan *Packet, error) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.packetConn != nil {
		return nil, fmt.Errorf("already started")
	}

	packetConn, err := net.ListenPacket("udp", fmt.Sprintf(
		":%d", p.localPort))
	if err != nil {
		return nil, fmt.Errorf("error starting port pair: %w", err)
	}

	packetConn.(*net.UDPConn).SetReadBuffer(p.maxBufferSize)
	packetConn.(*net.UDPConn).SetWriteBuffer(p.maxBufferSize)

	p.packetConn = packetConn

	p.packetChan = make(chan *Packet)
	p.quitChan = make(chan struct{})

	p.wg.Add(1)
	go p.loop()

	return p.packetChan, nil
}

func (p *PortPair) Stop() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.packetConn == nil {
		return fmt.Errorf("not started")
	}

	close(p.quitChan)
	p.packetConn.Close()

	p.packetConn = nil

	p.packetChan = nil
	p.quitChan = nil

	p.wg.Wait()

	return nil
}

func (p *PortPair) Send(ip net.IP, data []byte) error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.packetConn == nil {
		return fmt.Errorf("not started")
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		ip.String(), p.remotePort))
	if err != nil {
		return fmt.Errorf("error resolving remote address: %w", err)
	}

	_, err = p.packetConn.WriteTo(data, remoteAddr)
	if err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}

	return nil
}

func (p *PortPair) loop() {
	fullStop := false

	buffer := make([]byte, p.maxBufferSize)
L:
	for {
		n, addr, err := p.packetConn.ReadFrom(buffer)
		if err != nil {
			fullStop = true
			break L
		}

		packet := NewPacket(addr.(*net.UDPAddr).IP, buffer[:n])

		select {
		case <-p.quitChan:
			break L
		case p.packetChan <- packet:
			// Do nothing.
		}
	}

	close(p.packetChan)
	p.wg.Done()

	if fullStop {
		p.Stop()
	}
}
