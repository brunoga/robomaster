package value

type Value[T any] struct {
	Value T `json:"value"`
}
