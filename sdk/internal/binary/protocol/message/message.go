package message

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/command"
	"github.com/brunoga/robomaster/sdk/internal/binary/support"
)

const (
	msgEnvelopeSize = 13 // 11 bytes of header + 2 bytes of CRC16.
)

// Message is a message sent to or received from the robot.
type Message struct {
	data []byte
}

// New creates a new Message with the given sender, receiver and command.
func New(sender, receiver byte, cmd command.Request) *Message {
	protoData := cmd.Data()

	size := msgEnvelopeSize + len(protoData)

	data := make([]byte, size)

	data[0] = 0x55
	data[1] = byte(size & 0xff)
	data[2] = byte((size>>8)&0x3 | 0x4)
	data[3] = support.ComputeCRC8(data[0:3])
	data[4] = sender
	data[5] = receiver

	binary.LittleEndian.PutUint16(data[6:8], getSequenceID())

	data[8] = 2 << 5 // Needs ack.
	data[9] = cmd.Set()
	data[10] = cmd.ID()

	copy(data[11:], protoData)

	binary.LittleEndian.PutUint16(data[size-2:], support.ComputeCRC16(data[:size-2]))

	return &Message{
		data: data,
	}
}

// NewFromData creates a new Message from the given data. If the data is not
// enough to create a new Message, it returns io.EOF as the error and the
// remaining data.
func NewFromData(data []byte) (*Message, []byte, error) {
	if len(data) < msgEnvelopeSize {
		// Buffer is too small, try to get more data.
		return nil, data, io.EOF
	}

	if data[0] != 0x55 {
		// Magic number is invalid.
		//
		// TODO(bga): Maybe we should discard data?
		return nil, data, fmt.Errorf("NewMessageFromData: invalid magic number 0x%x", data[0])
	}

	if support.ComputeCRC8(data[0:3]) != data[3] {
		// CRC8 is invalid.
		//
		// TODO(bga): Maybe we should discard data?
		return nil, data, fmt.Errorf("NewMessageFromData: invalid CRC8")
	}

	msgLen := uint16(data[2]&0x3)<<8 | uint16(data[1])

	if len(data) < int(msgLen) {
		// Buffer is too small, try to get more data.
		return nil, data, io.EOF
	}

	if support.ComputeCRC16(data[:msgLen-2]) != binary.LittleEndian.Uint16(data[msgLen-2:]) {
		return nil, data, fmt.Errorf("NewMessageFromData: invalid CRC16")
	}

	m := &Message{
		data: data[:msgLen],
	}

	return m, data[msgLen:], nil
}

// Sender returns the sender of the message.
func (m *Message) Sender() byte {
	return m.data[4]
}

// Receiver returns the receiver of the message.
func (m *Message) Receiver() byte {
	return m.data[5]
}

// SequenceID returns the sequence ID of the message.
func (m *Message) SequenceID() uint16 {
	return binary.LittleEndian.Uint16(m.data[6:8])
}

// CmdSet returns the command set of the message.
func (m *Message) CmdSet() byte {
	return m.data[9]
}

// CmdID returns the command ID of the message.
func (m *Message) CmdID() byte {
	return m.data[10]
}

// IsResponse returns true if the message is a response.
func (m *Message) IsResponse() bool {
	return (m.data[8] >> 7) != 0
}

// Data returns the raw data of the message.
func (m *Message) Data() []byte {
	return m.data
}

// Command returns the Command associated with the message.
func (m *Message) Command() command.Command {
	isResponse := m.IsResponse()

	var cmd command.Command
	if isResponse {
		cmd = command.GetResponse(m.data[9], m.data[10], m.data[11:len(m.data)-2])
	} else {
		cmd = command.GetRequest(m.data[9], m.data[10], m.data[11:len(m.data)-2])
	}

	return cmd
}
