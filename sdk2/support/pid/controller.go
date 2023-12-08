package pid

// Controller is the interface for generic controllers that generate an output
// based on some input parameter.
type Controller interface {
	// Output returns an output associated with the given current error.
	// Generally speaking, the oujtput will be used to minimize this error.
	Output(currentError float64) float64

	// Adjust allows adjusting the internal state of the controller to
	// counteract known effects (integral windup, etc.).
	Adjust(adjustment float64)
}
