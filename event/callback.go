package event

// Callback is the prototype for functions that will be called when a reply to
// an event is received.
//
// For the tag value, it appears that only the bottom 6 bytes are used. From
// the other 2 bytes, it appears only the bottom one is currently used to hold
// information about the data type of the event (0 for string and 1 for uint64)
// but only when the event refers to a event type only and not to an actual
// specific event.
type Callback func(eventCode uint64, data []byte, tag uint64)
