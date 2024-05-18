package pid

// PIDController is a Controller with proportional, integral and derivative
// components.
type PIDController struct {
	pController Controller
	iController Controller
	dController Controller
	minOutput   float64
	maxOutput   float64
}

// NewPIDController returns a new PIDController instance that use the given kp,
// ki and kd values as multipliers for the associated individual components
// (proportional, integral and derivative). minOutput and maxOutput are used to
// clamp the output to known minimum and maximum values.
func NewPIDController(kp, ki, kd, minOutput, maxOutput float64) Controller {
	return &PIDController{
		NewPController(kp),
		NewIController(ki),
		NewDController(kd),
		minOutput,
		maxOutput,
	}
}

// Output returns the sum of the individual components (proportional, integral,
// derivative) associated with the given input after clamping the result to the
// minimum/maximum allowed values.
func (p *PIDController) Output(currentError float64) float64 {
	pOutput := p.pController.Output(currentError)
	dOutput := p.dController.Output(currentError)
	iOutput := p.iController.Output(currentError)

	output := pOutput + iOutput + dOutput

	/// Clamping.
	limitedOutput := output
	if output < p.minOutput {
		limitedOutput = p.minOutput
	} else if output > p.maxOutput {
		limitedOutput = p.maxOutput
	}

	// Integral windup handling.
	if output != limitedOutput {
		p.iController.Adjust(output - limitedOutput)
	}

	return limitedOutput
}

// Adjust does nothing for now.
func (p *PIDController) Adjust(adjustment float64) {}
