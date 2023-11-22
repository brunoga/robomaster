package connection

import (
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/finder"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
)

const (
	subTypeConnectionOpen = iota
	subTypeConnectionClose
	subTypeConnectionSetIP
	subTypeConnectionSetPort
)

// Connection provides support for managing the connection to the robot.
type Connection struct {
	*internal.BaseModule

	appID uint64

	f *finder.Finder
}

var _ module.Module = (*Connection)(nil)

// New creates a new Connection instance with the given UnityBridge instance and
// logger.
func New(ub unitybridge.UnityBridge,
	l *logger.Logger, appID uint64) (*Connection, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("connection_module").With(
		slog.Uint64("app_id", appID))

	return &Connection{
		BaseModule: internal.NewBaseModule(ub, l, "Connection",
			key.KeyRobomasterSystemConnection, nil),
		appID: appID,
		f:     finder.New(appID, l),
	}, nil
}

// Start starts the connection module. It will try to find a robot broadcasting
// in the network and connect to it.
func (cm *Connection) Start() error {
	err := cm.BaseModule.Start()
	if err != nil {
		return err
	}

	b, err := cm.f.Find(30 * time.Second)
	if err != nil {
		return err
	}

	cm.f.SendACK(b.SourceIp(), b.AppId())

	e := event.NewFromType(event.TypeConnection)

	e.ResetSubType(subTypeConnectionClose)
	err = cm.UB().SendEvent(e)
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionSetIP)
	err = cm.UB().SendEventWithString(e, b.SourceIp().String())
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionSetPort)
	err = cm.UB().SendEventWithUint64(e, 10607)
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionOpen)
	err = cm.UB().SendEvent(e)
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the connection module.
func (cm *Connection) Stop() error {
	e := event.NewFromType(event.TypeConnection)

	e.ResetSubType(subTypeConnectionClose)
	err := cm.UB().SendEvent(e)
	if err != nil {
		return err
	}

	return cm.BaseModule.Stop()
}
