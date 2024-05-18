package camera

type VideoFormat uint8

const (
	VideoFormat720p_30 VideoFormat = iota
	VideoFormat1080p_30
	VideoFormat720p_60  // this might not be actually supported
	VideoFormat1080p_60 // this might not be actually supported
	VideoFormatCount
)

// String returns the video format as a string.
func (vf VideoFormat) String() string {
	switch vf {
	case VideoFormat720p_30:
		return "1280x720@30"
	case VideoFormat1080p_30:
		return "1920x1080@30"
	case VideoFormat720p_60:
		return "1280x720@60"
	case VideoFormat1080p_60:
		return "1920x1080@60"
	default:
		return "Invalid"
	}
}

// Valid returns true if the video format is valid.
func (vf VideoFormat) Valid() bool {
	return vf < VideoFormatCount
}
