package main

import "testing"

var testMessage = []byte{
	0x55, 0x10, 0x04, 0x56, 0x02, 0x09, 0x2d, 0x27,
	0x40, 0x3f, 0x77, 0x01, 0x04, 0x01, 0x5b, 0x0e,
}

func TestMessage(t *testing.T) {
	m, err := NewMessage(testMessage)
	if err != nil {
		t.Fatalf("NewMessage() failed: %v", err)
	}

	if m.Sender() != 0x02 {
		t.Fatalf("Sender() failed: %02x != %02x", m.Sender(), 0x56)
	}

	if m.Receiver() != 0x09 {
		t.Fatalf("Receiver() failed: %02x != %02x", m.Receiver(), 0x02)
	}

	if m.Sequence() != 10029 {
		t.Fatalf("Sequence() failed: %d != %d", m.Sequence(), 10029)
	}

	if m.Attrs() != 0x40 {
		t.Fatalf("Attrs() failed: %02x != %02x", m.Attrs(), 0x27)
	}

	if !m.AttrIsAck() {
		t.Fatalf("AttrIsAck() failed: true != false")
	}

	if m.AttrNeedsAck() {
		t.Fatalf("AttrNeedsAck() failed: false != true")
	}

	if m.ProtoData() == nil {
		t.Fatalf("ProtoData() failed: nil")
	}

	if len(m.ProtoData()) != 3 {
		t.Fatalf("ProtoData() failed: %d != %d", len(m.ProtoData()), 3)
	}

	if m.ProtoCmdSet() != 0x3f {
		t.Fatalf("ProtoCmdSet() failed: %02x != %02x", m.ProtoCmdSet(), 0x3f)
	}

	if m.ProtoCmdID() != 0x77 {
		t.Fatalf("ProtoCmdID() failed: %02x != %02x", m.ProtoCmdID(), 0x77)
	}
}
