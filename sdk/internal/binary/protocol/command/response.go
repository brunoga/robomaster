package command

// Response is the interface that must be implemented by all responses to a
// request (Pull) command. It allows checking if the command was successful
// or not.
type Response interface {
	Command

	// Ok returns true if the command was successful.
	Ok() bool
}
