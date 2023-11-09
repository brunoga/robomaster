// go:build windows & amd64

package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

type CallbackHandler struct {
	eventFile *os.File
}

func (c *CallbackHandler) HandleCallback(eventCode uint64, data []byte,
	tag uint64) {
	var b bytes.Buffer

	// TOOD(bga): Add error checking.
	binary.Write(&b, binary.BigEndian, eventCode)
	binary.Write(&b, binary.BigEndian, tag)
	binary.Write(&b, binary.BigEndian, uint16(len(data)))
	b.Write(data)

	b.WriteTo(c.eventFile)
}
