package event

// Type represents a Unity Bridge event type.
type Type uint32

const (
	TypeSetValue          Type = iota // Sets the value of an attribute.
	TypeGetValue                      // Gets the value of an attribute.
	TypeGetAvailableValue             // Get cached value of an attribute.
	TypePerformAction                 // Perform an action.
	TypeStartListening                // Start listening for a specific event.
	TypeStopListening                 // Stop listening for a specific event.
	TypeActivation
	TypeLocalAlbum
	TypeFirmwareUpgrade

	TypeConnection         Type = 100 // Configure connection.
	TypeSecurity           Type = 101
	TypePrintLog           Type = 200
	TypeStartVideo         Type = 300 // Start video streaming.
	TypeStopVideo          Type = 301 // Stop video streaming.
	TypeRender             Type = 302
	TypeGetNativeTexture   Type = 303
	TypeVideoTransferSpeed Type = 304 // Video speed updates.
	TypeAudioDataRecv      Type = 305
	TypeVideoDataRecv      Type = 306
	TypeNativeFunctions    Type = 500
)

func (t Type) String() string {
	switch t {
	case TypeSetValue:
		return "SetValue"
	case TypeGetValue:
		return "GetValue"
	case TypeGetAvailableValue:
		return "GetAvailableValue"
	case TypePerformAction:
		return "PerformAction"
	case TypeStartListening:
		return "StartListening"
	case TypeStopListening:
		return "StopListening"
	case TypeActivation:
		return "Activation"
	case TypeLocalAlbum:
		return "LocalAlbum"
	case TypeFirmwareUpgrade:
		return "FirmwareUpgrade"
	case TypeConnection:
		return "Connection"
	case TypeSecurity:
		return "Security"
	case TypePrintLog:
		return "PrintLog"
	case TypeStartVideo:
		return "StartVideo"
	case TypeStopVideo:
		return "StopVideo"
	case TypeRender:
		return "Render"
	case TypeGetNativeTexture:
		return "GetNativeTexture"
	case TypeVideoTransferSpeed:
		return "VideoTransferSpeed"
	case TypeAudioDataRecv:
		return "AudioDataRecv"
	case TypeVideoDataRecv:
		return "VideoDataRecv"
	case TypeNativeFunctions:
		return "NativeFunctions"
	}

	return "Unknown"
}

func AllTypes() []Type {
	return []Type{
		TypeSetValue,
		TypeGetValue,
		TypeGetAvailableValue,
		TypePerformAction,
		TypeStartListening,
		TypeStopListening,
		TypeActivation,
		TypeLocalAlbum,
		TypeFirmwareUpgrade,
		TypeConnection,
		TypeSecurity,
		TypePrintLog,
		TypeStartVideo,
		TypeStopVideo,
		TypeRender,
		TypeGetNativeTexture,
		TypeVideoTransferSpeed,
		TypeAudioDataRecv,
		TypeVideoDataRecv,
		TypeNativeFunctions,
	}
}
