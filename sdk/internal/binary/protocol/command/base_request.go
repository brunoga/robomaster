package command

// baseRequest is the base struct for all request (Push) commands. It includes
// the implementation of the Data method.
type baseRequest struct {
	*base
}

// newBaseRequest creates a new baseRequest struct for a command with the given
// set, id, type and size of the raw data buffer.
func newBaseRequest(set, id byte, typ Type, size int) *baseRequest {
	return &baseRequest{
		base: newBase(set, id, typ, size),
	}
}

// Data returns a reference to the wire data for the command. The returned slice
// should not be changed.
func (r *baseRequest) Data() []byte {
	return r.data
}
