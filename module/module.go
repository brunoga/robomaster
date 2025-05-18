package module

import (
	"fmt"
	"time"

	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
)

// Module is the interface implemented by all modules.
type Module interface {
	fmt.Stringer

	// Add a listener for the given key. Listeners will be automatically
	// removed when the module is stopped.
	AddKeyListener(*key.Key, func(*result.Result), bool) error

	// Start starts the module.
	Start() error

	// Connected returns true if the module connection has been established.
	Connected() bool

	// WaitForConnection waits for the module connection to be established.
	WaitForConnection(timeout time.Duration) bool

	// Stop stops the module.
	Stop() error
}
