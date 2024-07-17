package main

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/brunoga/net/client"
	"github.com/brunoga/net/server"
	"github.com/brunoga/robomaster/support/finder"
)

func listenOnPort40927() *server.Server {
	s, err := server.New("udp", "0.0.0.0:40927", func(conn net.Conn) {
		buffer := make([]byte, 2048)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				break
			}

			fmt.Println("Received data on 40927:", string(buffer[:n-1]))
		}
	})
	if err != nil {
		panic(err)
	}

	err = s.Start()
	if err != nil {
		panic(err)
	}

	return s
}

func connectFromPort10608ToPort10607(robotIP net.IP) *client.Client {
	localAddr := &net.UDPAddr{IP: nil, Port: 10608}
	remoteAddr := &net.UDPAddr{IP: robotIP, Port: 10607}

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		panic(err)
	}

	var c *client.Client

	c, err = client.NewWithConn(conn, client.ScanFullBuffer, func(data []byte) {
		if n := bytes.Index(data, []byte{0x55}); n != -1 {
			m, err := NewMessage(data[n:])
			if err == nil {
				fmt.Printf("** Message: MsgSet:%02x MsgID:%02x\n", m.ProtoCmdSet(), m.ProtoCmdID())
				return
			}
			fmt.Println("Received packet with no message")
		}
	})
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}

	return c
}

func main() {
	// The robot sends data to this port. This is so we can check what it is.
	s := listenOnPort40927()
	defer s.Stop()

	// Connection starts at finding a robot broadcasting on the network. This
	// is already solved so we just reuse the code here.
	f := finder.New(0, nil)

	// Wait 30 seconds to find a robot.
	broadcast, err := f.Find(30 * time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found robot: ", broadcast)

	// This is the main communication channel.
	c := connectFromPort10608ToPort10607(broadcast.SourceIp())
	defer c.Stop()

	// Now we need to ACK the robot. Normally we would send our appID back but
	// we want to connect to any robot so we are just pretend to the robot we
	// are the app it is looking for.
	f.SendACK(broadcast.SourceIp(), broadcast.AppId())

	// Now is when things get interesting. We need to send the first packet
	// back to the robot and from now on we are going to hack things together
	// from scratch. Here is one example such packet I captured:
	//
	// 0000   30 80 8a 2a 00 00 00 10 50 54 64 00 64 00 c0 05
	// 0010   14 00 00 64 00 64 00 64 00 c0 05 14 00 00 64 00
	// 0030   14 00 64 00 c0 05 14 00 00 64 00 01 01 04 01 02
	//
	// The first 2 bytes encode the size in an interesting way. Up to a size of
	// 255, it appears the size is directly represented in the first byte and
	// the second byte is set to 0x80 (i.e. it's most significant bit is set).
	// For bigger sizes, the value is encoded in the first byte and second byte
	// by using all bits of the first byte and some (maybe 3?) lower bits of the
	// second byte. For example, a size of 1472 (0b10111000000) would have the 2
	// bytes set as 0xc0 (0b11000000) and 0x85 (0b10000101). By "appending" the
	// last 3 bits of the second byte to the first byte, we get the original
	// size value (1472). The formula should be something like this:
	//
	// size := 1427
	// b1 := byte(size & 0xff)
	// b2 := byte(0x80 | a>>8)
	//
	// Note this might not be the full picture. The formula above works for the
	// value I saw but there might be other non-obvious operations being done
	// (like making sure that only some bits of the second byte can be used.
	// Maybe 3?).
	//
	// The next 2 bytes are always the same but change each session. They look
	// like they are arbitrary and setting anything here might work.
	//
	// Some packets appear to have a protocol message as defined in Robomaster
	// SDK (https://github.com/dji-sdk/RoboMaster-SDK/blob/ff6646e115ab125af3207a4ed3df42cc76c795b2/src/robomaster/protocol.py#L185).
	// Lets analyze the packet bellow:
	//
	// 0000 24 80 8a 2a c8 54 05 9d b8 54 c8 54 00 00 00 00
	// 0010 0f 01 00 00 55 10 04 56 02 09 2d 27 40 3f 77 01
	// 0020 04 01 5b 0e
	//
	// Byte 21 is 0x55 which is the marker for a protocol message. The following
	// byte is 0x10 (0b10000) and the next one is 0x04 (0b00000100). According
	// to the code pointed above, this encodes a size value of 0x10 (16) bytes,
	// which matches the data size from the 55 value (inclusive) on. The next
	// byte is a crc (we already have the code for that TODO(bga): add it here).
	// The next 2 bytes are sender and receiver. The next one is the lower 8
	// bits of the sequence id, followed by the upper 8 bits. The an attribute
	// byte (is_ack, need_ack, enc). Then if there is a proto associated with
	// the message, the next 2 bytes are the cmdset and cmdid. Then we have
	// the actual proto data (which varies from proto to proto). Finally we
	// have 2 bytes of crc16.
	//
	// Apparently this is going to be considerably easier than it would be
	// otherwise but figuring out the non-msg packets will still be challenging.

	// So, let's try to just send one of the packets we captured again
	// and see what happens.
	err = c.Send([]byte{
		0x30, 0x80,
		0xb6, 0x05,
		0x00, 0x00, 0x00, 0x03, 0xb0, 0x2d,
		0x64, 0x00, 0x64, 0x00, 0xc0, 0x05, 0x14, 0x00, 0x00, 0x64, 0x00, 0x64, 0x00, 0x64, 0x00, 0xc0, 0x05, 0x14, 0x00, 0x00, 0x64, 0x00, 0x14, 0x00, 0x64, 0x00, 0xc0, 0x05, 0x14, 0x00, 0x00, 0x64, 0x00, 0x01, 0x01, 0x04, 0x01, 0x02,
	})
	if err != nil {
		panic(err)
	}

	// This worked and we got a lot of data back but then the connection is
	// closed. There must be some keepalive packets that are required.

	select {}
}
