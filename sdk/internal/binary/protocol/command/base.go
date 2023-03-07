package command

// base is the base struct for all commands. It keeps track of the set, id and
// type of the command and has a raw data buffer for the command itself.
type base struct {
	set byte
	id  byte
	typ Type

	data []byte
}

// newBase creates a new base struct for a command with the given set, id, type
// and size of the raw data buffer.
func newBase(set, id byte, typ Type, size int) *base {
	return &base{
		set:  set,
		id:   id,
		typ:  typ,
		data: make([]byte, size),
	}
}

// Set returns the set of the command.
func (b *base) Set() byte {
	return b.set
}

// ID returns the id of the command.
func (b *base) ID() byte {
	return b.id
}

// Type returns the type of the command.
func (b *base) Type() Type {
	return b.typ
}
