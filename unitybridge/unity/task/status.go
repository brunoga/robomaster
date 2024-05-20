package task

type Status int8

const (
	StatusRunning Status = iota
	StatusSuccess
	StatusFailure
	StatusCount
)
