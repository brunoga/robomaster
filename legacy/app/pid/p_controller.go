package pid

type PController struct {
	kp float64
}

func NewPController(kp float64) Controller {
	return &PController{
		kp,
	}
}

func (p *PController) Output(currentError float64) float64 {
	return p.kp * currentError
}
