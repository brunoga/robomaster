package camera

type VideoQuality uint8

const (
	VideoQualityNormal VideoQuality = iota
	VideoQualityGood
	VideoQualityBest
	VideoQualityCount
)

// String returns the video quality as a string.
func (vq VideoQuality) String() string {
	switch vq {
	case VideoQualityNormal:
		return "Normal"
	case VideoQualityGood:
		return "Good"
	case VideoQualityBest:
		return "Best"
	default:
		return "Invalid"
	}
}

// ToRate returns the video quality as a rate (in megabits per second).
func (vq VideoQuality) ToRate() float32 {
	switch vq {
	case VideoQualityNormal:
		return 2.4
	case VideoQualityGood:
		return 3.4
	case VideoQualityBest:
		return 6.0
	default:
		return 0.0
	}
}

// VideoQualityFromRate returns the video quality from a rate (in megabits per second).
func VideoQualityFromRate(rate float32) VideoQuality {
	switch rate {
	case 2.4:
		return VideoQualityNormal
	case 3.4:
		return VideoQualityGood
	case 6.0:
		return VideoQualityBest
	default:
		return VideoQualityCount // invalid
	}
}

// Valid returns true if the video quality is valid.
func (vq VideoQuality) Valid() bool {
	return vq < VideoQualityCount
}
