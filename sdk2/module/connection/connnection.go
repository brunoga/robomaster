package connection

import (
	"log/slog"
	"net"
	"sync/atomic"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/support/finder"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
	"github.com/brunoga/robomaster/unitybridge/unity/event"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

const (
	subTypeConnectionOpen = iota
	subTypeConnectionClose
	subTypeConnectionSetIP
	subTypeConnectionSetPort

	wifiDirectIPString = "192.168.2.1"
)

// Connection provides support for managing the connection to the robot.
type Connection struct {
	*internal.BaseModule

	appID uint64
	typ   Type

	f *finder.Finder

	signalQuality atomic.Uint64
}

var _ module.Module = (*Connection)(nil)

// New creates a new Connection instance with the given UnityBridge instance and
// logger.
func New(ub unitybridge.UnityBridge,
	l *logger.Logger, appID uint64, typ Type) (*Connection, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("connection_module").With(
		slog.Uint64("app_id", appID))

	return &Connection{
		BaseModule: internal.NewBaseModule(ub, l, "Connection",
			key.KeyAirLinkConnection, nil),
		appID: appID,
		typ:   typ,
		f:     finder.New(appID, l),
	}, nil
}

// Start starts the connection module. It will try to find a robot broadcasting
// in the network and connect to it.
func (c *Connection) Start() error {
	err := c.BaseModule.Start()
	if err != nil {
		return err
	}

	var ip net.IP = net.ParseIP(wifiDirectIPString)
	if c.typ == TypeRouter {
		b, err := c.f.Find(30 * time.Second)
		if err != nil {
			return err
		}

		c.f.SendACK(b.SourceIp(), b.AppId())

		ip = b.SourceIp()
	}

	e := event.NewFromType(event.TypeConnection)

	e.ResetSubType(subTypeConnectionClose)
	err = c.UB().SendEvent(e)
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionSetIP)
	err = c.UB().SendEventWithString(e, ip.String())
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionSetPort)
	err = c.UB().SendEventWithUint64(e, 10607)
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionOpen)
	err = c.UB().SendEvent(e)
	if err != nil {
		return err
	}

	c.UB().AddKeyListener(key.KeyAirLinkSignalQuality, func(r *result.Result) {
		c.signalQuality.Store(r.Value().(*value.Uint64).Value)
	}, false)

	return nil
}

// SignalQualityLevel returns the current signal quality level. 0 means no
// signal whatsoever and 60 appears to be the strongest value.
func (c *Connection) SignalQualityLevel() uint8 {
	return uint8(c.signalQuality.Load())
}

// SignalQualityBars returns the current signal quality as a number of bars (1
// to 4).
func (c *Connection) SignalQualityBars() uint8 {
	level := c.SignalQualityLevel()

	if level < 10 {
		return 1
	}
	if level < 25 {
		return 2
	}
	if level < 45 {
		return 3
	}

	return 4
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
