package chassis

import (
	"fmt"
	"sync"
)

type Speed struct {
	m sync.RWMutex
	x float64
	y float64
	z float64
}

func NewSpeed(x, y, z float64) *Speed {
	return &Speed{
		sync.RWMutex{},
		x,
		y,
		z,
	}
}

func NewSpeedFromData(data string) (*Speed, error) {
	s := &Speed{
		sync.RWMutex{},
		0.0,
		0.0,
		0.0,
	}

	err := s.UpdateFromData(data)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Speed) Update(x, y, z float64) {
	s.m.Lock()

	s.x = x
	s.y = y
	s.z = z

	s.m.Unlock()
}

func (s *Speed) UpdateFromData(data string) error {
	var x, y, z float64

	n, err := fmt.Sscanf(data, "%f %f %f %f %f %f %f", &x, &y, &z)
	if err != nil {
		fmt.Errorf("error parsing data: %w", err)
	}
	if n != 7 {
		fmt.Errorf(
			"unexpected number of entries in data: %w", err)
	}

	s.Update(x, y, z)

	return nil
}

func (s *Speed) X() float64 {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.x
}

func (s *Speed) Y() float64 {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.y
}

func (s *Speed) Z() float64 {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.z
}
