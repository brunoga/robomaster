//go:generate go run github.com/dmarkham/enumer -type=Type

package event

// Type represents a Unity Bridge event type.
type Type int32

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
