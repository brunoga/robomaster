package wrapper

const (
	FuncCreateUnityBridge byte = iota
	FuncDestroyUnityBridge
	FuncUnityBridgeInitialize
	FuncUnityBridgeUninitialize
	FuncUnitySetEventCallback
	FuncUnitySendEvent
	FuncUnitySendEventWithNumber
	FuncUnitySendEventWithString
)
