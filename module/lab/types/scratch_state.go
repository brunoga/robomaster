package types

type ScratchState int

const (
	ScratchStateIdle ScratchState = iota
	ScratchStateDownloading
	ScratchStateExecuting
	ScratchStateCustomSkillRunning
	ScratchStateOfflineSkillRunning
)
