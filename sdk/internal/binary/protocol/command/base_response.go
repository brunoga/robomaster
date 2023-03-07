package command

// baseResponse is the base struct for all responses to request (Push) commands.
// It includes the implementation of the Ok method.
type baseResponse struct {
	*base
}

// newBaseResponse creates a new baseResponse struct for a command with the given
// set, id, type and size of the raw data buffer.
func newBaseResponse(set, id byte, typ Type, size int) *baseResponse {
	return &baseResponse{
		base: newBase(set, id, typ, size),
	}
}

// Ok returns true if the command was successful.
func (b *baseResponse) Ok() bool {
	return b.data[0] == 0x00
}
