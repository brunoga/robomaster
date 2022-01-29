package pid

import (
	"time"
)

type IController struct {
	ki float64

	sum        float64
	lastOutput time.Time
}

func NewIController(ki float64) Controller {
	return &IController{
		ki,
		0.0,
		time.Time{},
	}
}

func (i *IController) Output(currentError float64) float64 {
	now := time.Now()

	var deltaTime time.Duration
	if !i.lastOutput.IsZero() {
		deltaTime = now.Sub(i.lastOutput)
	}

	deltaSeconds := deltaTime.Seconds()

	i.sum += currentError * float64(deltaSeconds)

	i.lastOutput = now

	return i.ki * i.sum
}
