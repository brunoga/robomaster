package value

import "github.com/brunoga/robomaster/unitybridge/unity/task"

type TaskStatus struct {
	TaskType task.Type   `json:"taskId"` // JSON name seems mismatched. Check.
	Percent  float64     `json:"percent"`
	Status   task.Status `json:"status"`
}
