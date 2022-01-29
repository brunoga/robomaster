package pid

type PIDController struct {
	pController Controller
	iController Controller
	dController Controller
	minOutput   float64
	maxOutput   float64

	disableIntegrator        bool
	previousIntegratorOutput float64
}

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

func (p *PIDController) Output(currentError float64) float64 {
	pOutput := p.pController.Output(currentError)
	dOutput := p.dController.Output(currentError)

	var iOutput float64
	if p.disableIntegrator {
		iOutput = p.iController.Output(0.0)
	} else {
		iOutput = p.iController.Output(currentError)
	}

	output := pOutput + iOutput + dOutput

	limitedOutput := output
	if output < p.minOutput {
		limitedOutput = p.minOutput
	} else if output > p.maxOutput {
		limitedOutput = p.maxOutput
	}

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
