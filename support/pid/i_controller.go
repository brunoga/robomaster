package pid

import (
	"time"
)

// IController is a Controller with only an integral component.
type IController struct {
	ki float64

	integral   float64
	lastOutput time.Time
}

// NewIController returns a new instance of an IController that uses the given
// ki multiplier.
func NewIController(ki float64) Controller {
	return &IController{
		ki,
		0.0,
		time.Time{},
	}
}

// Output returns the integral output associated with the given
// currentError.
func (i *IController) Output(currentError float64) float64 {
	now := time.Now()

	var deltaTime time.Duration
	if !i.lastOutput.IsZero() {
		deltaTime = now.Sub(i.lastOutput)
	}

	integral := (i.integral + currentError) * deltaTime.Seconds()

	i.lastOutput = now
	i.integral = integral

	return i.ki * integral
}

// Adjust modifies the integral term to counteract windup.
func (i *IController) Adjust(adjustment float64) {
	// The adjustment is subtracted because if the output is too high, we want
	// to reduce the integral and vice versa.
	i.integral -= adjustment / i.ki

	// Ensure integral does not become negative
	if i.integral < 0 {
		i.integral = 0
	}
}
