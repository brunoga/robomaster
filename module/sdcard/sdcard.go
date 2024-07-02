package sdcard

import (
	"fmt"
	"log/slog"

	"github.com/brunoga/robomaster/module/camera"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

type Module struct {
	*internal.BaseModule
}

// New creates a new SDCard module instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	conm *connection.Connection, camm *camera.Camera) (*Module, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("sdcard_module")

	m := &Module{}

	m.BaseModule = internal.NewBaseModule(ub, l, "SDCard", nil, func(r *result.Result) {
		if !r.Succeeded() {
			m.Logger().Error("Connection: Unsuccessfull result.", "result", r)
			return
		}

		connectedValue, ok := r.Value().(*value.Bool)
		if !ok {
			m.Logger().Error("Connection: Unexpected value.", "value", r.Value())
			return
		}

		if connectedValue.Value {
			m.Logger().Debug("Connected.")
		} else {
			m.Logger().Debug("Disconnected.")
		}
	}, conm, camm)

	return m, nil
}

func (m *Module) IsInserted() (bool, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraSDCardIsInserted, false)
	if err != nil {
		return false, err
	}

	if !r.Succeeded() {
		return false, fmt.Errorf("failed to get value for key: %v", r)
	}

	v, ok := r.Value().(*value.Bool)
	if !ok {
		return false, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return v.Value, nil
}

func (m *Module) Format() error {
	err := m.UB().PerformActionForKeySync(key.KeyCameraFormatSDCard, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) IsFormatting() (bool, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraSDCardIsFormatting, false)
	if err != nil {
		return false, err
	}

	if !r.Succeeded() {
		return false, fmt.Errorf("failed to get value for key: %v", r)
	}

	v, ok := r.Value().(*value.Bool)
	if !ok {
		return false, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return v.Value, nil
}

func (m *Module) IsFull() (bool, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraSDCardIsFull, false)
	if err != nil {
		return false, err
	}

	if !r.Succeeded() {
		return false, fmt.Errorf("failed to get value for key: %v", r)
	}

	v, ok := r.Value().(*value.Bool)
	if !ok {
		return false, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return v.Value, nil
}

func (m *Module) HasError() (bool, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraSDCardHasError, false)
	if err != nil {
		return false, err
	}

	if !r.Succeeded() {
		return false, fmt.Errorf("failed to get value for key: %v", r)
	}

	v, ok := r.Value().(*value.Bool)
	if !ok {
		return false, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return v.Value, nil
}

func (m *Module) TotalSpaceInMB() (uint64, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraSDCardTotalSpaceInMB, false)
	if err != nil {
		return 0, err
	}

	if !r.Succeeded() {
		return 0, fmt.Errorf("failed to get value for key: %v", r)
	}

	v, ok := r.Value().(*value.Uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return v.Value, nil
}

func (m *Module) RemainingSpaceInMB() (uint64, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraSDCardRemainingSpaceInMB, false)
	if err != nil {
		return 0, err
	}

	if !r.Succeeded() {
		return 0, fmt.Errorf("failed to get value for key: %v", r)
	}

	v, ok := r.Value().(*value.Uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return v.Value, nil
}

func (m *Module) AvailablePhotoCount() (uint64, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraSDCardAvailablePhotoCount, false)
	if err != nil {
		return 0, err
	}

	if !r.Succeeded() {
		return 0, fmt.Errorf("failed to get value for key: %v", r)
	}

	v, ok := r.Value().(*value.Uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return v.Value, nil
}

func (m *Module) AvailableRecordingTimeInSeconds() (uint64, error) {
	r, err := m.UB().GetKeyValueSync(key.KeyCameraSDCardAvailableRecordingTimeInSeconds, false)
	if err != nil {
		return 0, err
	}

	if !r.Succeeded() {
		return 0, fmt.Errorf("failed to get value for key: %v", r)
	}

	v, ok := r.Value().(*value.Uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return v.Value, nil
}
