package key

import (
	"fmt"
	"reflect"

	"github.com/brunoga/robomaster/unitybridge/unity/event"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

// Key represents a Unity Bridge event key which are actions that can be
// performed and attributes that can be read or written. This is used to control
// the robot through the bridge.
//
// All the known keys are listed in the var section below.
type Key struct {
	name        string
	subType     uint32
	accessType  AccessType
	resultValue any
}

const numKeys = 343

func init() {
	// Simple key validation. Make sure we have the right number of them.
	if len(keyBySubType) != numKeys {
		panic(fmt.Sprintf("Unexpected number of keys: %d (wanted %d)",
			len(keyBySubType), numKeys))
	}
}

var (
	keyBySubType = make(map[uint32]*Key, numKeys)

	KeyProductTest = newKey("KeyProductTest", 1, AccessTypeWrite, nil)
	KeyProductType = newKey("KeyProductType", 2, AccessTypeRead, nil)

	KeyCameraConnection                    = newKey("KeyCameraConnection", 16777217, AccessTypeRead, &value.Bool{})
	KeyCameraFirmwareVersion               = newKey("KeyCameraFirmwareVersion", 16777218, AccessTypeRead, nil)
	KeyCameraStartShootPhoto               = newKey("KeyCameraStartShootPhoto", 16777219, AccessTypeAction, nil)
	KeyCameraIsShootingPhoto               = newKey("KeyCameraIsShootingPhoto", 16777220, AccessTypeRead, nil)
	KeyCameraPhotoSize                     = newKey("KeyCameraPhotoSize", 16777221, AccessTypeRead|AccessTypeWrite, nil)
	KeyCameraStartRecordVideo              = newKey("KeyCameraStartRecordVideo", 16777222, AccessTypeAction, &value.Void{})
	KeyCameraStopRecordVideo               = newKey("KeyCameraStopRecordVideo", 16777223, AccessTypeAction, &value.Void{})
	KeyCameraIsRecording                   = newKey("KeyCameraIsRecording", 16777224, AccessTypeRead, &value.Bool{})
	KeyCameraCurrentRecordingTimeInSeconds = newKey("KeyCameraCurrentRecordingTimeInSeconds", 16777225, AccessTypeRead, &value.Uint64{})
	KeyCameraVideoFormat                   = newKey("KeyCameraVideoFormat", 16777226, AccessTypeRead|AccessTypeWrite, &value.Uint64{})
	KeyCameraMode                          = newKey("KeyCameraMode", 16777227, AccessTypeRead|AccessTypeWrite, &value.Uint64{})
	KeyCameraDigitalZoomFactor             = newKey("KeyCameraDigitalZoomFactor", 16777228, AccessTypeRead|AccessTypeWrite, nil)
	KeyCameraAntiFlicker                   = newKey("KeyCameraAntiFlicker", 16777229, AccessTypeRead|AccessTypeWrite, nil)
	KeyCameraSwitch                        = newKey("KeyCameraSwitch", 16777230, AccessTypeAction, nil)
	KeyCameraCurrentCameraIndex            = newKey("KeyCameraCurrentCameraIndex", 16777231, AccessTypeRead, nil)
	KeyCameraHasMainCamera                 = newKey("KeyCameraHasMainCamera", 16777232, AccessTypeRead, nil)
	KeyCameraHasSecondaryCamera            = newKey("KeyCameraHasSecondaryCamera", 16777233, AccessTypeRead, nil)
	KeyCameraIsTimeSynced                  = newKey("KeyCameraIsTimeSynced", 16777243, AccessTypeRead, nil)
	KeyCameraDate                          = newKey("KeyCameraDate", 16777244, AccessTypeRead|AccessTypeWrite, nil)
	KeyCameraVideoTransRate                = newKey("KeyCameraVideoTransRate", 16777245, AccessTypeWrite, &value.Float64{})
	KeyCameraRequestIFrame                 = newKey("KeyCameraRequestIFrame", 16777246, AccessTypeAction, nil)
	KeyCameraAntiLarsenAlgorithmEnable     = newKey("KeyCameraAntiLarsenAlgorithmEnable", 16777247, AccessTypeWrite, nil)

	KeyCameraFormatSDCard                          = newKey("KeyCameraFormatSDCard", 16777234, AccessTypeAction, &value.Bool{})
	KeyCameraSDCardIsFormatting                    = newKey("KeyCameraSDCardIsFormatting", 16777235, AccessTypeRead, &value.Bool{})
	KeyCameraSDCardIsFull                          = newKey("KeyCameraSDCardIsFull", 16777236, AccessTypeRead, &value.Bool{})
	KeyCameraSDCardHasError                        = newKey("KeyCameraSDCardHasError", 16777237, AccessTypeRead, &value.Bool{})
	KeyCameraSDCardIsInserted                      = newKey("KeyCameraSDCardIsInserted", 16777238, AccessTypeRead, &value.Bool{})
	KeyCameraSDCardTotalSpaceInMB                  = newKey("KeyCameraSDCardTotalSpaceInMB", 16777239, AccessTypeRead, &value.Uint64{})
	KeyCameraSDCardRemainingSpaceInMB              = newKey("KeyCameraSDCardRemaingSpaceInMB", 16777240, AccessTypeRead, &value.Uint64{})
	KeyCameraSDCardAvailablePhotoCount             = newKey("KeyCameraSDCardAvailablePhotoCount", 16777241, AccessTypeRead, &value.Uint64{})
	KeyCameraSDCardAvailableRecordingTimeInSeconds = newKey("KeyCameraSDCardAvailableRecordingTimeInSeconds", 16777242, AccessTypeRead, &value.Uint64{})

	KeyMainControllerConnection             = newKey("KeyMainControllerConnection", 33554433, AccessTypeRead, &value.Bool{})
	KeyMainControllerFirmwareVersion        = newKey("KeyMainControllerFirmwareVersion", 33554434, AccessTypeRead, nil)
	KeyMainControllerLoaderVersion          = newKey("KeyMainControllerLoaderVersion", 33554435, AccessTypeRead, nil)
	KeyMainControllerVirtualStick           = newKey("KeyMainControllerVirtualStick", 33554436, AccessTypeAction, nil)
	KeyMainControllerVirtualStickEnabled    = newKey("KeyMainControllerVirtualStickEnabled", 33554437, AccessTypeRead|AccessTypeWrite, &value.Uint64{}) // broken
	KeyMainControllerChassisSpeedMode       = newKey("KeyMainControllerChassisSpeedMode", 33554438, AccessTypeWrite, &value.Uint64{})
	KeyMainControllerChassisFollowMode      = newKey("KeyMainControllerChassisFollowMode", 33554439, AccessTypeWrite, &value.Uint64{})
	KeyMainControllerChassisCarControlMode  = newKey("KeyMainControllerChassisCarControlMode", 33554440, AccessTypeWrite, &value.Uint64{})
	KeyMainControllerRecordState            = newKey("KeyMainControllerRecordState", 33554441, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerGetRecordSetting       = newKey("KeyMainControllerGetRecordSetting", 33554442, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSetRecordSetting       = newKey("KeyMainControllerSetRecordSetting", 33554443, AccessTypeRead, nil)
	KeyMainControllerPlayRecordAttr         = newKey("KeyMainControllerPlayRecordAttr", 33554444, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerGetPlayRecordSetting   = newKey("KeyMainControllerGetPlayRecordSetting", 33554445, AccessTypeRead, nil)
	KeyMainControllerSetPlayRecordSetting   = newKey("KeyMainControllerSetPlayRecordSetting", 33554446, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerMaxSpeedForward        = newKey("KeyMainControllerMaxSpeedForward", 33554447, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerMaxSpeedBackward       = newKey("KeyMainControllerMaxSpeedBackward", 33554448, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerMaxSpeedLateral        = newKey("KeyMainControllerMaxSpeedLateral", 33554449, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSlopeY                 = newKey("KeyMainControllerSlopeY", 33554450, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSlopeX                 = newKey("KeyMainControllerSlopeX", 33554451, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSlopeBreakY            = newKey("KeyMainControllerSlopeBreakY", 33554452, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSlopeBreakX            = newKey("KeyMainControllerSlopeBreakX", 33554453, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerMaxSpeedForwardConfig  = newKey("KeyMainControllerMaxSpeedForwardConfig", 33554454, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerMaxSpeedBackwardConfig = newKey("KeyMainControllerMaxSpeedBackwardConfig", 33554455, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerMaxSpeedLateralConfig  = newKey("KeyMainControllerMaxSpeedLateralConfig", 33554456, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSlopSpeedYConfig       = newKey("KeyMainControllerSlopSpeedYConfig", 33554457, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSlopSpeedXConfig       = newKey("KeyMainControllerSlopSpeedXConfig", 33554458, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSlopBreakYConfig       = newKey("KeyMainControllerSlopBreakYConfig", 33554459, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerSlopBreakXConfig       = newKey("KeyMainControllerSlopBreakXConfig", 33554460, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerChassisPosition        = newKey("KeyMainControllerChassisPosition", 33554461, AccessTypeAction, &value.ChassisPosition{})
	KeyMainControllerWheelSpeed             = newKey("KeyMainControllerWheelSpeed", 33554462, AccessTypeWrite, nil)
	KeyMainControllerArmServoID             = newKey("KeyMainControllerArmServoID", 33554477, AccessTypeRead|AccessTypeWrite, nil)
	KeyMainControllerServoAddressing        = newKey("KeyMainControllerServoAddressing", 33554478, AccessTypeAction, nil)
	KeyMainControllerGetLinkAck             = newKey("KeyMainControllerGetLinkAck", 83886091, AccessTypeRead, nil)

	KeyRobomasterMainControllerEscEncodingStatus        = newKey("KeyRobomasterMainControllerEscEncodingStatus", 33554463, AccessTypeRead, nil)
	KeyRobomasterMainControllerEscEncodeFlag            = newKey("KeyRobomasterMainControllerEscEncodeFlag", 33554464, AccessTypeWrite, nil)
	KeyRobomasterMainControllerStartIMUCalibration      = newKey("KeyRobomasterMainControllerStartIMUCalibration", 33554465, AccessTypeAction, nil)
	KeyRobomasterMainControllerIMUCalibrationState      = newKey("KeyRobomasterMainControllerIMUCalibrationState", 33554466, AccessTypeRead, nil)
	KeyRobomasterMainControllerIMUCalibrationCurrSide   = newKey("KeyRobomasterMainControllerIMUCalibrationCurrSide", 33554467, AccessTypeRead, nil)
	KeyRobomasterMainControllerIMUCalibrationProgress   = newKey("KeyRobomasterMainControllerIMUCalibrationProgress", 33554468, AccessTypeRead, nil)
	KeyRobomasterMainControllerIMUCalibrationFailCode   = newKey("KeyRobomasterMainControllerIMUCalibrationFailCode", 33554469, AccessTypeRead, nil)
	KeyRobomasterMainControllerIMUCalibrationFinishFlag = newKey("KeyRobomasterMainControllerIMUCalibrationFinishFlag", 33554470, AccessTypeRead, nil)
	KeyRobomasterMainControllerStopIMUCalibration       = newKey("KeyRobomasterMainControllerStopIMUCalibration", 33554471, AccessTypeAction, nil)
	KeyRobomasterMainControllerRelativePosition         = newKey("KeyRobomasterMainControllerRelativePosition", 33554476, AccessTypeRead, nil)

	KeyRobomasterChassisMode              = newKey("KeyRobomasterChassisMode", 33554472, AccessTypeRead, nil)
	KeyRobomasterChassisSpeed             = newKey("KeyRobomasterChassisSpeed", 33554473, AccessTypeRead, nil)
	KeyRobomasterOpenChassisSpeedUpdates  = newKey("KeyRobomasterOpenChassisSpeedUpdates", 33554474, AccessTypeAction, nil)
	KeyRobomasterCloseChassisSpeedUpdates = newKey("KeyRobomasterCloseChassisSpeedUpdates", 33554475, AccessTypeAction, nil)

	KeyRobomasterSystemConnection                       = newKey("KeyRobomasterSystemConnection", 83886081, AccessTypeRead, &value.Bool{})
	KeyRobomasterSystemFirmwareVersion                  = newKey("KeyRobomasterSystemFirmwareVersion", 83886082, AccessTypeRead, nil)
	KeyRobomasterSystemCANFirmwareVersion               = newKey("KeyRobomasterSystemCANFirmwareVersion", 83886083, AccessTypeRead, nil)
	KeyRobomasterSystemScratchFirmwareVersion           = newKey("KeyRobomasterSystemScratchFirmwareVersion", 83886084, AccessTypeRead, nil)
	KeyRobomasterSystemSerialNumber                     = newKey("KeyRobomasterSystemSerialNumber", 83886085, AccessTypeRead, nil)
	KeyRobomasterSystemAbilitiesAttack                  = newKey("KeyRobomasterSystemAbilitiesAttack", 83886086, AccessTypeAction, nil)
	KeyRobomasterSystemUnderAbilitiesAttack             = newKey("KeyRobomasterSystemUnderAbilitiesAttack", 83886087, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemKill                             = newKey("KeyRobomasterSystemKill", 83886088, AccessTypeAction, nil)
	KeyRobomasterSystemRevive                           = newKey("KeyRobomasterSystemRevive", 83886089, AccessTypeAction, nil)
	KeyRobomasterSystemGet1860LinkAck                   = newKey("KeyRobomasterSystemGet1860LinkAck", 83886090, AccessTypeRead, nil)
	KeyRobomasterSystemGameRoleConfig                   = newKey("KeyRobomasterSystemGameRoleConfig", 83886093, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemGameColorConfig                  = newKey("KeyRobomasterSystemGameColorConfig", 83886094, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemGameStart                        = newKey("KeyRobomasterSystemGameStart", 83886095, AccessTypeAction, nil)
	KeyRobomasterSystemGameEnd                          = newKey("KeyRobomasterSystemGameEnd", 83886096, AccessTypeAction, nil)
	KeyRobomasterSystemDebugLog                         = newKey("KeyRobomasterSystemDebugLog", 83886097, AccessTypeRead, nil)
	KeyRobomasterSystemSoundEnabled                     = newKey("KeyRobomasterSystemSoundEnabled", 83886098, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemLeftHeadlightBrightness          = newKey("KeyRobomasterSystemLeftHeadlightBrightness", 83886099, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemRightHeadlightBrightness         = newKey("KeyRobomasterSystemRightHeadlightBrightness", 83886100, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemLEDColor                         = newKey("KeyRobomasterSystemLEDColor", 83886101, AccessTypeWrite, nil)
	KeyRobomasterSystemUploadScratch                    = newKey("KeyRobomasterSystemUploadScratch", 83886102, AccessTypeWrite, nil)
	KeyRobomasterSystemUploadScratchByFTP               = newKey("KeyRobomasterSystemUploadScratchByFTP", 83886103, AccessTypeWrite, nil)
	KeyRobomasterSystemUninstallScratchSkill            = newKey("KeyRobomasterSystemUninstallScratchSkill", 83886104, AccessTypeAction, nil)
	KeyRobomasterSystemInstallScratchSkill              = newKey("KeyRobomasterSystemInstallScratchSkill", 83886105, AccessTypeAction, nil)
	KeyRobomasterSystemInquiryDspMd5                    = newKey("KeyRobomasterSystemInquiryDspMd5", 83886106, AccessTypeWrite, nil)
	KeyRobomasterSystemInquiryDspMd5Ack                 = newKey("KeyRobomasterSystemInquiryDspMd5Ack", 83886107, AccessTypeWrite, nil)
	KeyRobomasterSystemInquiryDspResourceMd5            = newKey("KeyRobomasterSystemInquiryDspResourceMd5", 83886108, AccessTypeWrite, nil)
	KeyRobomasterSystemInquiryDspResourceMd5Ack         = newKey("KeyRobomasterSystemInquiryDspResourceMd5Ack", 83886109, AccessTypeWrite, nil)
	KeyRobomasterSystemLaunchSinglePlayerCustomSkill    = newKey("KeyRobomasterSystemLaunchSinglePlayerCustomSkill", 83886110, AccessTypeAction, nil)
	KeyRobomasterSystemStopSinglePlayerCustomSkill      = newKey("KeyRobomasterSystemStopSinglePlayerCustomSkill", 83886111, AccessTypeAction, nil)
	KeyRobomasterSystemControlScratch                   = newKey("KeyRobomasterSystemControlScratch", 83886112, AccessTypeAction, nil)
	KeyRobomasterSystemScratchState                     = newKey("KeyRobomasterSystemScratchState", 83886113, AccessTypeRead, nil)
	KeyRobomasterSystemScratchCallback                  = newKey("KeyRobomasterSystemScratchCallback", 83886114, AccessTypeRead, nil)
	KeyRobomasterSystemForesightPosition                = newKey("KeyRobomasterSystemForesightPosition", 83886115, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemPullLogFiles                     = newKey("KeyRobomasterSystemPullLogFiles", 83886116, AccessTypeRead, nil)
	KeyRobomasterSystemCurrentHP                        = newKey("KeyRobomasterSystemCurrentHP", 83886117, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemTotalHP                          = newKey("KeyRobomasterSystemTotalHP", 83886118, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemCurrentBullets                   = newKey("KeyRobomasterSystemCurrentBullets", 83886119, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemTotalBullets                     = newKey("KeyRobomasterSystemTotalBullets", 83886120, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemEquipments                       = newKey("KeyRobomasterSystemEquipments", 83886121, AccessTypeRead, nil)
	KeyRobomasterSystemBuffs                            = newKey("KeyRobomasterSystemBuffs", 83886122, AccessTypeRead, nil)
	KeyRobomasterSystemSkillStatus                      = newKey("KeyRobomasterSystemSkillStatus", 83886123, AccessTypeRead, nil)
	KeyRobomasterSystemGunCoolDown                      = newKey("KeyRobomasterSystemGunCoolDown", 83886124, AccessTypeRead, nil)
	KeyRobomasterSystemGameConfigList                   = newKey("KeyRobomasterSystemGameConfigList", 83886125, AccessTypeWrite, nil)
	KeyRobomasterSystemCarAndSkillID                    = newKey("KeyRobomasterSystemCarAndSkillID", 83886126, AccessTypeWrite, nil)
	KeyRobomasterSystemAppStatus                        = newKey("KeyRobomasterSystemAppStatus", 83886127, AccessTypeWrite, nil)
	KeyRobomasterSystemLaunchMultiPlayerSkill           = newKey("KeyRobomasterSystemLaunchMultiPlayerSkill", 83886128, AccessTypeAction, nil)
	KeyRobomasterSystemStopMultiPlayerSkill             = newKey("KeyRobomasterSystemStopMultiPlayerSkill", 83886129, AccessTypeAction, nil)
	KeyRobomasterSystemConfigSkillTable                 = newKey("KeyRobomasterSystemConfigSkillTable", 83886130, AccessTypeWrite, nil)
	KeyRobomasterSystemWorkingDevices                   = newKey("KeyRobomasterSystemWorkingDevices", 83886131, AccessTypeRead, &value.List[uint16]{})
	KeyRobomasterSystemExceptions                       = newKey("KeyRobomasterSystemExceptions", 83886132, AccessTypeRead, nil)
	KeyRobomasterSystemTaskStatus                       = newKey("KeyRobomasterSystemTaskStatus", 83886133, AccessTypeRead, &value.TaskStatus{})
	KeyRobomasterSystemReturnEnabled                    = newKey("KeyRobomasterSystemReturnEnabled", 83886134, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemSafeMode                         = newKey("KeyRobomasterSystemSafeMode", 83886135, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemScratchExecuteState              = newKey("KeyRobomasterSystemScratchExecuteState", 83886136, AccessTypeRead, nil)
	KeyRobomasterSystemAttitudeInfo                     = newKey("KeyRobomasterSystemAttitudeInfo", 83886137, AccessTypeRead, nil)
	KeyRobomasterSystemSightBeadPosition                = newKey("KeyRobomasterSystemSightBeadPosition", 83886138, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemSpeakerLanguage                  = newKey("KeyRobomasterSystemSpeakerLanguage", 83886139, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemSpeakerVolumn                    = newKey("KeyRobomasterSystemSpeakerVolumn", 83886140, AccessTypeRead|AccessTypeWrite, &value.Uint64{})
	KeyRobomasterSystemChassisSpeedLevel                = newKey("KeyRobomasterSystemChassisSpeedLevel", 83886141, AccessTypeRead|AccessTypeWrite, &value.Uint64{})
	KeyRobomasterSystemIsEncryptedFirmware              = newKey("KeyRobomasterSystemIsEncryptedFirmware", 83886142, AccessTypeRead, nil)
	KeyRobomasterSystemScratchErrorInfo                 = newKey("KeyRobomasterSystemScratchErrorInfo", 83886143, AccessTypeRead, nil)
	KeyRobomasterSystemScratchOutputInfo                = newKey("KeyRobomasterSystemScratchOutputInfo", 83886144, AccessTypeRead, nil)
	KeyRobomasterSystemBarrelCoolDown                   = newKey("KeyRobomasterSystemBarrelCoolDown", 83886145, AccessTypeAction, nil)
	KeyRobomasterSystemResetBarrelOverheat              = newKey("KeyRobomasterSystemResetBarrelOverheat", 83886146, AccessTypeAction, nil)
	KeyRobomasterSystemMobileAccelerInfo                = newKey("KeyRobomasterSystemMobileAccelerInfo", 83886147, AccessTypeWrite, nil)
	KeyRobomasterSystemMobileGyroAttitudeAngleInfo      = newKey("KeyRobomasterSystemMobileGyroAttitudeAngleInfo", 83886148, AccessTypeWrite, nil)
	KeyRobomasterSystemMobileGyroRotationRateInfo       = newKey("KeyRobomasterSystemMobileGyroRotationRateInfo", 83886149, AccessTypeWrite, nil)
	KeyRobomasterSystemEnableAcceleratorSubscribe       = newKey("KeyRobomasterSystemEnableAcceleratorSubscribe", 83886150, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemEnableGyroRotationRateSubscribe  = newKey("KeyRobomasterSystemEnableGyroRotationRateSubscribe", 83886151, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemEnableGyroAttitudeAngleSubscribe = newKey("KeyRobomasterSystemEnableGyroAttitudeAngleSubscribe", 83886152, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemDeactivate                       = newKey("KeyRobomasterSystemDeactivate", 83886153, AccessTypeAction, nil)
	KeyRobomasterSystemFunctionEnable                   = newKey("KeyRobomasterSystemFunctionEnable", 83886154, AccessTypeAction, &value.FunctionEnable{})
	KeyRobomasterSystemIsGameRunning                    = newKey("KeyRobomasterSystemIsGameRunning", 83886155, AccessTypeRead, nil)
	KeyRobomasterSystemIsActivated                      = newKey("KeyRobomasterSystemIsActivated", 83886156, AccessTypeRead, nil)
	KeyRobomasterSystemLowPowerConsumption              = newKey("KeyRobomasterSystemLowPowerConsumption", 83886157, AccessTypeRead|AccessTypeWrite, nil)
	KeyRobomasterSystemEnterLowPowerConsumption         = newKey("KeyRobomasterSystemEnterLowPowerConsumption", 83886158, AccessTypeAction, nil)
	KeyRobomasterSystemExitLowPowerConsumption          = newKey("KeyRobomasterSystemExitLowPowerConsumption", 83886159, AccessTypeAction, nil)
	KeyRobomasterSystemIsLowPowerConsumption            = newKey("KeyRobomasterSystemIsLowPowerConsumption", 83886160, AccessTypeRead, nil)
	KeyRobomasterSystemPushFile                         = newKey("KeyRobomasterSystemPushFile", 83886161, AccessTypeAction, nil)
	KeyRobomasterSystemPlaySound                        = newKey("KeyRobomasterSystemPlaySound", 83886162, AccessTypeAction, nil)
	KeyRobomasterSystemPlaySoundStatus                  = newKey("KeyRobomasterSystemPlaySoundStatus", 83886163, AccessTypeRead, nil)
	KeyRobomasterSystemCustomUIAttribute                = newKey("KeyRobomasterSystemCustomUIAttribute", 83886164, AccessTypeRead, nil)
	KeyRobomasterSystemCustomUIFunctionEvent            = newKey("KeyRobomasterSystemCustomUIFunctionEvent", 83886165, AccessTypeAction, nil)
	KeyRobomasterSystemTotalMileage                     = newKey("KeyRobomasterSystemTotalMileage", 83886166, AccessTypeRead, nil)
	KeyRobomasterSystemTotalDrivingTime                 = newKey("KeyRobomasterSystemTotalDrivingTime", 83886167, AccessTypeRead, nil)
	KeyRobomasterSystemSetPlayMode                      = newKey("KeyRobomasterSystemSetPlayMode", 83886168, AccessTypeWrite, nil)
	KeyRobomasterSystemCustomSkillInfo                  = newKey("KeyRobomasterSystemCustomSkillInfo", 83886169, AccessTypeRead, nil)
	KeyRobomasterSystemAddressing                       = newKey("KeyRobomasterSystemAddressing", 83886170, AccessTypeAction, nil)
	KeyRobomasterSystemLEDLightEffect                   = newKey("KeyRobomasterSystemLEDLightEffect", 83886171, AccessTypeAction, nil)
	KeyRobomasterSystemOpenImageTransmission            = newKey("KeyRobomasterSystemOpenImageTransmission", 83886172, AccessTypeAction, nil)
	KeyRobomasterSystemCloseImageTransmission           = newKey("KeyRobomasterSystemCloseImageTransmission", 83886173, AccessTypeAction, nil)

	KeyRobomasterWaterGunFirmwareVersion       = newKey("KeyRobomasterWaterGunFirmwareVersion", 167772161, AccessTypeRead, nil)
	KeyRobomasterWaterGunWaterGunFire          = newKey("KeyRobomasterWaterGunWaterGunFire", 167772162, AccessTypeAction, nil)
	KeyRobomasterWaterGunWaterGunFireWithTimes = newKey("KeyRobomasterWaterGunWaterGunFireWithTimes", 167772163, AccessTypeAction, nil)
	KeyRobomasterWaterGunShootSpeed            = newKey("KeyRobomasterWaterGunShootSpeed", 167772164, AccessTypeRead, nil)
	KeyRobomasterWaterGunShootFrequency        = newKey("KeyRobomasterWaterGunShootFrequency", 167772165, AccessTypeRead, nil)

	KeyRobomasterInfraredGunConnection      = newKey("KeyRobomasterInfraredGunConnection", 301989889, AccessTypeRead, nil)
	KeyRobomasterInfraredGunFirmwareVersion = newKey("KeyRobomasterInfraredGunFirmwareVersion", 301989890, AccessTypeRead, nil)
	KeyRobomasterInfraredGunInfraredGunFire = newKey("KeyRobomasterInfraredGunInfraredGunFire", 301989891, AccessTypeAction, nil)
	KeyRobomasterInfraredGunShootFrequency  = newKey("KeyRobomasterInfraredGunShootFrequency", 301989892, AccessTypeRead, nil)

	KeyRobomasterBatteryFirmwareVersion = newKey("KeyRobomasterBatteryFirmwareVersion", 218103809, AccessTypeRead, nil)
	KeyRobomasterBatteryPowerPercent    = newKey("KeyRobomasterBatteryPowerPercent", 218103810, AccessTypeRead, &value.Uint64{})
	KeyRobomasterBatteryVoltage         = newKey("KeyRobomasterBatteryVoltage", 218103811, AccessTypeRead, nil)
	KeyRobomasterBatteryTemperature     = newKey("KeyRobomasterBatteryTemperature", 218103812, AccessTypeRead, nil)
	KeyRobomasterBatteryCurrent         = newKey("KeyRobomasterBatteryCurrent", 218103813, AccessTypeRead, nil)
	KeyRobomasterBatteryShutdown        = newKey("KeyRobomasterBatteryShutdown", 218103814, AccessTypeAction, nil)
	KeyRobomasterBatteryReboot          = newKey("KeyRobomasterBatteryReboot", 218103815, AccessTypeAction, nil)

	KeyRobomasterGamePadConnection                   = newKey("KeyRobomasterGamePadConnection", 234881025, AccessTypeRead, &value.Bool{})
	KeyRobomasterGamePadFirmwareVersion              = newKey("KeyRobomasterGamePadFirmwareVersion", 234881026, AccessTypeRead, &value.String{})
	KeyRobomasterGamePadHasMouse                     = newKey("KeyRobomasterGamePadHasMouse", 234881027, AccessTypeRead, nil)
	KeyRobomasterGamePadHasKeyboard                  = newKey("KeyRobomasterGamePadHasKeyboard", 234881028, AccessTypeRead, nil)
	KeyRobomasterGamePadCtrlSensitivityX             = newKey("KeyRobomasterGamePadCtrlSensitivityX", 234881029, AccessTypeWrite, nil)
	KeyRobomasterGamePadCtrlSensitivityY             = newKey("KeyRobomasterGamePadCtrlSensitivityY", 234881030, AccessTypeWrite, nil)
	KeyRobomasterGamePadCtrlSensitivityYaw           = newKey("KeyRobomasterGamePadCtrlSensitivityYaw", 234881031, AccessTypeWrite, nil)
	KeyRobomasterGamePadCtrlSensitivityYawSlop       = newKey("KeyRobomasterGamePadCtrlSensitivityYawSlop", 234881032, AccessTypeWrite, nil)
	KeyRobomasterGamePadCtrlSensitivityYawDeadZone   = newKey("KeyRobomasterGamePadCtrlSensitivityYawDeadZone", 234881033, AccessTypeWrite, nil)
	KeyRobomasterGamePadCtrlSensitivityPitch         = newKey("KeyRobomasterGamePadCtrlSensitivityPitch", 234881034, AccessTypeWrite, nil)
	KeyRobomasterGamePadCtrlSensitivityPitchSlop     = newKey("KeyRobomasterGamePadCtrlSensitivityPitchSlop", 234881035, AccessTypeWrite, nil)
	KeyRobomasterGamePadCtrlSensitivityPitchDeadZone = newKey("KeyRobomasterGamePadCtrlSensitivityPitchDeadZone", 234881036, AccessTypeWrite, nil)
	KeyRobomasterGamePadMouseLeftButton              = newKey("KeyRobomasterGamePadMouseLeftButton", 234881037, AccessTypeRead, nil)
	KeyRobomasterGamePadMouseRightButton             = newKey("KeyRobomasterGamePadMouseRightButton", 234881038, AccessTypeRead, nil)
	KeyRobomasterGamePadC1                           = newKey("KeyRobomasterGamePadC1", 234881039, AccessTypeRead, nil)
	KeyRobomasterGamePadC2                           = newKey("KeyRobomasterGamePadC2", 234881040, AccessTypeRead, nil)
	KeyRobomasterGamePadFire                         = newKey("KeyRobomasterGamePadFire", 234881041, AccessTypeRead, nil)
	KeyRobomasterGamePadFn                           = newKey("KeyRobomasterGamePadFn", 234881042, AccessTypeRead, nil)
	KeyRobomasterGamePadNoCalibrate                  = newKey("KeyRobomasterGamePadNoCalibrate", 234881043, AccessTypeRead, nil)
	KeyRobomasterGamePadNotAtMiddle                  = newKey("KeyRobomasterGamePadNotAtMiddle", 234881044, AccessTypeRead, nil)
	KeyRobomasterGamePadBatteryWarning               = newKey("KeyRobomasterGamePadBatteryWarning", 234881045, AccessTypeRead, nil)
	KeyRobomasterGamePadBatteryPercent               = newKey("KeyRobomasterGamePadBatteryPercent", 234881046, AccessTypeRead, nil)
	KeyRobomasterGamePadActivationSettings           = newKey("KeyRobomasterGamePadActivationSettings", 234881047, AccessTypeRead|AccessTypeWrite, &value.GamePadActivationSettings{})
	KeyRobomasterGamePadControlEnabled               = newKey("KeyRobomasterGamePadControlEnabled", 234881048, AccessTypeWrite, &value.Bool{})

	KeyRobomasterClawConnection          = newKey("KeyRobomasterClawConnection", 251658241, AccessTypeRead, nil)
	KeyRobomasterClawFirmwareVersion     = newKey("KeyRobomasterClawFirmwareVersion", 251658242, AccessTypeRead, nil)
	KeyRobomasterClawCtrl                = newKey("KeyRobomasterClawCtrl", 251658243, AccessTypeAction, nil)
	KeyRobomasterClawStatus              = newKey("KeyRobomasterClawStatus", 251658244, AccessTypeRead, nil)
	KeyRobomasterClawInfoSubscribe       = newKey("KeyRobomasterClawInfoSubscribe", 251658245, AccessTypeRead, nil)
	KeyRobomasterEnableClawInfoSubscribe = newKey("KeyRobomasterEnableClawInfoSubscribe", 251658246, AccessTypeAction, nil)

	KeyRobomasterArmConnection          = newKey("KeyRobomasterArmConnection", 285212673, AccessTypeRead, nil)
	KeyRobomasterArmCtrl                = newKey("KeyRobomasterArmCtrl", 285212674, AccessTypeAction, nil)
	KeyRobomasterArmCtrlMode            = newKey("KeyRobomasterArmCtrlMode", 285212675, AccessTypeAction, nil)
	KeyRobomasterArmCalibration         = newKey("KeyRobomasterArmCalibration", 285212676, AccessTypeAction, nil)
	KeyRobomasterArmBlockedFlag         = newKey("KeyRobomasterArmBlockedFlag", 285212677, AccessTypeRead, nil)
	KeyRobomasterArmPositionSubscribe   = newKey("KeyRobomasterArmPositionSubscribe", 285212678, AccessTypeRead, nil)
	KeyRobomasterArmReachLimitX         = newKey("KeyRobomasterArmReachLimitX", 285212679, AccessTypeRead, nil)
	KeyRobomasterArmReachLimitY         = newKey("KeyRobomasterArmReachLimitY", 285212680, AccessTypeRead, nil)
	KeyRobomasterEnableArmInfoSubscribe = newKey("KeyRobomasterEnableArmInfoSubscribe", 285212681, AccessTypeAction, nil)
	KeyRobomasterArmControlMode         = newKey("KeyRobomasterArmControlMode", 285212682, AccessTypeRead|AccessTypeWrite, nil)

	KeyRobomasterTOFConnection          = newKey("KeyRobomasterTOFConnection", 318767105, AccessTypeRead, nil)
	KeyRobomasterTOFLEDColor            = newKey("KeyRobomasterTOFLEDColor", 318767106, AccessTypeWrite, nil)
	KeyRobomasterTOFOnlineModules       = newKey("KeyRobomasterTOFOnlineModules", 318767107, AccessTypeRead, nil)
	KeyRobomasterTOFInfoSubscribe       = newKey("KeyRobomasterTOFInfoSubscribe", 318767108, AccessTypeRead, nil)
	KeyRobomasterEnableTOFInfoSubscribe = newKey("KeyRobomasterEnableTOFInfoSubscribe", 318767109, AccessTypeAction, nil)
	KeyRobomasterTOFFirmwareVersion1    = newKey("KeyRobomasterTOFFirmwareVersion1", 318767110, AccessTypeRead, nil)
	KeyRobomasterTOFFirmwareVersion2    = newKey("KeyRobomasterTOFFirmwareVersion2", 318767111, AccessTypeRead, nil)
	KeyRobomasterTOFFirmwareVersion3    = newKey("KeyRobomasterTOFFirmwareVersion3", 318767112, AccessTypeRead, nil)
	KeyRobomasterTOFFirmwareVersion4    = newKey("KeyRobomasterTOFFirmwareVersion4", 318767113, AccessTypeRead, nil)

	KeyRobomasterServoConnection          = newKey("KeyRobomasterServoConnection", 335544321, AccessTypeRead, nil)
	KeyRobomasterServoLEDColor            = newKey("KeyRobomasterServoLEDColor", 335544322, AccessTypeWrite, nil)
	KeyRobomasterServoSpeed               = newKey("KeyRobomasterServoSpeed", 335544323, AccessTypeWrite, nil)
	KeyRobomasterServoOnlineModules       = newKey("KeyRobomasterServoOnlineModules", 335544324, AccessTypeRead, nil)
	KeyRobomasterServoInfoSubscribe       = newKey("KeyRobomasterServoInfoSubscribe", 335544325, AccessTypeRead, nil)
	KeyRobomasterEnableServoInfoSubscribe = newKey("KeyRobomasterEnableServoInfoSubscribe", 335544326, AccessTypeAction, nil)
	KeyRobomasterServoFirmwareVersion1    = newKey("KeyRobomasterServoFirmwareVersion1", 335544327, AccessTypeRead, nil)
	KeyRobomasterServoFirmwareVersion2    = newKey("KeyRobomasterServoFirmwareVersion2", 335544328, AccessTypeRead, nil)
	KeyRobomasterServoFirmwareVersion3    = newKey("KeyRobomasterServoFirmwareVersion3", 335544329, AccessTypeRead, nil)
	KeyRobomasterServoFirmwareVersion4    = newKey("KeyRobomasterServoFirmwareVersion4", 335544330, AccessTypeRead, nil)

	KeyRobomasterSensorAdapterConnection          = newKey("KeyRobomasterSensorAdapterConnection", 352321537, AccessTypeRead, nil)
	KeyRobomasterSensorAdapterOnlineModules       = newKey("KeyRobomasterSensorAdapterOnlineModules", 352321538, AccessTypeRead, nil)
	KeyRobomasterSensorAdapterInfoSubscribe       = newKey("KeyRobomasterSensorAdapterInfoSubscribe", 352321539, AccessTypeRead, nil)
	KeyRobomasterEnableSensorAdapterInfoSubscribe = newKey("KeyRobomasterEnableSensorAdapterInfoSubscribe", 352321540, AccessTypeAction, nil)
	KeyRobomasterSensorAdapterFirmwareVersion1    = newKey("KeyRobomasterSensorAdapterFirmwareVersion1", 352321541, AccessTypeRead, nil)
	KeyRobomasterSensorAdapterFirmwareVersion2    = newKey("KeyRobomasterSensorAdapterFirmwareVersion2", 352321542, AccessTypeRead, nil)
	KeyRobomasterSensorAdapterFirmwareVersion3    = newKey("KeyRobomasterSensorAdapterFirmwareVersion3", 352321543, AccessTypeRead, nil)
	KeyRobomasterSensorAdapterFirmwareVersion4    = newKey("KeyRobomasterSensorAdapterFirmwareVersion4", 352321544, AccessTypeRead, nil)
	KeyRobomasterSensorAdapterFirmwareVersion5    = newKey("KeyRobomasterSensorAdapterFirmwareVersion5", 352321545, AccessTypeRead, nil)
	KeyRobomasterSensorAdapterFirmwareVersion6    = newKey("KeyRobomasterSensorAdapterFirmwareVersion6", 352321546, AccessTypeRead, nil)
	KeyRobomasterSensorAdapterLEDColor            = newKey("KeyRobomasterSensorAdapterLEDColor", 352321547, AccessTypeWrite, nil)

	KeyRemoteControllerConnection = newKey("KeyRemoteControllerConnection", 50331649, AccessTypeRead, nil)

	KeyGimbalConnection              = newKey("KeyGimbalConnection", 67108865, AccessTypeRead, &value.Bool{})
	KeyGimbalESCFirmwareVersion      = newKey("KeyGimbalESCFirmwareVersion", 67108866, AccessTypeRead, nil)
	KeyGimbalFirmwareVersion         = newKey("KeyGimbalFirmwareVersion", 67108867, AccessTypeRead, nil)
	KeyGimbalWorkMode                = newKey("KeyGimbalWorkMode", 67108868, AccessTypeRead|AccessTypeWrite, &value.Uint64{})
	KeyGimbalControlMode             = newKey("KeyGimbalControlMode", 67108869, AccessTypeRead|AccessTypeWrite, &value.Uint64{})
	KeyGimbalResetPosition           = newKey("KeyGimbalResetPosition", 67108870, AccessTypeAction, &value.Void{})
	KeyGimbalResetPositionState      = newKey("KeyGimbalResetPositionState", 67108871, AccessTypeRead, &value.Uint64{})
	KeyGimbalCalibration             = newKey("KeyGimbalCalibration", 67108872, AccessTypeAction, nil)
	KeyGimbalSpeedRotation           = newKey("KeyGimbalSpeedRotation", 67108873, AccessTypeAction, &value.GimbalSpeedRotation{})
	KeyGimbalSpeedRotationEnabled    = newKey("KeyGimbalSpeedRotationEnabled", 67108874, AccessTypeWrite|AccessTypeAction /*Added*/, &value.Uint64{})
	KeyGimbalAngleIncrementRotation  = newKey("KeyGimbalAngleIncrementRotation", 67108875, AccessTypeAction, &value.GimbalAngleRotation{})
	KeyGimbalAngleFrontYawRotation   = newKey("KeyGimbalAngleFrontYawRotation", 67108876, AccessTypeAction, &value.GimbalAngleRotation{})
	KeyGimbalAngleFrontPitchRotation = newKey("KeyGimbalAngleFrontPitchRotation", 67108877, AccessTypeAction, &value.GimbalAngleRotation{})
	KeyGimbalAttitude                = newKey("KeyGimbalAttitude", 67108878, AccessTypeRead, &value.GimbalAttitude{})
	KeyGimbalAutoCalibrate           = newKey("KeyGimbalAutoCalibrate", 67108879, AccessTypeAction, nil)
	KeyGimbalCalibrationStatus       = newKey("KeyGimbalCalibrationStatus", 67108880, AccessTypeRead, nil)
	KeyGimbalCalibrationProgress     = newKey("KeyGimbalCalibrationProgress", 67108881, AccessTypeRead, nil)
	KeyGimbalOpenAttitudeUpdates     = newKey("KeyGimbalOpenAttitudeUpdates", 67108882, AccessTypeAction, &value.Void{})
	KeyGimbalCloseAttitudeUpdates    = newKey("KeyGimbalCloseAttitudeUpdates", 67108883, AccessTypeAction, &value.Void{})
	KeyGimbalGetLinkAck              = newKey("KeyGimbalGetLinkAck", 83886092, AccessTypeRead, nil)

	KeyVisionFirmwareVersion             = newKey("KeyVisionFirmwareVersion", 100663297, AccessTypeRead, nil)
	KeyVisionTrackingAutoLockTarget      = newKey("KeyVisionTrackingAutoLockTarget", 100663298, AccessTypeRead|AccessTypeWrite, nil)
	KeyVisionARParameters                = newKey("KeyVisionARParameters", 100663299, AccessTypeRead, nil)
	KeyVisionARTagEnabled                = newKey("KeyVisionARTagEnabled", 100663300, AccessTypeRead, nil)
	KeyVisionDebugRect                   = newKey("KeyVisionDebugRect", 100663301, AccessTypeRead, nil)
	KeyVisionLaserPosition               = newKey("KeyVisionLaserPosition", 100663302, AccessTypeRead, nil)
	KeyVisionDetectionEnable             = newKey("KeyVisionDetectionEnable", 100663303, AccessTypeRead|AccessTypeWrite, nil)
	KeyVisionMarkerRunningStatus         = newKey("KeyVisionMarkerRunningStatus", 100663304, AccessTypeRead, nil)
	KeyVisionTrackingRunningStatus       = newKey("KeyVisionTrackingRunningStatus", 100663305, AccessTypeRead, nil)
	KeyVisionAimbotRunningStatus         = newKey("KeyVisionAimbotRunningStatus", 100663306, AccessTypeRead, nil)
	KeyVisionHeadAndShoulderStatus       = newKey("KeyVisionHeadAndShoulderStatus", 100663307, AccessTypeRead, nil)
	KeyVisionHumanDetectionRunningStatus = newKey("KeyVisionHumanDetectionRunningStatus", 100663308, AccessTypeRead, nil)
	KeyVisionUserConfirm                 = newKey("KeyVisionUserConfirm", 100663309, AccessTypeAction, nil)
	KeyVisionUserCancel                  = newKey("KeyVisionUserCancel", 100663310, AccessTypeAction, nil)
	KeyVisionUserTrackingRect            = newKey("KeyVisionUserTrackingRect", 100663311, AccessTypeWrite, nil)
	KeyVisionTrackingDistance            = newKey("KeyVisionTrackingDistance", 100663312, AccessTypeWrite, nil)
	KeyVisionLineColor                   = newKey("KeyVisionLineColor", 100663313, AccessTypeWrite, nil)
	KeyVisionMarkerColor                 = newKey("KeyVisionMarkerColor", 100663314, AccessTypeWrite, nil)
	KeyVisionMarkerAdvanceStatus         = newKey("KeyVisionMarkerAdvanceStatus", 100663315, AccessTypeRead, nil)

	KeyPerceptionFirmwareVersion = newKey("KeyPerceptionFirmwareVersion", 184549377, AccessTypeRead, nil)
	KeyPerceptionMarkerEnable    = newKey("KeyPerceptionMarkerEnable", 184549378, AccessTypeRead|AccessTypeWrite, nil)
	KeyPerceptionMarkerResult    = newKey("KeyPerceptionMarkerResult", 184549379, AccessTypeRead, nil)

	KeyESCFirmwareVersion1 = newKey("KeyESCFirmwareVersion1", 201326593, AccessTypeRead, nil)
	KeyESCFirmwareVersion2 = newKey("KeyESCFirmwareVersion2", 201326594, AccessTypeRead, nil)
	KeyESCFirmwareVersion3 = newKey("KeyESCFirmwareVersion3", 201326595, AccessTypeRead, nil)
	KeyESCFirmwareVersion4 = newKey("KeyESCFirmwareVersion4", 201326596, AccessTypeRead, nil)
	KeyESCMotorInfomation1 = newKey("KeyESCMotorInfomation1", 201326597, AccessTypeRead, nil)
	KeyESCMotorInfomation2 = newKey("KeyESCMotorInfomation2", 201326598, AccessTypeRead, nil)
	KeyESCMotorInfomation3 = newKey("KeyESCMotorInfomation3", 201326599, AccessTypeRead, nil)
	KeyESCMotorInfomation4 = newKey("KeyESCMotorInfomation4", 201326600, AccessTypeRead, nil)

	KeyWiFiLinkFirmwareVersion         = newKey("KeyWiFiLinkFirmwareVersion", 134217729, AccessTypeRead, nil)
	KeyWiFiLinkDebugInfo               = newKey("KeyWiFiLinkDebugInfo", 134217730, AccessTypeRead, nil)
	KeyWiFiLinkMode                    = newKey("KeyWiFiLinkMode", 134217731, AccessTypeRead, nil)
	KeyWiFiLinkSSID                    = newKey("KeyWiFiLinkSSID", 134217732, AccessTypeRead|AccessTypeWrite, nil)
	KeyWiFiLinkPassword                = newKey("KeyWiFiLinkPassword", 134217733, AccessTypeRead|AccessTypeWrite, nil)
	KeyWiFiLinkAvailableChannelNumbers = newKey("KeyWiFiLinkAvailableChannelNumbers", 134217734, AccessTypeRead, nil)
	KeyWiFiLinkCurrentChannelNumber    = newKey("KeyWiFiLinkCurrentChannelNumber", 134217735, AccessTypeRead|AccessTypeWrite, nil)
	KeyWiFiLinkSNR                     = newKey("KeyWiFiLinkSNR", 134217736, AccessTypeRead, nil)
	KeyWiFiLinkSNRPushEnabled          = newKey("KeyWiFiLinkSNRPushEnabled", 134217737, AccessTypeWrite, nil)
	KeyWiFiLinkReboot                  = newKey("KeyWiFiLinkReboot", 134217738, AccessTypeAction, nil)
	KeyWiFiLinkChannelSelectionMode    = newKey("KeyWiFiLinkChannelSelectionMode", 134217739, AccessTypeRead|AccessTypeWrite, nil)
	KeyWiFiLinkInterference            = newKey("KeyWiFiLinkInterference", 134217740, AccessTypeRead, nil)
	KeyWiFiLinkDeleteNetworkConfig     = newKey("KeyWiFiLinkDeleteNetworkConfig", 134217741, AccessTypeAction, nil)

	KeySDRLinkSNR                  = newKey("KeySDRLinkSNR", 268435457, AccessTypeRead, nil)
	KeySDRLinkBandwidth            = newKey("KeySDRLinkBandwidth", 268435458, AccessTypeRead|AccessTypeWrite, nil)
	KeySDRLinkChannelSelectionMode = newKey("KeySDRLinkChannelSelectionMode", 268435459, AccessTypeRead|AccessTypeWrite, nil)
	KeySDRLinkCurrentFreqPoint     = newKey("KeySDRLinkCurrentFreqPoint", 268435460, AccessTypeRead|AccessTypeWrite, nil)
	KeySDRLinkCurrentFreqBand      = newKey("KeySDRLinkCurrentFreqBand", 268435461, AccessTypeRead|AccessTypeWrite, nil)
	KeySDRLinkIsDualFreqSupported  = newKey("KeySDRLinkIsDualFreqSupported", 268435462, AccessTypeRead, nil)
	KeySDRLinkUpdateConfigs        = newKey("KeySDRLinkUpdateConfigs", 268435463, AccessTypeAction, nil)

	KeyAirLinkConnection         = newKey("KeyAirLinkConnection", 117440513, AccessTypeRead, &value.Bool{})
	KeyAirLinkSignalQuality      = newKey("KeyAirLinkSignalQuality", 117440514, AccessTypeRead, &value.Uint64{})
	KeyAirLinkCountryCode        = newKey("KeyAirLinkCountryCode", 117440515, AccessTypeWrite, nil)
	KeyAirLinkCountryCodeUpdated = newKey("KeyAirLinkCountryCodeUpdated", 117440516, AccessTypeRead, nil)

	KeyArmorFirmwareVersion1 = newKey("KeyArmorFirmwareVersion1", 150994945, AccessTypeRead, nil)
	KeyArmorFirmwareVersion2 = newKey("KeyArmorFirmwareVersion2", 150994946, AccessTypeRead, nil)
	KeyArmorFirmwareVersion3 = newKey("KeyArmorFirmwareVersion3", 150994947, AccessTypeRead, nil)
	KeyArmorFirmwareVersion4 = newKey("KeyArmorFirmwareVersion4", 150994948, AccessTypeRead, nil)
	KeyArmorFirmwareVersion5 = newKey("KeyArmorFirmwareVersion5", 150994949, AccessTypeRead, nil)
	KeyArmorFirmwareVersion6 = newKey("KeyArmorFirmwareVersion6", 150994950, AccessTypeRead, nil)
	KeyArmorUnderAttack      = newKey("KeyArmorUnderAttack", 150994951, AccessTypeRead, nil)
	KeyArmorEnterResetID     = newKey("KeyArmorEnterResetID", 150994952, AccessTypeAction, nil)
	KeyArmorCancelResetID    = newKey("KeyArmorCancelResetID", 150994953, AccessTypeAction, nil)
	KeyArmorSkipCurrentID    = newKey("KeyArmorSkipCurrentID", 150994954, AccessTypeAction, nil)
	KeyArmorResetStatus      = newKey("KeyArmorResetStatus", 150994955, AccessTypeRead, nil)
)

// String returns a string representation of the key.
func (k *Key) String() string {
	if k == nil {
		return "[UNKNOWN KEY]"
	}
	return k.name
}

// SubType returns the sub-type associated wit hhis key. Used for events.
func (k *Key) SubType() uint32 {
	if k == nil {
		return 0
	}
	return k.subType
}

// AccessType returns the access type of the
func (k *Key) AccessType() AccessType {
	if k == nil {
		return AccessTypeNone
	}

	return k.accessType
}

func (k *Key) ResultValue() any {
	if k.resultValue == nil {
		panic(fmt.Sprintf("Unknown result value for key %s.", k.name))
	}

	valueType := reflect.TypeOf(k.resultValue).Elem()

	return reflect.New(valueType).Interface()
}

// FromEvent returns a Key associated with the given event. It returns
// an error in case the key can not be inferred.
func FromEvent(ev *event.Event) (*Key, error) {
	return FromSubType(ev.SubType())
}

// FromSubType returns a Key associated with the given sub-type. It returns
// an error in case the key can not be inferred.
func FromSubType(subType uint32) (*Key, error) {
	k, ok := keyBySubType[subType]
	if !ok {
		return nil, fmt.Errorf("event sub-type does not match any key: %d",
			subType)
	}

	return k, nil
}

func newKey(name string, subType uint32, accessType AccessType,
	resultValue any) *Key {
	if resultValue != nil && reflect.ValueOf(resultValue).Kind() !=
		reflect.Ptr {
		panic("resultValue must be a pointer")
	}

	k := &Key{
		name:        name,
		subType:     subType,
		accessType:  accessType,
		resultValue: resultValue,
	}

	keyBySubType[subType] = k

	return k
}
