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
	var x float64 = 0.5
	if s != nil {
		x = s.X
	}

	return uint64(lerp(minBound, maxBound, (x+1)*0.5) +
		offset)
}

func (s *StickPosition) InterpolatedY() uint64 {
	var y float64 = 0.5
	if s != nil {
		y = s.Y
	}

	y = -y

	return uint64(lerp(minBound, maxBound, float64(y+1)*0.5) +
		offset)
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}
