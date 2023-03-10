package video

import (
	"fmt"
	"image"
	"net"
	"sync"

	"github.com/brunoga/robomaster/sdk/internal/h264"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/control"
	"github.com/brunoga/robomaster/sdk/modules/video"
)

const (
	videoAddrPort = ":40921"
)

// Video handles starting a robot's video stream, receiving and decoding the
// data from it and sending to all registered VideoHandlers and stopping the
// video stream. The decoding relies on GoCV (https://gocv.io).
type Video struct {
	control *control.Control
	decoder *h264.Decoder

	m            sync.Mutex
	quitChan     chan struct{}
	waitChan     chan struct{}
	videoHandler video.Handler
}

// New creates a new Video instance. The control parameter is used to start
// stop the video stream and setup the video connection address.
func New(control *control.Control) (*Video, error) {
	v := &Video{
		control,
		nil,
		sync.Mutex{},
		nil,
		nil,
		nil,
	}

	decoder, err := h264.NewDecoder(v.frameCallback)
	if err != nil {
		return nil, fmt.Errorf("error creating h264 decoder: %w", err)
	}

	v.decoder = decoder

	return v, nil
}

func (v *Video) Start(resolution video.Resolution,
	videoHandler video.Handler) error {
	v.m.Lock()
	defer v.m.Unlock()

	if videoHandler == nil {
		return fmt.Errorf("video handler must not be nil")
	}

	v.videoHandler = videoHandler
	v.quitChan = make(chan struct{})
	v.waitChan = make(chan struct{})

	go v.videoLoop()

	return nil
}

func (v *Video) Stop() error {
	v.m.Lock()
	defer v.m.Unlock()

	if v.videoHandler == nil {
		return fmt.Errorf("already stopped")
	}

	close(v.quitChan)

	<-v.waitChan

	v.videoHandler = nil
	v.quitChan = nil
	v.waitChan = nil

	return nil
}

func (v *Video) videoLoop() {
	err := v.control.SendDataExpectOk("stream on;")
	if err != nil {
		// TODO(bga): Log this.
		return
	}

	ip, err := v.control.IP()
	if err != nil {
		// TODO(bga): Log this.
		return
	}

	videoAddr := ip.String() + videoAddrPort

	videoConn, err := net.Dial("tcp", videoAddr)
	if err != nil {
		return
	}
	defer videoConn.Close()

	err = v.decoder.Open()
	if err != nil {
		return
	}

	readBuffer := make([]byte, 16384)

L:
	for {
		select {
		case <-v.quitChan:
			break L
		default:
			n, err := videoConn.Read(readBuffer)
			if err != nil {
				break L
			}

			v.decoder.SendData(readBuffer[:n])
		}
	}

	err = v.control.SendDataExpectOk("stream off;")
	if err != nil {
		// TODO(bga): Log this.
		return
	}

	_ = v.decoder.Close()

	close(v.waitChan)
}

func (v *Video) frameCallback(data []byte, width, height int) {
	frameRGBA := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})

	copy(frameRGBA.Pix, data)

	v.videoHandler(frameRGBA)
}
