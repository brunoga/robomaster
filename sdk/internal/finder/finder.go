package finder

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/brunoga/robomaster/sdk/modules/finder"
	"github.com/brunoga/robomaster/sdk/support/logger"
)

// Finder is the generic implementation of the SDK Finder interface. It is used
// by both the binary and text protocols.
type Finder struct {
	l              *logger.Logger
	port           int
	dataParserFunc DataParserFunc
	packetConn     net.PacketConn

	m sync.Mutex
}

func New(port int, dataParserFunc DataParserFunc,
	l *logger.Logger) *Finder {
	return &Finder{
		l:              l,
		port:           port,
		dataParserFunc: dataParserFunc,
		packetConn:     nil,
	}
}

func (f *Finder) Find(timeout time.Duration,
	finderFunc finder.Func) ([]*finder.Data, error) {
	if !f.m.TryLock() {
		return nil, fmt.Errorf("finder find: already looking for robots")
	}
	defer f.m.Unlock()

	var err error
	f.packetConn, err = net.ListenPacket("udp", fmt.Sprintf(":%d", f.port))
	if err != nil {
		return nil, fmt.Errorf("finder find: error listening for packets: %w", err)
	}

	return f.findLoop(timeout, finderFunc)
}

func (f *Finder) findLoop(timeout time.Duration,
	finderFunc finder.Func) ([]*finder.Data, error) {
	go time.AfterFunc(timeout, func() {
		f.packetConn.Close()
	})

	ipsMap := map[string]bool{}
	datas := []*finder.Data{}

	buf := make([]byte, 50)
	for {
		n, addr, err := f.packetConn.ReadFrom(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			} else if err, ok := err.(net.Error); ok && err.Timeout() {
				break
			}

			return nil, fmt.Errorf("finder find: error reading packet: %w", err)
		}

		ip := addr.(*net.UDPAddr).IP

		data, err := f.dataParserFunc(ip, buf[:n])
		if err != nil {
			return nil, fmt.Errorf("finder find: error parsing response: %w", err)
		}

		_, ok := ipsMap[ip.String()]
		if !ok {
			ipsMap[ip.String()] = true
			if finderFunc == nil {
				datas = append(datas, data)
				continue
			} else {
				add, cont := finderFunc(ip, buf[:n])
				if add {
					datas = append(datas, data)
				}
				if !cont {
					break
				}
			}

			break
		}
	}

	return datas, nil
}
