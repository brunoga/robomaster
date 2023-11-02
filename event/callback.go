package event

// Callback is the prototype for functions that will be called when a reply to
// an event is received.
type Callback func(eventCode uint64, data []byte, tag uint64)
