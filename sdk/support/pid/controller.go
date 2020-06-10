package pid

// Controller is the interface for generic controllers that generate an output
// based on some input parameter.
type Controller interface {
	// Output returns an output associated with the given input.
	Output(float64) float64
}
