package value

import "github.com/brunoga/robomaster/module/lab/types"

type ScratchCallback struct {
	Key        string              `json:"key"`
	ResultCode types.ScratchResult `json:"resultCode"`
}
