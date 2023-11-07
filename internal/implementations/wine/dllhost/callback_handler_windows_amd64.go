package main

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/brunoga/unitybridge/unity/event"
)

type CallbackHandler struct {
	eventFile *os.File
}

func (c *CallbackHandler) HandleCallback(e *event.Event, data []byte,
	tag uint64) {
	var b bytes.Buffer

	// TOOD(bga): Add error checking.
	binary.Write(&b, binary.BigEndian, e.Code())
	binary.Write(&b, binary.BigEndian, tag)
	binary.Write(&b, binary.BigEndian, uint16(len(data)))
	b.Write(data)

	b.WriteTo(c.eventFile)
}
