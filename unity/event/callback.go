package event

// Callback is the prototype for functions that will be called when a reply to
// an event is received.
//
// For the tag value, the top byte indicates the data type of the data (either
// string or number) and the remaining 7 bytes are the actual tag value.
type Callback func(t *Event, data []byte, tag uint64)
