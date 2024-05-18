package mobile

import (
	"github.com/brunoga/robomaster/unitybridge/support/token"

	"github.com/brunoga/robomaster/module/camera"
)

const (
	CameraHorizontalResolutionPoints = camera.HorizontalResolutionPoints
	CameraVerticalResolutionPoints   = camera.VerticalResolutionPoints
	CameraHorizontalFOVDegrees       = camera.HorizontalFOVDegrees
	CameraVerticalFOVDegrees         = camera.VerticalFOVDegrees
)

// VideoHandler is the interface that must be implemented by types that want
// to handle video frames from the camera.
type VideoHandler interface {
	// HandleVideo is called when a new video frame is received from the camera.
	// rgb24Data, as the name name implies, is the raw data for a RGB24 image and
	// its dimensions are 1280x720.
	HandleVideo(rgb24Data []byte)
}

// Camera allows controlling the robot camera.
type Camera struct {
	c *camera.Camera
}

// AddVideoHandler adds a new video handler to the camera. If this is the first
// video handler added, the camera will start sending video frames.
func (c *Camera) AddVideoHandler(handler VideoHandler) (int64, error) {
	t, err := c.c.AddVideoCallback(func(frame *camera.RGB) {
		handler.HandleVideo(frame.Pix)
	})

	return int64(t), err
}

// RemoveVideoHandler removes a video handler from the camera. If this is the
// last video handler removed, the camera will stop sending video frames.
func (c *Camera) RemoveVideoHandler(t int64) error {
	return c.c.RemoveVideoCallback(token.Token(t))
}

// StartRecordingVideo starts recording video from the camera to the robot's
// SD card.
func (c *Camera) StartRecordingVideo() error {
	return c.c.StartRecordingVideo()
}

// IsRecordingVideo returns true if the camera is currently recording video.
func (c *Camera) IsRecordingVideo() (bool, error) {
	return c.c.IsRecordingVideo()
}

// RecordingTimeInSeconds returns the current recording time in seconds.
func (c *Camera) RecordingTimeInSeconds() int64 {
	return int64(c.c.RecordingTime().Seconds())
}

// StopRecordingVideo stops recording video from the camera to the robot's SD
// card.
func (c *Camera) StopRecordingVideo() error {
	return c.c.StopRecordingVideo()
}
