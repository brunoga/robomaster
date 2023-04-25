package mobile

import (
	"sync"

	"github.com/brunoga/robomaster/sdk/modules/video"
)

type Video struct {
	v *video.Video
}

type VideoHandler interface {
	Handle([]byte, *WaitGroup)
}

func (v *Video) StartStream(videoHandler VideoHandler) (int, error) {
	h := func(data []byte, wg *sync.WaitGroup) {
		videoHandler.Handle(data, &WaitGroup{wg: wg})
	}
	return v.v.StartStream(h)
}

func (v *Video) StopStream(token int) error {
	return v.v.StopStream(token)
}
