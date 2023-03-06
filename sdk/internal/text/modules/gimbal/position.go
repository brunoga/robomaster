package gimbal

import (
	"fmt"
	"sync"
)

type Position struct {
	m     sync.RWMutex
	pitch float64
	yaw   float64
}

func NewPosition(pitch, yaw float64) *Position {
	return &Position{
		sync.RWMutex{},
		pitch,
		yaw,
	}
}

func NewPositionFromData(data string) (*Position, error) {
	p := &Position{
		sync.RWMutex{},
		0.0,
		0.0,
	}

	err := p.UpdateFromData(data)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Position) Update(pitch, yaw float64) {
	p.m.Lock()

	p.pitch = pitch
	p.yaw = yaw

	p.m.Unlock()
}

func (p *Position) UpdateFromData(data string) error {
	var pitch, yaw float64
	n, err := fmt.Sscanf(data, "%f %f", &pitch, &yaw)
	if err != nil {
		return fmt.Errorf("error parsing data: %w", err)
	}
	if n != 2 {
		return fmt.Errorf(
			"unexpected number of entries in data: %w", err)
	}

	p.Update(pitch, yaw)

	return nil
}

func (p *Position) Pitch() float64 {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.pitch
}

func (p *Position) Yaw() float64 {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.yaw
}
