package notification

type Connection interface {
	Open() error
	Read([]byte) (int, error)
	Close() error
}
