package camera

import (
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

// Camera provides support for managing the camera attached to the robot.
type Camera struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger

	gntToken token.Token
	vtsToken token.Token
	vdrToken token.Token

	ccToken token.Token
	crToken token.Token

	tg token.Generator

	connRL *support.ResultListener

	m             sync.RWMutex
	callbacks     map[token.Token]VideoCallback
	recordingTime time.Duration
}

var _ module.Module = (*Camera)(nil)

// New creates a new Camera instance with the given UnityBridge instance and
// logger.
func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Camera, error) {
	return &Camera{
		ub:        ub,
		l:         l,
		tg:        *token.NewGenerator(),
		callbacks: make(map[token.Token]VideoCallback),
		connRL:    support.NewResultListener(ub, l, key.KeyCameraConnection),
	}, nil
}

// Start starts the camera manager.
func (c *Camera) Start() error {
	var err error

	err = c.connRL.Start(func(r *result.Result) {
		if r.ErrorCode() != 0 {
			return
		}

		if r.Value().(bool) {
			// Ask for video texture information.
			err = c.ub.SendEvent(event.NewFromType(event.TypeGetNativeTexture))
			if err != nil {
				return
			}
		}
	})
	if err != nil {
		return err
	}

	c.gntToken, err = c.ub.AddEventTypeListener(event.TypeGetNativeTexture,
		c.onGetNativeTexture)
	if err != nil {
		return err
	}

	c.vtsToken, err = c.ub.AddEventTypeListener(event.TypeVideoTransferSpeed,
		c.onVideoTransferSpeed)
	if err != nil {
		return err
	}

	c.vdrToken, err = c.ub.AddEventTypeListener(event.TypeVideoDataRecv,
		c.onVideoDataRecv)
	if err != nil {
		return err
	}

	return nil
}

func (c *Camera) WaitForConnection() bool {
	connected, ok := c.connRL.Result().Value().(bool)
	if ok && connected {
		return true
	}

	return c.connRL.WaitForNewResult(5 * time.Second).Value().(bool)
}

// Connected returns whether the camera is physically attached to the robot and
// ready to use. Note this is based on information provided by the robot itself.
func (c *Camera) Connected() bool {
	ok, connected := c.connRL.Result().Value().(bool)
	if !ok {
		return false
	}

	return connected
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

	if len(c.callbacks) == 0 {
		err := c.ub.SendEvent(event.NewFromType(event.TypeStartVideo))
		if err != nil {
			return 0, err
		}
	}

	t := c.tg.Next()

	c.callbacks[t] = vc

	return t, nil
}

// RemoveVideoCallback removes the callback function associated with the given
// token.
func (c *Camera) RemoveVideoCallback(t token.Token) error {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.callbacks) == 0 {
		return fmt.Errorf("no callbacks added")
	}

	_, ok := c.callbacks[t]
	if !ok {
		return fmt.Errorf("no callback added for token %d", t)
	}

	delete(c.callbacks, t)

	if len(c.callbacks) == 0 {
		err := c.ub.SendEvent(event.NewFromType(event.TypeStopVideo))
		if err != nil {
			return err
		}
	}

	return nil
}

// VideoFormat returns the currently set video format.
func (c *Camera) VideoFormat() (VideoFormat, error) {
	var value VideoFormat

	return value, c.ub.GetKeyValueSync(key.KeyCameraVideoFormat, true, &value)
}

// SetVideoFormat sets the video resolution.
//
// TODO(bga): Other then  actually limiting the available resolutions, it looks
// like changing resolutions is not working. Need to investigate further as
// there might be some setup that is needed and is not being done. It might be
// that this is only for the video recorded in the robot and not for the
// video being streamed from it.
func (c *Camera) SetVideoFormat(format VideoFormat) error {
	return c.ub.SetKeyValueSync(key.KeyCameraVideoFormat, format)
}

// VideoQuality returns the currently set video quality.
func (c *Camera) VideoQuality() (VideoQuality, error) {
	var value VideoQuality

	return value, c.ub.GetKeyValueSync(key.KeyCameraVideoTransRate,
		true, &value)
}

