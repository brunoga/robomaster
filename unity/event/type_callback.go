package event

// TypeCallback is the type of the callback function that will be called
// when a type event is received. It will not be called for events that
// have sub-types.
type TypeCallback func(data []byte, dataType DataType)
