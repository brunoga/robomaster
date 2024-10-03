package value

type GamePadActivationSettings struct {
	IsActivated  bool   `json:"IsActivated"`
	ActivateTime int64  `json:"ActivateTime"`
	SerialNumber string `json:"SerialNumber"`
}
