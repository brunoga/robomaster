package pid

import (
	"time"
)

// DController is a Controller with only a derivative component.
type DController struct {
	kd float64

	lastError  float64
	lastOutput time.Time
}

// NewDController returns a new instance of a DController that uses the given kd
// multiplier.
func NewDController(kd float64) Controller {
	return &DController{
		kd,
		0.0,
		time.Time{},
	}
}

// Output returns the derivative output associated with the given
// currentError.
func (d *DController) Output(currentError float64) float64 {
	now := time.Now()

	var deltaTime time.Duration
	if !d.lastOutput.IsZero() {
		deltaTime = now.Sub(d.lastOutput)
	} else {
		d.lastOutput = now
		d.lastError = currentError
		return 0.0
	}

	derivative := (currentError - d.lastError) / deltaTime.Seconds()

	d.lastOutput = now
	d.lastError = currentError

	return d.kd * derivative
}

// Adjust does nothing for now.
func (d *DController) Adjust(adjustment float64) {}
