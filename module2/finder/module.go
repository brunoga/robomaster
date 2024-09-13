package finder

import (
	"log/slog"

	"github.com/brunoga/robomaster/module2"
	"github.com/brunoga/robomaster/module2/internal"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/support/token"
)

type Module struct {
	*internal.BaseModule

	tg                        token.Generator
	broadcastMessageCallbacks map[token.Token]func(*Broadcast)
}

var _ module2.Module = &Module{}

func New(l *logger.Logger, appID uint64) *Module {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("finder_module").With("app_id", appID)

	return &Module{
		BaseModule: internal.NewBaseModule(nil, l, "Finder", nil),
	}
}

func (m *Module) AddBroadcastMessageCallback(callback func(*Broadcast)) {
}
