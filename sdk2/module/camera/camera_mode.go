package camera

type Mode uint8

const (
	ModePhoto Mode = iota
	ModeVideo
	ModeCount
)

// String returns the camera mode as a string.
func (cm Mode) String() string {
	switch cm {
	case ModePhoto:
		return "Photo"
	case ModeVideo:
		return "Video"
	default:
		return "Invalid"
	}
}

// Valid returns true if the camera mode is valid.
func (cm Mode) Valid() bool {
	return cm < ModeCount
}
