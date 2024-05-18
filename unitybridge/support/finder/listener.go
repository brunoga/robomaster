//go:build !windows

package finder

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func listener(addr string) (*net.UDPConn, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				unix.SetsockoptInt(int(fd), unix.SOL_SOCKET,
					unix.SO_REUSEADDR, 1)
			})
		},
	}

	conn, err := lc.ListenPacket(context.Background(), "udp4", addr)
	if err != nil {
		return nil, err
	}

	return conn.(*net.UDPConn), nil
}
