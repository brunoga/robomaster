package gimbal

type Attitude Position

func NewAttitude(pitch, yaw float64) *Attitude {
	position := NewPosition(pitch, yaw)

	return (*Attitude)(position)
}
