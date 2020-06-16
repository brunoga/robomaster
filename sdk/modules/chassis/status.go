package chassis

import (
	"fmt"
	"sync"
)

type Status struct {
	m          sync.RWMutex
	static     bool
	upHill     bool
	downHill   bool
	onSlope    bool
	pickUp     bool
	slip       bool
	impactX    bool
	impactY    bool
	impactZ    bool
	rollOver   bool
	hillStatic bool
}

func NewStatusFromData(data string) (*Status, error) {
	s := &Status{
		sync.RWMutex{},
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
	}

	err := s.UpdateFromData(data)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Status) UpdateFromData(data string) error {
	var static, upHill, downHill, onSlope, pickUp, slip, impactX, impactY,
		impactZ, rollOver, hillStatic uint8
	n, err := fmt.Sscanf(data, "%u %u %u", &static, &upHill, &downHill,
		&onSlope, &pickUp, &slip, &impactX, &impactY, &impactZ, &rollOver,
		&hillStatic)
	if err != nil {
		fmt.Errorf("error parsing data: %w", err)
	}
	if n != 11 {
		fmt.Errorf("unexpected number of entries in data: %w", err)
	}

	s.m.Lock()

	s.static = (static == 1)
	s.upHill = (upHill == 1)
	s.downHill = (downHill == 1)
	s.onSlope = (onSlope == 1)
	s.pickUp = (pickUp == 1)
	s.slip = (slip == 1)
	s.impactX = (impactX == 1)
	s.impactY = (impactY == 1)
	s.impactZ = (impactZ == 1)
	s.rollOver = (rollOver == 1)
	s.hillStatic = (hillStatic == 1)

	s.m.Unlock()

	return nil
}

func (s *Status) IsStatic() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.static
}

func (s *Status) IsUphill() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.upHill
}

func (s *Status) IsDownhill() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.downHill
}

func (s *Status) IsOnSlope() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.onSlope
}

func (s *Status) IsPickedUp() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.pickUp
}

func (s *Status) IsSlipping() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.slip
}

func (s *Status) XImpactDetected() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.impactX
}

func (s *Status) YImpactDetected() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.impactY
}

func (s *Status) ZImpactDetected() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.impactZ
}

func (s *Status) IsRolledOver() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.rollOver
}

func (s *Status) IsStaticOnHill() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.hillStatic
}

func (s *Status) String() string {
	str := "Robot status: "

	unknown := true
	if s.IsStatic() {
		unknown = false
		str += "[static]"
	}
	if s.IsUphill() {
		unknown = false
		str += "[uphill]"
	}
	if s.IsDownhill() {
		unknown = false
		str += "[downhill]"
	}
	if s.IsOnSlope() {
		unknown = false
		str += "[on slope]"
	}
	if s.IsPickedUp() {
		unknown = false
		str += "[picked up]"
	}
	if s.IsSlipping() {
		unknown = false
		str += "[slipping]"
	}
	if s.XImpactDetected() {
		unknown = false
		str += "[x axis impact detected]"
	}
	if s.YImpactDetected() {
		unknown = false
		str += "[y axis impact detected]"
	}
	if s.ZImpactDetected() {
		unknown = false
		str += "[z axis impact detected]"
	}
	if s.IsRolledOver() {
		unknown = false
		str += "[rolled over]"
	}
	if s.IsStaticOnHill() {
		unknown = false
		str += "[static on hill]"
	}

	if unknown {
		str += "[unknown]"
	}

	return str
}
