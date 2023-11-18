package module

import "fmt"

// Module is the interface implemented by all modules.
type Module interface {
	fmt.Stringer

	Start() error

	WaitForConnection() bool

	Stop() error
}
