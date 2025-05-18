package value

import "github.com/brunoga/robomaster/module/lab/types"

type ScratchState struct {
	Key                    string                          `json:"key"`
	State                  types.ScratchState              `json:"state"`
	BlockID                string                          `json:"blockId"`
	BlockExecuteState      int                             `json:"blockExecuteState"`
	BlockType              int                             `json:"blockType"`
	BlockExecutePercentage int                             `json:"blockExecutePercentage"`
	Variables              []types.ScratchVariableInfo     `json:"variables"`
	Lists                  []types.ScratchListVariableInfo `json:"lists"`
}
