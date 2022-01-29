package pid

import (
	"time"
)

type DController struct {
	kd float64

	lastError  float64
	lastOutput time.Time
}

func NewDController(kd float64) Controller {
	return &DController{
		kd,
		0.0,
		time.Time{},
	}
}

func (d *DController) Output(currentError float64) float64 {
	now := time.Now()

	var deltaTime time.Duration
	if !d.lastOutput.IsZero() {
		deltaTime = now.Sub(d.lastOutput)
	}

	deltaSeconds := deltaTime.Seconds()

	deltaError := currentError - d.lastError

	d.lastOutput = now
	d.lastError = currentError

	return d.kd * deltaError / float64(deltaSeconds)
}
