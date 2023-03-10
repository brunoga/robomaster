package video

type Video interface {
	Start(Resolution, Handler) error
	Stop() error
}
