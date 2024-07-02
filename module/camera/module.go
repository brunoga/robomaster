package camera

import (
	"encoding/json"
	"fmt"
	"image"
	"log/slog"
	"sync"
	"sync/atomic"
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

// Module provides support for managing the camera attached to the robot.
type Module struct {
	*internal.BaseModule

	gntToken token.Token
	vtsToken token.Token
	vdrToken token.Token

	crToken token.Token

	tg *token.Generator

	recordingTime atomic.Pointer[time.Duration]

	glTextureData atomic.Pointer[value.GLTextureData]

	m         sync.RWMutex
	callbacks map[token.Token]VideoCallback
}

var _ module.Module = (*Module)(nil)

// New creates a new Camera instance with the given UnityBridge instance and
// logger.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection) (*Module, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("camera_module")

	m := &Module{
		tg:        token.NewGenerator(),
		callbacks: make(map[token.Token]VideoCallback),
	}

	m.BaseModule = internal.NewBaseModule(ub, l, "Camera",
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
				l.Debug("Camera Connected.")
			} else {
				l.Debug("Camera Disconnected.")
			}
		}, cm)

	return m, nil
}

// Start starts the camera module.
func (m *Module) Start() error {
	var err error

	m.gntToken, err = m.UB().AddEventTypeListener(event.TypeGetNativeTexture,
		m.onGetNativeTexture)
	if err != nil {
		return err
	}

	m.vtsToken, err = m.UB().AddEventTypeListener(event.TypeVideoTransferSpeed,
		m.onVideoTransferSpeed)
	if err != nil {
		return err
	}

	m.vdrToken, err = m.UB().AddEventTypeListener(event.TypeVideoDataRecv,
		m.onVideoDataRecv)
	if err != nil {
		return err
	}

	return m.BaseModule.Start()
}

// AddVideoCallback adds a callback function to be called when a new video frame
// is received from the robot. The callback function will be called in a
// separate goroutine. Returns a token that can be used to remove the callback
// later.
func (m *Module) AddVideoCallback(vc VideoCallback) (token.Token, error) {
	if vc == nil {
		return 0, fmt.Errorf("callback must not be nil")
	}

	m.m.Lock()
	defer m.m.Unlock()

	t := m.tg.Next()

	m.callbacks[t] = vc

	if len(m.callbacks) == 1 {
		// We just added the first callback. Start video stream.
		err := m.UB().SendEvent(event.NewFromType(event.TypeStartVideo))
		if err != nil {
			return 0, err
		}
	}

	return t, nil
}

// RemoveVideoCallback removes the callback function associated with the given
// token.
func (m *Module) RemoveVideoCallback(t token.Token) error {
	m.m.Lock()
	defer m.m.Unlock()

	_, ok := m.callbacks[t]
	if !ok {
		return fmt.Errorf("no callback added for token %d", t)
	}

	delete(m.callbacks, t)

	if len(m.callbacks) == 0 {
		// We just removed the last callback. Stop video stream.
		err := m.UB().SendEvent(event.NewFromType(event.TypeStopVideo))
		if err != nil {
			return err
		}
	}

	return nil
}

// VideoFormat returns the currently set video format.
func (m *Module) VideoFormat() (VideoFormat, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraVideoFormat, true)
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
func (m *Module) SetVideoFormat(format VideoFormat) error {
	return m.UB().SetKeyValueSync(key.KeyCameraVideoFormat, format)
}

// SetVideoQuality sets the video quality.
func (m *Module) SetVideoQuality(quality VideoQuality) error {
	return m.UB().SetKeyValueSync(key.KeyCameraVideoTransRate, &value.Float64{Value: float64(quality)})
}

// Mode returns the current camera mode.
func (m *Module) Mode() (Mode, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraMode, true)
	if err != nil {
		return 0, err
	}

	return Mode(r.Value().(*value.Uint64).Value), nil
}

// SetMode sets the camera mode.
func (m *Module) SetMode(mode Mode) error {
	return m.UB().SetKeyValueSync(key.KeyCameraMode,
		&value.Uint64{Value: uint64(mode)})
}

// ExposureMode returns the current digital zoom factor.
func (m *Module) DigitalZoomFactor() (uint64, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraDigitalZoomFactor,
		true)
	if err != nil {
		return 0, err
	}

	return uint64(r.Value().(float64)), nil
}

// SetDigitalZoomFactor sets the digital zoom factor.
func (m *Module) SetDigitalZoomFactor(factor uint64) error {
	return m.UB().SetKeyValueSync(key.KeyCameraDigitalZoomFactor, factor)
}

