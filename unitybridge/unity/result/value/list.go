package value

// List is a result value that holds a list of values.
type List[T any] struct {
	List []T `json:"list"`
}
