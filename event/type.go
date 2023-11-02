//go:generate go run github.com/dmarkham/enumer -type=Type

package event

// Type represents a Unity Bridge event type.
type Type int32

const (
	SetValue          Type = iota // Sets the value of an attribute.
	GetValue                      // Gets the value of an attribute.
	GetAvailableValue             // Get cached value of an attribute.
	PerformAction                 // Perform an action.
	StartListening                // Start listening for a specific event.
	StopListening                 // Stop listening for a specific event.
	Activation
	LocalAlbum
	FirmwareUpgrade

	Connection         Type = 100 // Configure connection.
	Security           Type = 101
	PrintLog           Type = 200
	StartVideo         Type = 300 // Start video streaming.
	StopVideo          Type = 301 // Stop video streaming.
	Render             Type = 302
	GetNativeTexture   Type = 303
	VideoTransferSpeed Type = 304 // Video speed updates.
	AudioDataRecv      Type = 305
	VideoDataRecv      Type = 306
	NativeFunctions    Type = 500
)
