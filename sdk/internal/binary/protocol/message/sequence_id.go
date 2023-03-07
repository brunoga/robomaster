package message

import "sync"

const (
	sequenceIDMin uint16 = 10000
	sequenceIDMax uint16 = 20000
)

var (
	m              sync.Mutex
	nextSequenceID = sequenceIDMin
)

func getSequenceID() uint16 {
	m.Lock()
	nextSequenceID = sequenceIDMin + (nextSequenceID-sequenceIDMin+1)%(sequenceIDMax-sequenceIDMin)
	sequenceID := nextSequenceID
	m.Unlock()

	return sequenceID
}
