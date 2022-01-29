package dji

import "fmt"

type Key int

const (
	KeyNone Key = iota
	KeyAirLinkConnection
	KeyGimbalAngleFrontYawRotation
	KeyGimbalAngleFrontPitchRotation
	KeyGimbalConnection
	KeyGimbalControlMode
	KeyGimbalOpenAttitudeUpdates
	KeyGimbalResetPosition
	KeyMainControllerConnection
	KeyMainControllerVirtualStick
	KeyRobomasterOpenChassisSpeedUpdates
	KeyRobomasterSystemConnection
	KeyRobomasterSystemFunctionEnable
	KeysCount
	// TODO(bga): Add keys here as needed.
)

type DataType int

const (
	KeyDataTypeBool DataType = iota
	KeyDataTypeLong
	KeyDataTypeAbsoluteRotationParameter
	KeyDataTypeControllerParameter
	KeyDataTypeFunctionEnableParameter
	KeyDataTypeVoid
	// TODO(bga): Add data types here as needed.
)

type AccessType int

const (
	KeyAccessTypeNone AccessType = 0
	KeyAccessTypeRead AccessType = 1 << (iota - 1)
	KeyAccessTypeWrite
	KeyAccessTypeAction
)

func (k Key) DataType() DataType {
	return keyAttributeMap[k].dataType
}

func (k Key) Value() uint32 {
	return keyAttributeMap[k].value
}

func (k Key) AccessType() AccessType {
	return keyAttributeMap[k].accessType
}

type keyAttributes struct {
	value      uint32
	dataType   DataType
	accessType AccessType
}

var (
	keyAttributeMap = map[Key]keyAttributes{
		KeyAirLinkConnection: keyAttributes{117440513, KeyDataTypeBool,
			KeyAccessTypeRead},
		KeyGimbalConnection: keyAttributes{67108865, KeyDataTypeBool,
			KeyAccessTypeRead},
		KeyGimbalControlMode: keyAttributes{67108869, KeyDataTypeLong,
			KeyAccessTypeRead | KeyAccessTypeWrite},
		KeyGimbalAngleFrontYawRotation: keyAttributes{67108876,
			KeyDataTypeAbsoluteRotationParameter,
			KeyAccessTypeAction},
		KeyGimbalAngleFrontPitchRotation: keyAttributes{67108877,
			KeyDataTypeAbsoluteRotationParameter,
			KeyAccessTypeAction},
		KeyGimbalOpenAttitudeUpdates: keyAttributes{67108882,
			KeyDataTypeVoid, KeyAccessTypeAction},
		KeyGimbalResetPosition: keyAttributes{67108870, KeyDataTypeLong,
			KeyAccessTypeAction},
		KeyMainControllerConnection: keyAttributes{33554433,
			KeyDataTypeBool, KeyAccessTypeRead},
		KeyMainControllerVirtualStick: keyAttributes{33554436,
			KeyDataTypeControllerParameter, KeyAccessTypeAction},
		KeyRobomasterOpenChassisSpeedUpdates: keyAttributes{33554474,
			KeyDataTypeVoid, KeyAccessTypeAction},
		KeyRobomasterSystemConnection: keyAttributes{83886081,
			KeyDataTypeBool, KeyAccessTypeRead},
		KeyRobomasterSystemFunctionEnable: keyAttributes{83886147,
			KeyDataTypeFunctionEnableParameter,
			KeyAccessTypeAction},
		// TODO(bga): Add other attributes here as needed. Needs to be
		//  kept in sync with existing Keys.
	}

	keyByValueMap = map[int]Key{
		117440513: KeyAirLinkConnection,
		67108865:  KeyGimbalConnection,
		67108869:  KeyGimbalControlMode,
		67108876:  KeyGimbalAngleFrontYawRotation,
		67108877:  KeyGimbalAngleFrontPitchRotation,
		67108882:  KeyGimbalOpenAttitudeUpdates,
		67108870:  KeyGimbalResetPosition,
		33554433:  KeyMainControllerConnection,
		33554436:  KeyMainControllerVirtualStick,
		33554474:  KeyRobomasterOpenChassisSpeedUpdates,
		83886081:  KeyRobomasterSystemConnection,
		83886147:  KeyRobomasterSystemFunctionEnable,
	}
)

func keyByValue(value int) Key {
	key, ok := keyByValueMap[value]
	if !ok {
		panic(fmt.Sprintf("Can't get key for value %d.", value))
	}

	return key
}
