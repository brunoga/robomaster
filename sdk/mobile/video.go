package mobile

import (
	"sync"

	"github.com/brunoga/robomaster/sdk/modules/video"
)

type Video struct {
	v *video.Video
}

type VideoHandler interface {
	Handle([]byte, *sync.WaitGroup)
}

func (v *Video) StartStream(videoHandler VideoHandler) (int, error) {
	return v.v.StartStream(videoHandler.Handle)
}

func (v *Video) StopStream(token int) error {
	return v.v.StopStream(token)
}
