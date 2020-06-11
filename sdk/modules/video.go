package modules

import (
	"fmt"
	"sync"

	"gocv.io/x/gocv"
)

const (
	videoAddrPort = ":40921"
)

// VideoHandler is a handler for video streams. A handler should do its work and
// return as fast as possible. If a handler needs to modify the given frame in
// any way, it should first make a copy of it as the given frame is shared by
// all handlers.
//
// Even if you are not modifying the given frame but you are doing some expensive
// image processing, it might be better to copy the Mat internally and return
// ASAP, deferring the processing to a separate goroutine (you might need to
// queue frames).
//
// After doing its work, a VideoHandler must call wg.Done() before returning.
type VideoHandler func(frame *gocv.Mat, wg *sync.WaitGroup)

// Video handles starting a robot's video stream, receiving and decoding the
// data from it and sending to all registered VideoHandlers and stopping the
// video stream. The decoding relies on GoCV (https://gocv.io).
type Video struct {
	control *Control

	m sync.Mutex
	quitChan chan struct{}
	videoHandlers map[int]VideoHandler
}

// NewVideo creates a new Video instance. The control parameter is used to start
// stop the video stream and setup the video connection address.
func NewVideo(control *Control) *Video {
	return &Video{
		control,
		sync.Mutex{},
		nil,
		make(map[int]VideoHandler),
	}
}

// StartStream starts the video stream (if it has not started yet) and starts
// sending video frames to the given videoHandler. It returns a positive int
// token (used to stop the stream when needed) and a nil error on success and
// a non-nil error on failure.
func (v *Video) StartStream(videoHandler VideoHandler) (int, error) {
	v.m.Lock()
	defer v.m.Unlock()

	if len(v.videoHandlers) == 0 {
		go v.videoLoop()
	}

	for i := 0; i < len(v.videoHandlers) + 1; i++ {
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

	// Use GoCV for decoding the video stream from the robot.
	stream, err := gocv.VideoCaptureFile("tcp://" + videoAddr)
	if err != nil {
		fmt.Println("VideoCaptureFile error", err)
		return
	}
	defer stream.Close()

	frame := gocv.NewMat()
	defer frame.Close()

L:
	for {
		select {
		case <-v.quitChan:
			break L
		default:
			if !stream.Read(&frame) {
				break L
			}

			var wg sync.WaitGroup

			// Send frame to all video handlers.
			v.m.Lock()
			for _, videoHandler := range v.videoHandlers {
				wg.Add(1)
				go videoHandler(&frame, &wg)
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