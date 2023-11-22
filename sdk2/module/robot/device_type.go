package robot

// DeviceType is the type of a device connected to the robot.
type DeviceType int16

const (
	DeviceTypeImageTransmission DeviceType = 256
	DeviceTypeCamera            DeviceType = 260
	DeviceTypeChassis           DeviceType = 768
	DeviceTypeBattery           DeviceType = 778
	DeviceTypeESC0              DeviceType = 788
	DeviceTypeESC1              DeviceType = 789
	DeviceTypeESC2              DeviceType = 790
	DeviceTypeESC3              DeviceType = 791
	DeviceTypeServo1            DeviceType = 792
	DeviceTypeServo2            DeviceType = 793
	DeviceTypeServo3            DeviceType = 794
	DeviceTypeServo4            DeviceType = 795
	DeviceTypeArm               DeviceType = 798
	DeviceTypeGimbal            DeviceType = 1024
	DeviceTypeTOF1              DeviceType = 4609
	DeviceTypeTOF2              DeviceType = 4610
	DeviceTypeTOF3              DeviceType = 4611
	DeviceTypeTOF4              DeviceType = 4612
	DeviceTypeSensorAdapter1    DeviceType = 5633
	DeviceTypeSensorAdapter2    DeviceType = 5634
	DeviceTypeSensorAdapter3    DeviceType = 5635
	DeviceTypeSensorAdapter4    DeviceType = 5636
	DeviceTypeSensorAdapter5    DeviceType = 5637
	DeviceTypeSensorAdapter6    DeviceType = 5638
	DeviceTypeWaterGun          DeviceType = 5888
	DeviceTypeInfraredGun       DeviceType = 5889
	DeviceTypeBackArmor         DeviceType = 6145
	DeviceTypeFrontArmor        DeviceType = 6146
	DeviceTypeLeftArmor         DeviceType = 6147
	DeviceTypeRightArmor        DeviceType = 6148
	DeviceTypeLeftHeadArmor     DeviceType = 6149
	DeviceTypeRightHeadArmor    DeviceType = 6150
)

func (d DeviceType) String() string {
	switch d {
	case DeviceTypeImageTransmission:
		return "ImageTransmission"
	case DeviceTypeCamera:
		return "Camera"
	case DeviceTypeChassis:
		return "Chassis"
	case DeviceTypeBattery:
		return "Battery"
	case DeviceTypeESC0:
		return "ESC0"
	case DeviceTypeESC1:
		return "ESC1"
	case DeviceTypeESC2:
		return "ESC2"
	case DeviceTypeESC3:
		return "ESC3"
	case DeviceTypeServo1:
		return "Servo1"
	case DeviceTypeServo2:
		return "Servo2"
	case DeviceTypeServo3:
		return "Servo3"
	case DeviceTypeServo4:
		return "Servo4"
	case DeviceTypeArm:
		return "Arm"
	case DeviceTypeGimbal:
		return "Gimbal"
	case DeviceTypeTOF1:
		return "TOF1"
	case DeviceTypeTOF2:
		return "TOF2"
	case DeviceTypeTOF3:
		return "TOF3"
	case DeviceTypeTOF4:
		return "TOF4"
	case DeviceTypeSensorAdapter1:
		return "SensorAdapter1"
	case DeviceTypeSensorAdapter2:
		return "SensorAdapter2"
	case DeviceTypeSensorAdapter3:
		return "SensorAdapter3"
	case DeviceTypeSensorAdapter4:
		return "SensorAdapter4"
	case DeviceTypeSensorAdapter5:
		return "SensorAdapter5"
	case DeviceTypeSensorAdapter6:
		return "SensorAdapter6"
	case DeviceTypeWaterGun:
		return "WaterGun"
	case DeviceTypeInfraredGun:
		return "InfraredGun"
	case DeviceTypeBackArmor:
		return "BackArmor"
	case DeviceTypeFrontArmor:
		return "FrontArmor"
	case DeviceTypeLeftArmor:
		return "LeftArmor"
	case DeviceTypeRightArmor:
		return "RightArmor"
	case DeviceTypeLeftHeadArmor:
		return "LeftHeadArmor"
	case DeviceTypeRightHeadArmor:
		return "RightHeadArmor"
	default:
		return "Unknown"
	}
}
