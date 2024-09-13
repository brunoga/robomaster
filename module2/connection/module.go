package connection

import (
	"log/slog"
	"net"

	"github.com/brunoga/robomaster/module2"
	"github.com/brunoga/robomaster/module2/internal"
	"github.com/brunoga/robomaster/support/finder"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/support/token"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/event"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
)

type Module struct {
	f     *finder.Finder
	appID uint64
	t     Type

	*internal.BaseModule
}

var _ module2.Module = (*Module)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger, f *finder.Finder,
	appID uint64, t Type) *Module {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("connection_module").With("app_id", appID)

	return &Module{
		f:     f,
		appID: appID,
		t:     t,
		BaseModule: internal.NewBaseModule(ub, l, "Connection",
			key.KeyAirLinkConnection),
	}
}

func (m *Module) Start() error {
	if err := m.BaseModule.Start(); err != nil {
		return err
	}

	var findCh chan *finder.Broadcast
	err := m.f.StartFinding(findCh)
	if err != nil {
		return err
	}

	var ch chan *finder.Broadcast
	err = m.f.StartFinding(ch)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			err := m.f.StopFinding()
			if err != nil {
				m.L().Error("Error stopping finding", "error", err)
			}
		}()

		bm := <-ch

		m.f.SendACK(bm.SourceIp(), bm.AppId())

		err = m.sendSetIP(bm.SourceIp())
		if err != nil {
			m.L().Error("Error sending set ip", "ip", bm.SourceIp().String(), "error", err)
			return
		}

		err = m.sendSetPort(10607)
		if err != nil {
			m.L().Error("Error sending set port", "port", 10607, "error", err)
			return
		}

		err = m.sendOpen()
		if err != nil {
			m.L().Error("Error sending open", "error", err)
		}
	}()

	return nil
}

func (m *Module) AddSignalQualityCallback(callback result.Callback) (token.Token, error) {
	return m.UB().AddKeyListener(key.KeyAirLinkSignalQuality, callback, true)
}

func (m *Module) RemoveSignalQualityCallback(token token.Token) error {
	return m.UB().RemoveKeyListener(key.KeyAirLinkSignalQuality, token)
}

func (m *Module) Stop() error {
	// Remove all signal quality callbacks.
	err := m.RemoveSignalQualityCallback(0)
	if err != nil {
		return err
	}

	return m.BaseModule.Stop()
}

var connectionEvent = event.NewFromType(event.TypeConnection)

func (m *Module) sendOpen() error {
	connectionEvent.ResetSubType(uint32(actionOpen))
	return m.UB().SendEvent(connectionEvent)
}

func (m *Module) sendClose() error {
	connectionEvent.ResetSubType(uint32(actionClose))
	return m.UB().SendEvent(connectionEvent)
}

func (m *Module) sendSetIP(ip net.IP) error {
	connectionEvent.ResetSubType(uint32(actionSetIP))
	return m.UB().SendEventWithString(connectionEvent, ip.String())
}

func (m *Module) sendSetPort(port uint64) error {
	connectionEvent.ResetSubType(uint32(actionSetPort))
	return m.UB().SendEventWithUint64(connectionEvent, port)
}
