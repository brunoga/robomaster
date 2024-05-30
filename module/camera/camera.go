package camera

import (
	"fmt"
	"image"
	"log/slog"
	"sync"
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/support/token"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/event"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

// Camera provides support for managing the camera attached to the robot.
type Camera struct {
	*internal.BaseModule

	gntToken token.Token
	vtsToken token.Token
	vdrToken token.Token

	crToken token.Token

	tg *token.Generator

	m             sync.RWMutex
	callbacks     map[token.Token]VideoCallback
	recordingTime time.Duration
}

var _ module.Module = (*Camera)(nil)

// New creates a new Camera instance with the given UnityBridge instance and
// logger.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection) (*Camera, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("camera_module")

	c := &Camera{
		tg:        token.NewGenerator(),
		callbacks: make(map[token.Token]VideoCallback),
	}

	c.BaseModule = internal.NewBaseModule(ub, l, "Camera",
		key.KeyCameraConnection, func(r *result.Result) {
			if !r.Succeeded() {
				l.Error("Camera Connection: Unsuccessfull result.", "result", r)
				return
			}

			connectedValue, ok := r.Value().(*value.Bool)

			if !ok {
				l.Error("Camera Connection: Unexpected value.", "value", r.Value())
				return
			}

			if connectedValue.Value {
				// Ask for video texture information.
				if err := c.UB().SendEvent(event.NewFromType(
					event.TypeGetNativeTexture)); err != nil {
					l.Error("error sending event", "event",
						event.TypeGetNativeTexture.String(), "error", err)
				}
				l.Debug("Camera Connected.")
			} else {
				l.Debug("Camera Disconnected.")
			}
		}, cm)

	return c, nil
}

// Start starts the camera manager.
func (c *Camera) Start() error {
	var err error

	c.gntToken, err = c.UB().AddEventTypeListener(event.TypeGetNativeTexture,
		c.onGetNativeTexture)
	if err != nil {
		return err
	}

	c.vtsToken, err = c.UB().AddEventTypeListener(event.TypeVideoTransferSpeed,
		c.onVideoTransferSpeed)
	if err != nil {
		return err
	}

	c.vdrToken, err = c.UB().AddEventTypeListener(event.TypeVideoDataRecv,
		c.onVideoDataRecv)
	if err != nil {
		return err
	}

	return c.BaseModule.Start()
}

// AddVideoCallback adds a callback function to be called when a new video frame
// is received from the robot. The callback function will be called in a
// separate goroutine. Returns a token that can be used to remove the callback
// later.
func (c *Camera) AddVideoCallback(vc VideoCallback) (token.Token, error) {
	if vc == nil {
		return 0, fmt.Errorf("callback must not be nil")
	}

	c.m.Lock()
	defer c.m.Unlock()

	t := c.tg.Next()

	c.callbacks[t] = vc

	if len(c.callbacks) == 1 {
		// We just added the first callback. Start video stream.
		err := c.UB().SendEvent(event.NewFromType(event.TypeStartVideo))
		if err != nil {
			return 0, err
		}
	}

	return t, nil
}

// RemoveVideoCallback removes the callback function associated with the given
// token.
func (c *Camera) RemoveVideoCallback(t token.Token) error {
	c.m.Lock()
	defer c.m.Unlock()

	_, ok := c.callbacks[t]
	if !ok {
		return fmt.Errorf("no callback added for token %d", t)
	}

	delete(c.callbacks, t)

	if len(c.callbacks) == 0 {
		// We just removed the last callback. Stop video stream.
		err := c.UB().SendEvent(event.NewFromType(event.TypeStopVideo))
		if err != nil {
			return err
		}
	}

	return nil
}

// VideoFormat returns the currently set video format.
func (c *Camera) VideoFormat() (VideoFormat, error) {
	r, err := c.UB().GetKeyValueSync(key.KeyCameraVideoFormat, true)
	if err != nil {
		return 0, err
	}

	return VideoFormat(r.Value().(float64)), nil
}

// SetVideoFormat sets the video resolution.
//
// TODO(bga): Other then  actually limiting the available resolutions, it looks
// like changing resolutions is not working. Need to investigate further as
// there might be some setup that is needed and is not being done. It might be
// that this is only for the video recorded in the robot and not for the
// video being streamed from it.
func (c *Camera) SetVideoFormat(format VideoFormat) error {
	return c.UB().SetKeyValueSync(key.KeyCameraVideoFormat, format)
}

