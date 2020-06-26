package notification

type connection interface {
	Open() error
	Read([]byte) (int, error)
	Close() error
}
