package sdk

import (
	"fmt"

	"github.com/brunoga/robomaster/sdk/finder"
	"github.com/brunoga/robomaster/sdk/internal/binary"
	"github.com/brunoga/robomaster/sdk/internal/text"
)

// Mode represents an SDK mode. Currently text and binary protocols are
// supported.
type Mode uint8

const (
	Text   Mode = iota // Text mode protocol
	Binary             // Binary mode protocol
)

// Sdk is the interface for SDK implementations.
type Sdk interface {
	finder.Finder
}

// New returns a new SDK instance with the selected mode and a nil error on
// success and nil and a non-nil error on failure
func New(m Mode) (Sdk, error) {
	switch m {
	case Text:
		return text.NewFinder(), nil
	case Binary:
		return binary.NewFinder(), nil
	}

	return nil, fmt.Errorf("invalid sdk mode")
}
