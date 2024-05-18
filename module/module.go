package module

import (
	"fmt"
	"time"
)

// Module is the interface implemented by all modules.
type Module interface {
	fmt.Stringer

	// Start starts the module.
	Start() error

	// Connected returns true if the module connection has been established.
	Connected() bool

	// WaitForConnection waits for the module connection to be established.
	WaitForConnection(timeout time.Duration) bool

	// Stop stops the module.
	Stop() error
}
