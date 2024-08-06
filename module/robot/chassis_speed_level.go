package robot

// ChassisSpeedLevel represents the speed level of the chassis.
type ChassisSpeedLevel int8

const (
	// ChassisSpeedLevelFast is the fast speed level.
	ChassisSpeedLevelFast ChassisSpeedLevel = iota
	// ChassisSpeedLevelMedium is the medium speed level.
	ChassisSpeedLevelMedium
	// ChassisSpeedLevelSlow is the slow speed level.
	ChassisSpeedLevelSlow
	// ChassisSpeedLevelCustom is the custom speed level. This enables
	// individually setting speeds for each direction/axis. This is not
	// supported yet.
	//
	// TODO(bga): Add support for this.
	ChassisSpeedLevelCustom
	ChassisSpeedLevelTypeCount
)
