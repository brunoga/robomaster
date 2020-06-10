package pid

// PController is a Controller with only a proportional component.
type PController struct {
	kp float64
}

// NewPController returns a new instance of a PController that uses the given kp
// multiplier.
func NewPController(kp float64) Controller {
	return &PController{
		kp,
	}
}

// Output returns the proportional output associated with the given
// currentError.
func (p *PController) Output(currentError float64) float64 {
	return p.kp * currentError
}
