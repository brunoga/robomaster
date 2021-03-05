package protocol

const (
	setSdkModeReqSize = 1
)

type SetSdkMode struct {
	*Data

	enabled bool
}

func NewSetSdkMode() *SetSdkMode {
	return &SetSdkMode{
		NewData("SetSdkMode", 0x3f, 0x01),
		true,
	}
}

func (s *SetSdkMode) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func (s *SetSdkMode) PackReq() []byte {
	buf := make([]byte, setSdkModeReqSize)
	if s.enabled {
		buf[0] = 1
	}

	return buf
}
