package key

// AccessType is the type of access allowed for a specific key.
type AccessType int32

const (
	AccessTypeNone = iota
	AccessTypeRead = 1 << (iota - 1)
	AccessTypeWrite
	AccessTypeAction
)
