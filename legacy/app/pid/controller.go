package pid

type Controller interface {
	Output(currentError float64) float64
}
