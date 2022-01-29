package pairing

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/udp"
)

const (
	listenerLocalPort  = 45678
	listenerRemotePort = 56789
	maxBufferSize      = 256
)

var (
	logger = log.New(os.Stdout, "PairingListener: ", log.LstdFlags)
)

type Listener struct {
	appId    uint64
	portPair *udp.PortPair

	m          sync.Mutex
	packetChan <-chan *udp.Packet
	eventChan  chan *Event
	clientMap  map[string]bool
}

func NewListener(appId uint64) *Listener {
	portPair := udp.NewPortPair(listenerLocalPort, listenerRemotePort,
		maxBufferSize)
	return &Listener{
		appId,
		portPair,
		sync.Mutex{},
		nil,
		nil,
		nil,
	}
}

func (l *Listener) Start() (<-chan *Event, error) {
	l.m.Lock()
	defer l.m.Unlock()

	packetChan, err := l.portPair.Start()
	if err != nil {
		return nil, err
	}

	logger.Printf("Starting on port %d.", listenerLocalPort)

	l.packetChan = packetChan
	l.eventChan = make(chan *Event)
	l.clientMap = make(map[string]bool)

	go l.loop()

	return l.eventChan, nil
}

func (l *Listener) Stop() error {
	l.m.Lock()
	defer l.m.Unlock()

	err := l.portPair.Stop()
	if err != nil {
		return err
	}

	logger.Printf("Stopping on port %d.", listenerLocalPort)

	l.packetChan = nil
	l.clientMap = nil

	return nil
}

func (l *Listener) SendACK(ip net.IP) error {
	logger.Printf("Sending ACK to %s:%d.\n", ip.String(), listenerRemotePort)

	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, l.appId)

	err := l.portPair.Send(ip, buffer)
	if err != nil {
		return fmt.Errorf("error sending ack: %w", err)
	}

	return nil
}

func (l *Listener) maybeGenerateEvent(ip net.IP, data []byte) *Event {
	bm, err := ParseBroadcastMessageData(data)
	if err != nil {
		logger.Printf("Error parsing broadcast message: %s.", err)
		return nil
	}

	l.m.Lock()
	defer l.m.Unlock()

	if l.clientMap[ip.String()] {
		if bm.AppId() != l.appId {
			l.clientMap[ip.String()] = false
			return NewEvent(EventRemove, bm.SourceIp(),
				bm.SourceMac())
		}
	} else {
		if bm.IsPairing() && bm.AppId() == l.appId {
			l.clientMap[ip.String()] = true
			return NewEvent(EventAdd, bm.SourceIp(), bm.SourceMac())
		}
	}

	return nil
}

func (l *Listener) loop() {
	for packet := range l.packetChan {
		event := l.maybeGenerateEvent(packet.IP(), packet.Data())
		if event == nil {
			continue
		}

		// TODO(bga): Add a quit channel.
		l.eventChan <- event
	}

	logger.Println("Existing read loop.")

	close(l.eventChan)
	l.eventChan = nil
}
