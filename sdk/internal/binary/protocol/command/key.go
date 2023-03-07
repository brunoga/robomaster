package command

// Key returns an identifier for the command based on its set and id.
func Key(cmdSet, cmdID byte) uint16 {
	return (uint16(cmdSet) << 8) | uint16(cmdID)
}
