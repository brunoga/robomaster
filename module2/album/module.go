package album

import (
	"fmt"
	"log/slog"

	"github.com/brunoga/robomaster/module2"
	"github.com/brunoga/robomaster/module2/internal"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/support/token"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/event"
)

type Module struct {
	*internal.BaseModule

	t token.Token
}

var _ module2.Module = (*Module)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger) *Module {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("album_module")

	return &Module{
		BaseModule: internal.NewBaseModule(ub, l, "Album",
			nil),
	}
}

func (m *Module) Start() error {
	err := m.BaseModule.Start()
	if err != nil {
		return err
	}

	e := event.NewFromTypeAndSubType(event.TypeLocalAlbum, 0)

	err = m.UB().SendEvent(e)
	if err != nil {
		return err
	}

	e.ResetSubType(1)
	err = m.UB().SendEventWithString(e, "[PATH]")
	if err != nil {
		return err
	}

	m.t, err = m.UB().AddEventTypeListener(event.TypeLocalAlbum, m.onLocalAlbumEvent)

	return err
}

func (m *Module) Stop() error {
	err := m.StopRecordingVideo()
	if err != nil {
		return err
	}

	err = m.UB().RemoveEventTypeListener(event.TypeLocalAlbum, m.t)
	if err != nil {
		return err
	}

	return m.BaseModule.Stop()
}

func (m *Module) StartRecordingVideo() error {
	e := event.NewFromTypeAndSubType(event.TypeLocalAlbum, 2)

	return m.UB().SendEvent(e)
}

func (m *Module) StopRecordingVideo() error {
	e := event.NewFromTypeAndSubType(event.TypeLocalAlbum, 3)

	return m.UB().SendEvent(e)
}

func (m *Module) TakePhoto() error {
	e := event.NewFromTypeAndSubType(event.TypeLocalAlbum, 4)

	return m.UB().SendEvent(e)
}

func (m *Module) DeleteFiles() error {
	return fmt.Errorf("implement me")
}

func (m *Module) LoadOriginalFile() error {
	return fmt.Errorf("implement me")
}

func (m *Module) LoadAlbumFiles() error {
	return fmt.Errorf("implement me")
}

func (m *Module) onLocalAlbumEvent(e *event.Event, data []byte, dataType event.DataType) {
	m.L().Debug("onLocalAlbumEvent", "event", e, "data", data, "dataType", dataType)
	switch e.SubType() {
	case 6:
		// Download. data = filename (string)
	case 7:
		// Record time. data = time in seconds (int64)
	case 8:
		// Fetch thumbnail success. data = nil
	case 9:
		// Fetch photo success. data = key (string)
	case 10:
		// Fetch photo failure. data = key (string)
	default:
		m.L().Error("Unexpected subtype", "subType", e.SubType())
	}
}
