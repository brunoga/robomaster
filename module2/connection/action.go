package connection

type action int8

const (
	actionOpen action = iota
	actionClose
	actionSetIP
	actionSetPort
	actionCount
)
