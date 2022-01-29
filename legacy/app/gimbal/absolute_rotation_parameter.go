package gimbal

type absoluteRotationParameter struct {
	Pitch int16 `json:"pitch"`
	Yaw   int16 `json:"yaw"`
	Time  int16 `json:"time"`
}
