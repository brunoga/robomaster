package value

type ChassisPosition struct {
	TaskID      uint8   `json:"taskId"`
	IsCancel    uint8   `json:"isCancel"`
	ControlMode uint8   `json:"controlMode"`
	X           float32 `json:"positionX"`
	Y           float32 `json:"positionY"`
	Z           float32 `json:"positionYaw"`
}
