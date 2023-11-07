package key

type AccessType int32

const (
	AccessTypeNone = iota
	AccessTypeRead
	AccessTypeWrite
	AccessTypeAction
)
