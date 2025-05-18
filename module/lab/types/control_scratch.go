package types

type ControlScratch int

const (
	StartControlScratch ControlScratch = iota + 2
	StopControlScratch
	ContinueControlScratch
	ExitControlScratch
	DeleteControlScratch
	DeleteAllControlScratch
)
