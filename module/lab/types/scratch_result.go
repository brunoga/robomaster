package types

type ScratchResult int

const (
	SuccessScratchResult ScratchResult = iota
	FailScratchResult
	TimeoutScratchResult
	InstallSkillSuccessScratchResult
	InstallSkillFailScratchResult
	InstallSkillTimeoutScratchResult
	UninstallSkillSuccessScratchResult
	UninstallSkillFailScratchResult
	UninstallSkillTimeoutScratchResult
)
