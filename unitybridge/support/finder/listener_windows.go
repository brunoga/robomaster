package finder

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/windows"
)

func listener(addr string) (*net.UDPConn, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET,
					windows.SO_REUSEADDR, 1)
			})
		},
	}

	conn, err := lc.ListenPacket(context.Background(), "udp4", addr)
	if err != nil {
		return nil, err
	}

	return conn.(*net.UDPConn), nil
}
