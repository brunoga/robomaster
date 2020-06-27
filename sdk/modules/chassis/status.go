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
	n, err := fmt.Sscanf(data, "%d %d %d %d %d %d %d %d %d %d %d", &static,
		&upHill, &downHill, &onSlope, &pickUp, &slip, &impactX, &impactY, &impactZ,
		&rollOver, &hillStatic)
	if err != nil {
		fmt.Errorf("error parsing data: %w", err)
	}
	if n != 11 {
		fmt.Errorf("unexpected number of entries in data: %w", err)
	}

	s.m.Lock()

	s.static = static == 1
	s.upHill = upHill == 1
	s.downHill = downHill == 1
	s.onSlope = onSlope == 1
	s.pickUp = pickUp == 1
	s.slip = slip == 1
	s.impactX = impactX == 1
	s.impactY = impactY == 1
	s.impactZ = impactZ == 1
	s.rollOver = rollOver == 1
	s.hillStatic = hillStatic == 1

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

func (s *Status) IsUnknown() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return !s.static && !s.upHill && !s.downHill && !s.onSlope && !s.pickUp &&
		!s.slip && !s.impactX && !s.impactY && !s.impactZ && !s.rollOver &&
		!s.hillStatic
}

func (s *Status) Equals(s2 *Status) bool {
	s.m.RLock()
	s2.m.RLock()
	defer s.m.RUnlock()
	defer s2.m.RUnlock()

	return s.static == s2.static && s.upHill == s2.upHill &&
		s.downHill == s2.downHill && s.onSlope == s2.onSlope &&
		s.pickUp == s2.pickUp && s.slip == s2.slip &&
		s.impactX == s2.impactX && s.impactY == s2.impactY &&
		s.impactZ == s2.impactZ && s.rollOver == s2.rollOver &&
		s.hillStatic == s2.hillStatic
}

func (s *Status) String() string {
	s.m.RLock()
	defer s.m.RUnlock()

	str := "Robot status: "

	unknown := true
	if s.static {
		unknown = false
		str += "[static]"
	}
	if s.upHill {
		unknown = false
		str += "[uphill]"
	}
	if s.downHill {
		unknown = false
		str += "[downhill]"
	}
	if s.onSlope {
		unknown = false
		str += "[on slope]"
	}
	if s.pickUp {
		unknown = false
		str += "[picked up]"
	}
	if s.slip {
		unknown = false
		str += "[slipping]"
	}
	if s.impactX {
		unknown = false
		str += "[x axis impact detected]"
	}
	if s.impactY {
		unknown = false
		str += "[y axis impact detected]"
	}
	if s.impactZ {
		unknown = false
		str += "[z axis impact detected]"
	}
	if s.rollOver {
		unknown = false
		str += "[rolled over]"
	}
	if s.hillStatic {
		unknown = false
		str += "[static on hill]"
	}

	if unknown {
		str += "[unknown]"
	}

	return str
}
