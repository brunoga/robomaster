package command

// Request is the interface that must be implemented by all request (Push)
// commands. It allows access to the wire date for the command which is
// used for sending it to the robot.
type Request interface {
	Command

	// Data returns a reference to the wire data for the command. The
	// returned slice should not be changed.
	Data() []byte
}
