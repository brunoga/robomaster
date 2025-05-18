package value

type ScratchErrorInfo struct {
	Key       string `json:"key"`
	ErrorLine int    `json:"errorLine"`
	ErrorMsg  string `json:"errorMsg"`
}
