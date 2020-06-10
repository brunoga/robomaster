package pid

// PIDController is a Controller with proportional, integral and derivative
// components.
type PIDController struct {
	pController Controller
	iController Controller
	dController Controller
	minOutput   float64
	maxOutput   float64

	disableIntegrator        bool
	previousIntegratorOutput float64
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
		false,
		0.0,
	}
}

// Output returns the sum of the individual components (proportional, integral,
// derivative) associated with the given input after clamping the result to the
// minimum/maximum allowed values.
func (p *PIDController) Output(currentError float64) float64 {
	pOutput := p.pController.Output(currentError)
	dOutput := p.dController.Output(currentError)

	// Disable integrator in case integral windup is in effect.
	var iOutput float64
	if p.disableIntegrator {
		iOutput = p.iController.Output(0.0)
	} else {
		iOutput = p.iController.Output(currentError)
	}

	output := pOutput + iOutput + dOutput

	/// Clamping.
	limitedOutput := output
	if output < p.minOutput {
		limitedOutput = p.minOutput
	} else if output > p.maxOutput {
		limitedOutput = p.maxOutput
	}

	// Integral windup.
	if output != limitedOutput {
		if (iOutput < 0 && currentError < 0) ||
			(iOutput > 0 && currentError > 0) {
			p.disableIntegrator = true
		} else {
			p.disableIntegrator = false
		}
	} else {
		p.disableIntegrator = false
	}

	return limitedOutput
}
