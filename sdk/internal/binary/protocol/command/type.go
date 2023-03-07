package command

// Type is the type of command.
type Type byte

const (
	// TypePull is for commands that are a initiated from client and sent to the
	// robot. It always have an associated response.
	TypePull Type = iota
	// TypePush is for commands that are initiated by the robot and sent to the
	// client. There is no associated response.
	TypePush
)
