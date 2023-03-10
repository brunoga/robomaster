package video

type Resolution byte

const (
	Resolution720p Resolution = iota
	Resolution360p
	Resolution540p
	ResolutionInvalid
)