// StartRecordingVideo starts recording video to the robot's internal storage.
func (m *Module) StartRecordingVideo() error {
	var err error

	currentMode, err := m.Mode()
	if err != nil {
		return err
	}

	if currentMode != ModeVideo {
		err = m.SetMode(ModeVideo)
		if err != nil {
			return err
		}
	}

	err = m.UB().PerformActionForKeySync(key.KeyCameraStartRecordVideo, nil)
	if err != nil {
		return err
	}

	m.crToken, err = m.UB().AddKeyListener(
		key.KeyCameraCurrentRecordingTimeInSeconds,
		func(r *result.Result) {
			if !r.Succeeded() {
				m.Logger().Error("error getting current recording time", "error",
					r.ErrorDesc())
				return
			}

			duration := time.Duration(r.Value().(*value.Uint64).Value) * time.Second
			m.recordingTime.Store(&duration)
		}, true)

	return err
}

// IsRecordingVideo returns whether the robot is currently recording video to
// its internal storage.
func (m *Module) IsRecordingVideo() (bool, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraIsRecording, true)
	if err != nil {
		return false, err
	}

	return r.Value().(*value.Bool).Value, nil
}

// RecordingTime returns the current recording time in seconds.
func (m *Module) RecordingTime() time.Duration {
	return *m.recordingTime.Load()
}

// StopRecordingVideo stops recording video to the robot's internal storage.
func (m *Module) StopRecordingVideo() error {
	err := m.UB().PerformActionForKeySync(key.KeyCameraStopRecordVideo, nil)
	if err != nil {
		return err
	}

	return m.UB().RemoveKeyListener(key.KeyCameraCurrentRecordingTimeInSeconds,
		m.crToken)
}

// RenderNextFrame requests the next frame to be rendered. This is used by iOS
// and the frame will be rendered to a texture associated with an OpenGLES 2.0
// context that was current when Start() is called. This should be called for
// for each frame to be rendered (up to 60 times per second).
func (m *Module) RenderNextFrame() {
	m.UB().RenderNextFrame()
}

// GLTextureData returns information about the current texture used for
// rendering frames. See RenderNextFrame() above.
func (c *Module) GLTextureData() (value.GLTextureData, error) {
	glTextureData := c.glTextureData.Load()
	if glTextureData == nil || *glTextureData == (value.GLTextureData{}) {
		return value.GLTextureData{}, fmt.Errorf("no GLTextureData available. Did " +
			"you call RenderNextFrame?")
	}

	return *glTextureData, nil
}

// Stop stops the camera manager.
func (m *Module) Stop() error {
	m.m.Lock()

	if len(m.callbacks) > 0 {
		m.callbacks = make(map[token.Token]VideoCallback)

		err := m.UB().SendEvent(event.NewFromType(event.TypeStopVideo))
		if err != nil {
			m.m.Unlock()
			return err
		}
	}

	m.m.Unlock()

	err := m.UB().RemoveEventTypeListener(event.TypeGetNativeTexture,
		m.gntToken)
	if err != nil {
		return err
	}

	err = m.UB().RemoveEventTypeListener(event.TypeVideoTransferSpeed,
		m.vtsToken)
	if err != nil {
		return err
	}

	err = m.UB().RemoveEventTypeListener(event.TypeVideoDataRecv, m.vdrToken)
	if err != nil {
		return err
	}

	return m.BaseModule.Stop()
}

func (m *Module) onGetNativeTexture(data []byte, dataType event.DataType) {
	endTrace := m.Logger().Trace("onGetNativeTexture", "data", string(data), "dataType", dataType)
	defer endTrace()

	var glTextureData value.GLTextureData
	err := json.Unmarshal(data, &glTextureData)
	if err != nil {
		m.Logger().Error("onGetNativeTexture", "error", err)
		return
	}

	m.glTextureData.Store(&glTextureData)
}

func (m *Module) onVideoTransferSpeed(data []byte, dataType event.DataType) {
	m.Logger().Debug("onVideoTransferSpeed", "data", data, "dataType", dataType)
}

func (m *Module) onVideoDataRecv(data []byte, dataType event.DataType) {
	rgb := NewRGBFromBytes(data, image.Rect(0, 0, 1280, 720))

	m.m.RLock()

	for _, vc := range m.callbacks {
		go vc(rgb)
	}

	m.m.RUnlock()
}
