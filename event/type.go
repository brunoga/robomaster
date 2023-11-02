//go:generate go run github.com/dmarkham/enumer -type=Type

package event

type Type int32

const (
	SetValue Type = iota
	GetValue
	GetAvailableValue
	PerformAction
	StartListening
	StopListening
	Activation
	LocalAlbum
	FirmwareUpgrade

	Connection         Type = 100
	Security           Type = 101
	PrintLog           Type = 200
	StartVideo         Type = 300
	StopVideo          Type = 301
	Render             Type = 302
	GetNativeTexture   Type = 303
	VideoTransferSpeed Type = 304
	AudioDataRecv      Type = 305
	VideoDataRecv      Type = 306
	NativeFunctions    Type = 500
)
