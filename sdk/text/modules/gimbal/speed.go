package gimbal

import (
	"fmt"
	"sync"
)

type Speed struct {
	m     sync.RWMutex
	pitch float64
	yaw   float64
}

func NewSpeed(pitch, yaw float64) *Speed {
	return &Speed{
		sync.RWMutex{},
		pitch,
		yaw,
	}
}

func NewSpeedFromData(data string) (*Speed, error) {
	s := &Speed{
		sync.RWMutex{},
		0.0,
		0.0,
	}

	err := s.UpdateFromData(data)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Speed) Update(pitch, yaw float64) {
	s.m.Lock()

	s.pitch = pitch
	s.yaw = yaw

	s.m.Unlock()
}

func (s *Speed) UpdateFromData(data string) error {
	var pitch, yaw float64

	n, err := fmt.Sscanf(data, "%f %f", &pitch, &yaw)
	if err != nil {
		fmt.Errorf("error parsing data: %w", err)
	}
	if n != 2 {
		fmt.Errorf(
			"unexpected number of entries in data: %w", err)
	}

	s.Update(pitch, yaw)

	return nil
}

func (s *Speed) Pitch() float64 {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.pitch
}

func (s *Speed) Yaw() float64 {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.yaw
}
