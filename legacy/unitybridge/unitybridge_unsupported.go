//go:build !(darwin && amd64) && !(android && arm) && !(android && arm64) && !(ios && arm64) && !(windows && amd64)

package unitybridge

/*
#include <unitybridge.h>
*/
import "C"
import (
	"fmt"
	"runtime"
)

func init() {
	// We are on an unsupported platform, so panic as there is nothing we can do.
	panic(fmt.Sprintf("Platform \"%s/%s\" not supported by Unity Bridge",
		runtime.GOOS, runtime.GOARCH))
}

// The following functions are stubs for unsupported platforms. This allow code
// to compile but it will crash immediatelly when run due ton the panic above.

func Create(string, bool, string)                  {}
func Destroy()                                     {}
func Initialize() bool                             { return false }
func Uninitialize()                                {}
func SetEventCallback(uint64, C.EventCallbackFunc) {}
func SendEvent(uint64, []byte, uint64)             {}
func SendEventWithString(uint64, string, uint64)   {}
func SendEventWithNumber(uint64, uint64, uint64)   {}
func GetSecurityKeyByKeyChainIndex(int) uintptr    { return uintptr(0) }
