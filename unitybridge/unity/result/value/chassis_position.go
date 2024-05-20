package value

import "github.com/brunoga/robomaster/unitybridge/unity/task"

type ChassisPosition struct {
	TaskType    task.Type `json:"taskId"`
	IsCancel    uint8     `json:"isCancel"`
	ControlMode uint8     `json:"controlMode"`
	X           float32   `json:"positionX"`
	Y           float32   `json:"positionY"`
	Z           float32   `json:"positionYaw"`
}
