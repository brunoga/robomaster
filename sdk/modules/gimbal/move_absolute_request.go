package gimbal

type MoveAbsoluteRequest struct {
	PitchAngleDegrees          float64
	YawAngleDegrees            float64
	PitchSpeedDegreesPerSecond float64
	YawSpeedDegreesPerSecond   float64
}
