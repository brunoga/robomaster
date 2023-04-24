package video

import (
	"fmt"
	"net"
	"sync"

	"github.com/brunoga/robomaster/sdk/modules/control"
)

const (
	videoAddrPort = ":40921"
)

// VideoHandler is a handler for video streams. A handler should do its work and
// return as fast as possible. The data passed in is the encoded h264 stream
// directly from the robot and needs to be decoded before being used.
//
// After doing its work, a VideoHandler must call wg.Done() before returning.
type Handler func(data []byte, wg *sync.WaitGroup)

// Video handles starting a robot's video stream, receiving the encoded h264
// data from it and sending to all registered VideoHandlers and stopping the
// video stream.
type Video struct {
	control *control.Control

	m             sync.Mutex
	quitChan      chan struct{}
	videoHandlers map[int]Handler
}

// New creates a new Video instance. The control parameter is used to start
// stop the video stream and setup the video connection address.
func New(control *control.Control) (*Video, error) {
	return &Video{
		control,
		sync.Mutex{},
		nil,
		make(map[int]Handler),
	}, nil
}

// StartStream starts the video stream (if it has not started yet) and starts
// sending video frames to the given videoHandler. It returns a positive int
// token (used to stop the stream when needed) and a nil error on success and
// a non-nil error on failure.
func (v *Video) StartStream(videoHandler Handler) (int, error) {
	v.m.Lock()
	defer v.m.Unlock()

	if len(v.videoHandlers) == 0 {
		go v.videoLoop()
	}

	for i := 0; i < len(v.videoHandlers)+1; i++ {
		_, ok := v.videoHandlers[i]
		if ok {
			continue
		}

		v.videoHandlers[i] = videoHandler

		return i, nil
	}

	return -1, fmt.Errorf("video handler tokens exhausted")
}

// StopStream stops sending frames to the VideoHandler associated with the
// given token and remove it from the list of VideoHandlers. If it is the last
// VideoHandler in the list, the robot's video stream is stopped.
func (v *Video) StopStream(token int) error {
	v.m.Lock()
	defer v.m.Unlock()

	_, ok := v.videoHandlers[token]
	if !ok {
		return fmt.Errorf("invalid stream handler token")
	}

	delete(v.videoHandlers, token)

	if len(v.videoHandlers) == 0 {
		close(v.quitChan)
	}

	return nil
}

func (v *Video) videoLoop() {
	v.quitChan = make(chan struct{})

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
			var wg sync.WaitGroup

			// Send frame to all video handlers.
			v.m.Lock()
			for _, videoHandler := range v.videoHandlers {
				wg.Add(1)
				go videoHandler(readBuffer[:n], &wg)
			}
			v.m.Unlock()

			// Wait for all video handlers to notify they finished processing
			// the frame.
			wg.Wait()
		}
	}

	err = v.control.SendDataExpectOk("stream off;")
	if err != nil {
		// TODO(bga): Log this.
		return
	}

	v.quitChan = nil
}
