package chassis

import (
	"fmt"
	"sync"
)

// Attitude represents chassis attitude information.
type Attitude struct {
	m     sync.RWMutex
	pitch float64
	roll  float64
	yaw   float64
}

// NewAttitude returns a new Attitude instance with the given pitch, roll and
// yaw values (in degrees).
func NewAttitude(pitch, roll, yaw float64) *Attitude {
	return &Attitude{
		sync.RWMutex{},
		pitch,
		roll,
		yaw,
	}
}

// NewAttitudeFromData returns a new Attitude instance based on the given data
// (which usually comes from push events) and a nil error on success and a
// non-nil error on failure.
func NewAttitudeFromData(data string) (*Attitude, error) {
	a := &Attitude{
		sync.RWMutex{},
		0.0,
		0.0,
		0.0,
	}

	err := a.UpdateFromData(data)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Update updates the Attitude instance with the given pitch, roll and yaw angle
// information (in degrees).
func (a *Attitude) Update(pitch, roll, yaw float64) {
	a.m.Lock()

	a.pitch = pitch
	a.roll = roll
	a.yaw = yaw

	a.m.Unlock()
}

// UpdateFromData updates the Attitude based on the given data (which usually
// comes from push events). Returns a nil error on success and a non-nil error
// on failure.
func (a *Attitude) UpdateFromData(data string) error {
	var pitch, roll, yaw float64
	n, err := fmt.Sscanf(data, "%f %f %f", &pitch, &roll, &yaw)
	if err != nil {
		return fmt.Errorf("error parsing data: %w", err)
	}
	if n != 3 {
		return fmt.Errorf(
			"unexpected number of entries in data: %w", err)
	}

	a.m.Lock()

	a.pitch = pitch
	a.roll = roll
	a.yaw = yaw

	a.m.Unlock()

	return nil
}

// Pitch returns the attitude instance pitch information (in degrees).
func (a *Attitude) Pitch() float64 {
	a.m.RLock()
	defer a.m.RUnlock()

	return a.pitch
}

// Roll returns the attitude instance roll information (in degrees).
func (a *Attitude) Roll() float64 {
	a.m.RLock()
	defer a.m.RUnlock()

	return a.roll
}

// Yaw returns the attitude instance yaw information (in degrees).
func (a *Attitude) Yaw() float64 {
	a.m.RLock()
	defer a.m.RUnlock()

	return a.yaw
}
