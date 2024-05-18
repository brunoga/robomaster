package finder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/brunoga/robomaster/unitybridge/support"
)

var (
	// Expected broadcast message length.
	broadcastLen = 24

	// Expected decoded broadcast message header.
	broadcastHeader = []byte{90, 91}
)

// Broadcast represents a broadcast message sent by a Robomaster robot.
type Broadcast struct {
	isPairing bool
	sourceIp  net.IP
	sourceMac net.HardwareAddr
	appId     uint64
}

// ParseBroadcast parses the given data as a BroadcastMessage. It
// returns the associated BroadcastMessage instance pointer and a nil error on
// success and a nil BroadcastMessage and a non-nil error on failure.
func ParseBroadcast(data []byte) (*Broadcast, error) {
	if len(data) != broadcastLen {
		return nil, fmt.Errorf("unexpected broadcast message length")
	}

	// Decode incoming data.
	support.SimpleEncryptDecrypt(data)

	if !bytes.HasPrefix(data, broadcastHeader) {
		return nil, fmt.Errorf("invalid broadcast message header")
	}

	// First byte tells us if this is a pairing message.
	isPairing := (data[2] & 1) > 0

	// Then we get the rest of the data trivially.
	sourceIp := data[6:10]
	sourceMac := data[10:16]
	appId := binary.LittleEndian.Uint64(data[16:])

	return &Broadcast{
		isPairing,
		sourceIp,
		sourceMac,
		appId,
	}, nil
}

func (b *Broadcast) IsPairing() bool {
	return b.isPairing
}

func (b *Broadcast) SourceIp() net.IP {
	return b.sourceIp
}

func (b *Broadcast) SourceMac() net.HardwareAddr {
	return b.sourceMac
}

func (b *Broadcast) AppId() uint64 {
	return b.appId
}

func (b *Broadcast) String() string {
	return fmt.Sprintf("IsPairing:%t, SourceIp:%s, SourceMac:%s, AppId:%d",
		b.isPairing, b.sourceIp, b.sourceMac, b.appId)
}
