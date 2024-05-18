package value

type GamePadActivationSettings struct {
	IsActivated  bool   `json:"isActivated"`
	ActivateTime int64  `json:"activateTime"`
	SerialNumber string `json:"serialNumber"`
}