// SetVideoQuality sets the video quality.
func (c *Camera) SetVideoQuality(quality VideoQuality) error {
	return c.ub.SetKeyValueSync(key.KeyCameraVideoTransRate, quality)
}

// Mode returns the current camera mode.
func (c *Camera) Mode() (Mode, error) {
	var value Mode

	return value, c.ub.GetKeyValueSync(key.KeyCameraMode, true, &value)
}

// SetMode sets the camera mode.
func (c *Camera) SetMode(mode Mode) error {
	return c.ub.SetKeyValueSync(key.KeyCameraMode, mode)
}

// ExposureMode returns the current digital zoom factor.
func (c *Camera) DigitalZoomFactor() (uint64, error) {
	var value uint64

	return value, c.ub.GetKeyValueSync(key.KeyCameraDigitalZoomFactor,
		true, &value)
}

// SetDigitalZoomFactor sets the digital zoom factor.
func (c *Camera) SetDigitalZoomFactor(factor uint64) error {
	return c.ub.SetKeyValueSync(key.KeyCameraDigitalZoomFactor, factor)
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

	err = c.ub.PerformActionForKeySync(key.KeyCameraStartRecordVideo, nil)
	if err != nil {
		return err
	}

	c.crToken, err = c.ub.AddKeyListener(key.KeyCameraCurrentRecordingTimeInSeconds,
		func(r *result.Result) {
			if r.ErrorCode() != 0 {
				c.l.Warn("error getting current recording time", "err", r.ErrorDesc())
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
	var value bool

	return value, c.ub.GetKeyValueSync(key.KeyCameraIsRecording, true, &value)
}

// RecordingTime returns the current recording time in seconds.
func (c *Camera) RecordingTime() time.Duration {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.recordingTime
}

// StopRecordingVideo stops recording video to the robot's internal storage.
func (c *Camera) StopRecordingVideo() error {
	err := c.ub.PerformActionForKeySync(key.KeyCameraStopRecordVideo, nil)
	if err != nil {
		return err
	}

	return c.ub.RemoveKeyListener(key.KeyCameraCurrentRecordingTimeInSeconds,
		c.crToken)
}

// Stop stops the camera manager.
func (c *Camera) Stop() error {
	c.m.Lock()

	if len(c.callbacks) > 0 {
		c.callbacks = make(map[token.Token]VideoCallback)

		err := c.ub.SendEvent(event.NewFromType(event.TypeStopVideo))
		if err != nil {
			c.m.Unlock()
			return err
		}
	}

	c.m.Unlock()

	err := c.connRL.Stop()
	if err != nil {
		return err
	}

	err = c.ub.RemoveEventTypeListener(event.TypeGetNativeTexture, c.gntToken)
	if err != nil {
		return err
	}

	err = c.ub.RemoveEventTypeListener(event.TypeVideoTransferSpeed, c.vtsToken)
	if err != nil {
		return err
	}

	err = c.ub.RemoveEventTypeListener(event.TypeVideoDataRecv, c.vdrToken)
	if err != nil {
		return err
	}

	return nil
}

func (c *Camera) String() string {
	return "Camera"
}

func (c *Camera) onGetNativeTexture(data []byte, dataType event.DataType) {
	c.l.Debug("onGetNativeTexture", "data", data, "dataType", dataType)
}

func (c *Camera) onVideoTransferSpeed(data []byte, dataType event.DataType) {
	c.l.Debug("onVideoTransferSpeed", "data", data, "dataType", dataType)
}

func (c *Camera) onVideoDataRecv(data []byte, dataType event.DataType) {
	//c.l.Debug("onVideoDataRecv", "len(data)", len(data), "dataType", dataType)

	c.m.RLock()

	rgb := NewRGBFromBytes(data, image.Rect(0, 0, 1280, 720))

	for _, vc := range c.callbacks {
		go vc(rgb)
	}

	c.m.RUnlock()
}
