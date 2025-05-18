package value

import "github.com/brunoga/robomaster/module/lab/types"

type ScratchOutputInfo struct {
	InfoType    types.ScratchInfo      `json:"infoType"`
	InfoLevel   types.ScratchInfoLevel `json:"infoLevel"`
	InfoContent string                 `json:"infoContent"`
}
