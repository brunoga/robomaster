package blaster

import (
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules"
	"strconv"
)

type Blaster struct {
	control *modules.Control
}

func New(control *modules.Control) *Blaster {
	return &Blaster{
		control,
	}
}

func (b *Blaster) SetNumBeads(beads int) error {
	return b.control.SendDataExpectOk(fmt.Sprintf(
		"blaster bead %d;", beads))
}

func (b *Blaster) GetNumBeads() (int, error) {
	data, err := b.control.SendAndReceiveData("blaster bead ?;")
	if err != nil {
		return -1, fmt.Errorf("error sending sdk command: %w", err)
	}

	numBeads, err := strconv.Atoi(data)
	if err != nil {
		return -1, fmt.Errorf("error parsing data: %w", err)
	}

	return numBeads, nil
}

func (b *Blaster) Fire() error {
	return b.control.SendDataExpectOk("blaster fire;")
}