// SetVideoQuality sets the video quality.
func (c *Camera) SetVideoQuality(quality VideoQuality) error {
	return c.UB().SetKeyValueSync(key.KeyCameraVideoTransRate, &value.Float64{Value: float64(quality)})
}

// Mode returns the current camera mode.
func (c *Camera) Mode() (Mode, error) {
	r, err := c.UB().GetKeyValueSync(key.KeyCameraMode, true)
	if err != nil {
		return 0, err
	}

	return Mode(r.Value().(float64)), nil
}

// SetMode sets the camera mode.
func (c *Camera) SetMode(mode Mode) error {
	return c.UB().SetKeyValueSync(key.KeyCameraMode, mode)
}

// ExposureMode returns the current digital zoom factor.
func (c *Camera) DigitalZoomFactor() (uint64, error) {
	r, err := c.UB().GetKeyValueSync(key.KeyCameraDigitalZoomFactor,
		true)
	if err != nil {
		return 0, err
	}

	return uint64(r.Value().(float64)), nil
}

// SetDigitalZoomFactor sets the digital zoom factor.
func (c *Camera) SetDigitalZoomFactor(factor uint64) error {
	return c.UB().SetKeyValueSync(key.KeyCameraDigitalZoomFactor, factor)
}

// StartRecordingVideo starts recording video to the robot's internal storage.
func (c *Camera) StartRecordingVideo() error {
	var err error

	currentMode, err := c.Mode()
	if err != nil {
		return err
	}

	if currentMode != 1 {
		err = c.SetMode(1)
		if err != nil {
			return err
		}
	}

	err = c.UB().PerformActionForKeySync(key.KeyCameraStartRecordVideo, nil)
	if err != nil {
		return err
	}

	c.crToken, err = c.UB().AddKeyListener(
		key.KeyCameraCurrentRecordingTimeInSeconds,
		func(r *result.Result) {
			if r.ErrorCode() != 0 {
				c.Logger().Warn("error getting current recording time", "error",
					r.ErrorDesc())
			}

			c.m.Lock()
			c.recordingTime = time.Duration(r.Value().(float64)) * time.Second
			c.m.Unlock()
		}, true)

	return err
}

// IsRecordingVideo returns whether the robot is currently recording video to
// its internal storage.
func (c *Camera) IsRecordingVideo() (bool, error) {
	r, err := c.UB().GetKeyValueSync(key.KeyCameraIsRecording, true)
	if err != nil {
		return false, err
	}

	return r.Value().(bool), nil
}

// RecordingTime returns the current recording time in seconds.
func (c *Camera) RecordingTime() time.Duration {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.recordingTime
}

// StopRecordingVideo stops recording video to the robot's internal storage.
func (c *Camera) StopRecordingVideo() error {
	err := c.UB().PerformActionForKeySync(key.KeyCameraStopRecordVideo, nil)
	if err != nil {
		return err
	}

	return c.UB().RemoveKeyListener(key.KeyCameraCurrentRecordingTimeInSeconds,
		c.crToken)
}

// Stop stops the camera manager.
func (c *Camera) Stop() error {
	c.m.Lock()

	if len(c.callbacks) > 0 {
		c.callbacks = make(map[token.Token]VideoCallback)

		err := c.UB().SendEvent(event.NewFromType(event.TypeStopVideo))
		if err != nil {
			c.m.Unlock()
			return err
		}
	}

	c.m.Unlock()

	err := c.UB().RemoveEventTypeListener(event.TypeGetNativeTexture,
		c.gntToken)
	if err != nil {
		return err
	}

	err = c.UB().RemoveEventTypeListener(event.TypeVideoTransferSpeed,
		c.vtsToken)
	if err != nil {
		return err
	}

	err = c.UB().RemoveEventTypeListener(event.TypeVideoDataRecv, c.vdrToken)
	if err != nil {
		return err
	}

	return c.BaseModule.Stop()
}

func (c *Camera) onGetNativeTexture(data []byte, dataType event.DataType) {
	c.Logger().Debug("onGetNativeTexture", "data", data, "dataType", dataType)
}

func (c *Camera) onVideoTransferSpeed(data []byte, dataType event.DataType) {
	c.Logger().Debug("onVideoTransferSpeed", "data", data, "dataType", dataType)
}

func (c *Camera) onVideoDataRecv(data []byte, dataType event.DataType) {
	rgb := NewRGBFromBytes(data, image.Rect(0, 0, 1280, 720))

	c.m.RLock()

	for _, vc := range c.callbacks {
		go vc(rgb)
	}

	c.m.RUnlock()
}
