package value

type GimbalAttitude struct {
	Pitch       float32 `json:"pitch"`
	Yaw         float32 `json:"yaw"`
	Roll        float32 `json:"roll"`
	YawOpposite float32 `json:"yawOpposite"`
	PitchSpeed  float32 `json:"pitchSpeed"`
	YawSpeed    float32 `json:"yawSpeed"`
	RollSpeed   float32 `json:"rollSpeed"`
}
