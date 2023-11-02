package unitybridge

import (
	"github.com/brunoga/unitybridge/event"
	"github.com/brunoga/unitybridge/internal/implementations"
)

// UnityBridge is the interface to the UnityBridge library. It will always wrap
// the specific implementation of the UnityBridge librart for the specific
// platform it is being compiled on (currently there are native implementations
// for windows_amd64, darwin_amd64, ios_arm64, android_arm64 and android_arm.
// There is also a linux_amd64 implementation through Wine).
type UnityBridge interface {
	// Create sets up the UnityBridge using the given name (apparently only
	// "Robomaster" is supported), debuggable status (true to enable log
	// debugging) and path for log files.
	Create(name string, debuggable bool, logPath string)

	// Initialize tries to initialize the UnityBridge. Returns true if
	// successful.
	Initialize() bool

	// SetEventCallback sets the callback function for given eventCode.
	// Note that eventCode can also represent an event type (by making sure
	// subType is zero) and then any events with that type will be sent to the
	// given callback.
	SetEventCallback(eventCode uint64, callback event.Callback)

	// SendEvent sends an event with given eventCode, data, and tag. The data
	// field in this case will usually be uintptr(0) to indicate no data. As
	// there is no way to express the actual length of the data being sent.
	SendEvent(eventCode uint64, data uintptr, tag uint64)

	// SendEventWithString sends an event with given eventCode, string data, and
	// tag.
	SendEventWithString(eventCode uint64, data string, tag uint64)

	// SendEventWithNumber sends an event with given eventCode, number data, and
	/// tag.
	SendEventWithNumber(eventCode, data, tag uint64)

	// GetSecurityKeyByKeyChainIndex returns the security key associated with
	// the given index.
	GetSecurityKeyByKeyChainIndex(index int) string

	// Unitialize uninitializes the UnityBridge.
	Uninitialize()

	// Destroy destroys the UnityBridge.
	Destroy()
}

// Get returns a platform specific singleton instance of the UnityBridge interface.
func Get() UnityBridge {
	return UnityBridge(implementations.UnityBridgeImpl)
}
