package wrapper

// EventCallback is the type required for functions that want to handle bridge
// events.
type EventCallback func(eventCode uint64, data []byte, tag uint64)

// Wrapper is a Go wrapper for DJI's Robomaster unitybridge.dll functions.
// Descriptions of methods below are mostly what could be seem through
// observation and might be incomplete or even incorrect. Also, the interface
// does not cover all exported methods (more specifically, it does not cover
// GetRenderEventFunc, UnityPluginLoad and UnityPluginUnload), only the ones
// actually used to interface with the Robomaster S1.
type Wrapper interface {
	// CreateUnityBridge creates the app/Robomaster bridge. The name
	// parameter *MUST* be "Robomaster" otherwise the bridge will simply not
	// do anything it seems. The debuggable parameter enables or disables
	// debugging to the console/terminal. This has to be called before any
	// other methods.
	CreateUnityBridge(name string, debuggable bool)

	// DestroyUnityBridge removes the bridge. It should only be called after
	// the bridge was unitialized (or never initialized).
	DestroyUnityBridge()

	// Initializes the bridge, effectivelly starting it. Returns true on
	// success or false on failure.
	UnityBridgeInitialize() bool

	// UnityBridgeUninitialize uninitializes the bridge.
	UnityBridgeUninitialize()

	// UnitySetEventCallback sets a callback function for a specific event
	// code.
	UnitySetEventCallback(eventCode uint64, eventCallback EventCallback)

	// UnitySendEvent sends an event wieth the given eventCode, data and
	// tag to the bridge. The tag parameter is used to track responses to a
	// specific request (i.e. An event sent with tag set to X will result
	// in one or more callbacks with tag also set to X).
	UnitySendEvent(eventCode uint64, data []byte, tag uint64)

	// UnitySendEventWithString same as above, but data is a string.
	UnitySendEventWithString(eventCode uint64, data string, tag uint64)

	// UnitySendEventWithString same as above, but data is a uint64.
	UnitySendEventWithNumber(eventCode uint64, data uint64, tag uint64)
}
