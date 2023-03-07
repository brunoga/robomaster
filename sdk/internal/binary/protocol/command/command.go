package command

// Command is the interface that must be implemented by all commands. It
// contains the bare minimum functionality.
type Command interface {
	// Set returns the command set.
	Set() byte
	// ID returns the command ID.
	ID() byte
	// Type returns the command type.
	Type() Type

	// New returns a new instance of the command. Used mostly for creating
	// commands based on a received response.
	New() Command
}
