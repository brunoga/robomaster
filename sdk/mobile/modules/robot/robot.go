package robot

import "github.com/brunoga/robomaster/sdk/modules/robot"

type Robot struct {
	*robot.Robot
}

func New(r *robot.Robot) *Robot {
	return &Robot{r}
}
