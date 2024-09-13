package module2

import (
	"fmt"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/support/token"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
)

// Module is the interface implemented by all modules.
type Module interface {
	fmt.Stringer

	// Start starts the module.
	Start() error

	// Connected returns true if the module connection has been established.
	Connected() bool

	// WaitForConnectionStatus waits for the module connection status to be the
	// given one (connected or disconnected).
	WaitForConnectionStatus(connected bool)

	// AddConnectionCallback adds a connection callback to the module. The given
	// callback will be called whenever the module connection status changes.
	AddConnectionCallback(callback result.Callback) (token.Token, error)

	// RemoveConnectionCallback removes a connection callback from the module.
	RemoveConnectionCallback(token token.Token) error

	// Stop stops the module. All connection callbacks are removed.
	Stop() error

	// UB returns the Unity Bridge used by the module.
	UB() unitybridge.UnityBridge

	// L returns the logger used by the module.
	L() *logger.Logger
}
