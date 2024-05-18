package dsp

import (
	"testing"
)

func TestSignature(t *testing.T) {
	f, err := New("Anonymous", "Untitled-1")
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	f.dji.Attribute.CreationDate = "2019/09/23"
	f.dji.Attribute.ModifyTime = "09/23/2019 12:59:27"
	f.dji.Attribute.Guid = "92dc6d5736184b9198baa09cf8fd4624"

	f.computeSignature()

	expected := "acec2ac38ea29d0c"
	if f.dji.Attribute.Sign != expected {
		t.Fatalf("expected %q, got %q", expected, f.dji.Attribute.Sign)
	}
}
