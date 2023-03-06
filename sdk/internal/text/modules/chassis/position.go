package chassis

import (
	"fmt"
	"sync"
)

type Position struct {
	m sync.RWMutex
	x float64
	y float64
	z float64
}

func NewPosition(x, y, z float64) *Position {
	return &Position{
		sync.RWMutex{},
		x,
		y,
		z,
	}
}

func NewPositionFromData(data string) (*Position, error) {
	p := &Position{
		sync.RWMutex{},
		0.0,
		0.0,
		0.0,
	}

	err := p.UpdateFromData(data)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Position) Update(x, y, z float64) {
	p.m.Lock()

	p.x = x
	p.y = y
	p.z = z

	p.m.Unlock()
}

func (p *Position) UpdateFromData(data string) error {
	var x, y, z float64
	n, err := fmt.Sscanf(data, "%f %f %f", &x, &y, &z)
	if err != nil {
		fmt.Errorf("error parsing data: %w", err)
	}
	if n != 3 {
		fmt.Errorf(
			"unexpected number of entries in data: %w", err)
	}

	p.Update(x, y, z)

	return nil
}

func (p *Position) X() float64 {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.x
}

func (p *Position) Y() float64 {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.y
}

func (p *Position) Z() float64 {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.z
}
