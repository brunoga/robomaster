package controller

const (
	minBound = float64(-660)
	maxBound = float64(660)
	offset   = float64(1024)
)

type StickPosition struct {
	X float64
	Y float64
}

func (s *StickPosition) InterpolatedX() uint64 {
	if s != nil {
		return uint64(lerp(minBound, maxBound, float64(s.X+1)*0.5) +
			offset)
	}

	return 0
}

func (s *StickPosition) InterpolatedY() uint64 {
	if s != nil {
		return uint64(lerp(minBound, maxBound, float64(s.Y+1)*0.5) +
			offset)
	}

	return 0
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}
