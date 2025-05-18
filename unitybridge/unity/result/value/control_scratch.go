package value

import "github.com/brunoga/robomaster/module/lab/types"

type ControlScratch struct {
	Type types.ControlScratch `json:"type"`
	GUID string               `json:"guid"`
	Sign string               `json:"sign"`
}
