package types

type ScratchInfoLevel int

const (
	NoneScratchInfoLevel ScratchInfoLevel = iota
	InfoScratcInfoLevel
	DebugScratchInfoLevel
	ErrorScratchInfoLevel
	FatalScratchInfoLevel
)
