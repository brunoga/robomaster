package callback

// Callback is the prototype for functions that will be called when a reply to
// an registered event type is received.
//
// For the tag value, the top byte indicates the data type of the data (either
// string, 0,  or number, 1) and the remaining 7 bytes are the actual tag value.
type Callback func(eventCode uint64, data []byte, tag uint64)
