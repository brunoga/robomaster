package key

import (
	"fmt"

	"github.com/brunoga/unitybridge/unity/event"
)

// Key represents a Unity Bridge event key which are actions that can be
// performed and attributes that can be read or written. This is used to control
// the robot through the bridge.
//
// All the known keys are listed in the var section below.
type Key struct {
	name       string
	subType    uint32
	accessType AccessType
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

	KeyProductTest = newKey("KeyProductTest", 1, AccessTypeWrite)
	KeyProductType = newKey("KeyProductType", 2, AccessTypeRead)

	KeyCameraConnection                            = newKey("KeyCameraConnection", 16777217, AccessTypeRead)
	KeyCameraFirmwareVersion                       = newKey("KeyCameraFirmwareVersion", 16777218, AccessTypeRead)
	KeyCameraStartShootPhoto                       = newKey("KeyCameraStartShootPhoto", 16777219, AccessTypeAction)
	KeyCameraIsShootingPhoto                       = newKey("KeyCameraIsShootingPhoto", 16777220, AccessTypeRead)
	KeyCameraPhotoSize                             = newKey("KeyCameraPhotoSize", 16777221, AccessTypeRead|AccessTypeWrite)
	KeyCameraStartRecordVideo                      = newKey("KeyCameraStartRecordVideo", 16777222, AccessTypeAction)
	KeyCameraStopRecordVideo                       = newKey("KeyCameraStopRecordVideo", 16777223, AccessTypeAction)
	KeyCameraIsRecording                           = newKey("KeyCameraIsRecording", 16777224, AccessTypeRead)
	KeyCameraCurrentRecordingTimeInSeconds         = newKey("KeyCameraCurrentRecordingTimeInSeconds", 16777225, AccessTypeRead)
	KeyCameraVideoFormat                           = newKey("KeyCameraVideoFormat", 16777226, AccessTypeRead|AccessTypeWrite)
	KeyCameraMode                                  = newKey("KeyCameraMode", 16777227, AccessTypeRead|AccessTypeWrite)
	KeyCameraDigitalZoomFactor                     = newKey("KeyCameraDigitalZoomFactor", 16777228, AccessTypeRead|AccessTypeWrite)
	KeyCameraAntiFlicker                           = newKey("KeyCameraAntiFlicker", 16777229, AccessTypeRead|AccessTypeWrite)
	KeyCameraSwitch                                = newKey("KeyCameraSwitch", 16777230, AccessTypeAction)
	KeyCameraCurrentCameraIndex                    = newKey("KeyCameraCurrentCameraIndex", 16777231, AccessTypeRead)
	KeyCameraHasMainCamera                         = newKey("KeyCameraHasMainCamera", 16777232, AccessTypeRead)
	KeyCameraHasSecondaryCamera                    = newKey("KeyCameraHasSecondaryCamera", 16777233, AccessTypeRead)
	KeyCameraFormatSDCard                          = newKey("KeyCameraFormatSDCard", 16777234, AccessTypeAction)
	KeyCameraSDCardIsFormatting                    = newKey("KeyCameraSDCardIsFormatting", 16777235, AccessTypeRead)
	KeyCameraSDCardIsFull                          = newKey("KeyCameraSDCardIsFull", 16777236, AccessTypeRead)
	KeyCameraSDCardHasError                        = newKey("KeyCameraSDCardHasError", 16777237, AccessTypeRead)
	KeyCameraSDCardIsInserted                      = newKey("KeyCameraSDCardIsInserted", 16777238, AccessTypeRead)
	KeyCameraSDCardTotalSpaceInMB                  = newKey("KeyCameraSDCardTotalSpaceInMB", 16777239, AccessTypeRead)
	KeyCameraSDCardRemaingSpaceInMB                = newKey("KeyCameraSDCardRemaingSpaceInMB", 16777240, AccessTypeRead)
	KeyCameraSDCardAvailablePhotoCount             = newKey("KeyCameraSDCardAvailablePhotoCount", 16777241, AccessTypeRead)
	KeyCameraSDCardAvailableRecordingTimeInSeconds = newKey("KeyCameraSDCardAvailableRecordingTimeInSeconds", 16777242, AccessTypeRead)
	KeyCameraIsTimeSynced                          = newKey("KeyCameraIsTimeSynced", 16777243, AccessTypeRead)
	KeyCameraDate                                  = newKey("KeyCameraDate", 16777244, AccessTypeRead|AccessTypeWrite)
	KeyCameraVideoTransRate                        = newKey("KeyCameraVideoTransRate", 16777245, AccessTypeWrite)
	KeyCameraRequestIFrame                         = newKey("KeyCameraRequestIFrame", 16777246, AccessTypeAction)
	KeyCameraAntiLarsenAlgorithmEnable             = newKey("KeyCameraAntiLarsenAlgorithmEnable", 16777247, AccessTypeWrite)

	KeyMainControllerConnection             = newKey("KeyMainControllerConnection", 33554433, AccessTypeRead)
	KeyMainControllerFirmwareVersion        = newKey("KeyMainControllerFirmwareVersion", 33554434, AccessTypeRead)
	KeyMainControllerLoaderVersion          = newKey("KeyMainControllerLoaderVersion", 33554435, AccessTypeRead)
	KeyMainControllerVirtualStick           = newKey("KeyMainControllerVirtualStick", 33554436, AccessTypeAction)
	KeyMainControllerVirtualStickEnabled    = newKey("KeyMainControllerVirtualStickEnabled", 33554437, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerChassisSpeedMode       = newKey("KeyMainControllerChassisSpeedMode", 33554438, AccessTypeWrite)
	KeyMainControllerChassisFollowMode      = newKey("KeyMainControllerChassisFollowMode", 33554439, AccessTypeWrite)
	KeyMainControllerChassisCarControlMode  = newKey("KeyMainControllerChassisCarControlMode", 33554440, AccessTypeWrite)
	KeyMainControllerRecordState            = newKey("KeyMainControllerRecordState", 33554441, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerGetRecordSetting       = newKey("KeyMainControllerGetRecordSetting", 33554442, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSetRecordSetting       = newKey("KeyMainControllerSetRecordSetting", 33554443, AccessTypeRead)
	KeyMainControllerPlayRecordAttr         = newKey("KeyMainControllerPlayRecordAttr", 33554444, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerGetPlayRecordSetting   = newKey("KeyMainControllerGetPlayRecordSetting", 33554445, AccessTypeRead)
	KeyMainControllerSetPlayRecordSetting   = newKey("KeyMainControllerSetPlayRecordSetting", 33554446, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerMaxSpeedForward        = newKey("KeyMainControllerMaxSpeedForward", 33554447, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerMaxSpeedBackward       = newKey("KeyMainControllerMaxSpeedBackward", 33554448, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerMaxSpeedLateral        = newKey("KeyMainControllerMaxSpeedLateral", 33554449, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSlopeY                 = newKey("KeyMainControllerSlopeY", 33554450, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSlopeX                 = newKey("KeyMainControllerSlopeX", 33554451, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSlopeBreakY            = newKey("KeyMainControllerSlopeBreakY", 33554452, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSlopeBreakX            = newKey("KeyMainControllerSlopeBreakX", 33554453, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerMaxSpeedForwardConfig  = newKey("KeyMainControllerMaxSpeedForwardConfig", 33554454, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerMaxSpeedBackwardConfig = newKey("KeyMainControllerMaxSpeedBackwardConfig", 33554455, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerMaxSpeedLateralConfig  = newKey("KeyMainControllerMaxSpeedLateralConfig", 33554456, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSlopSpeedYConfig       = newKey("KeyMainControllerSlopSpeedYConfig", 33554457, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSlopSpeedXConfig       = newKey("KeyMainControllerSlopSpeedXConfig", 33554458, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSlopBreakYConfig       = newKey("KeyMainControllerSlopBreakYConfig", 33554459, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerSlopBreakXConfig       = newKey("KeyMainControllerSlopBreakXConfig", 33554460, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerChassisPosition        = newKey("KeyMainControllerChassisPosition", 33554461, AccessTypeAction)
	KeyMainControllerWheelSpeed             = newKey("KeyMainControllerWheelSpeed", 33554462, AccessTypeWrite)
	KeyMainControllerArmServoID             = newKey("KeyMainControllerArmServoID", 33554477, AccessTypeRead|AccessTypeWrite)
	KeyMainControllerServoAddressing        = newKey("KeyMainControllerServoAddressing", 33554478, AccessTypeAction)
	KeyMainControllerGetLinkAck             = newKey("KeyMainControllerGetLinkAck", 83886091, AccessTypeRead)

	KeyRobomasterMainControllerEscEncodingStatus        = newKey("KeyRobomasterMainControllerEscEncodingStatus", 33554463, AccessTypeRead)
	KeyRobomasterMainControllerEscEncodeFlag            = newKey("KeyRobomasterMainControllerEscEncodeFlag", 33554464, AccessTypeWrite)
	KeyRobomasterMainControllerStartIMUCalibration      = newKey("KeyRobomasterMainControllerStartIMUCalibration", 33554465, AccessTypeAction)
	KeyRobomasterMainControllerIMUCalibrationState      = newKey("KeyRobomasterMainControllerIMUCalibrationState", 33554466, AccessTypeRead)
	KeyRobomasterMainControllerIMUCalibrationCurrSide   = newKey("KeyRobomasterMainControllerIMUCalibrationCurrSide", 33554467, AccessTypeRead)
	KeyRobomasterMainControllerIMUCalibrationProgress   = newKey("KeyRobomasterMainControllerIMUCalibrationProgress", 33554468, AccessTypeRead)
	KeyRobomasterMainControllerIMUCalibrationFailCode   = newKey("KeyRobomasterMainControllerIMUCalibrationFailCode", 33554469, AccessTypeRead)
	KeyRobomasterMainControllerIMUCalibrationFinishFlag = newKey("KeyRobomasterMainControllerIMUCalibrationFinishFlag", 33554470, AccessTypeRead)
	KeyRobomasterMainControllerStopIMUCalibration       = newKey("KeyRobomasterMainControllerStopIMUCalibration", 33554471, AccessTypeAction)
	KeyRobomasterMainControllerRelativePosition         = newKey("KeyRobomasterMainControllerRelativePosition", 33554476, AccessTypeRead)

	KeyRobomasterChassisMode              = newKey("KeyRobomasterChassisMode", 33554472, AccessTypeRead)
	KeyRobomasterChassisSpeed             = newKey("KeyRobomasterChassisSpeed", 33554473, AccessTypeRead)
	KeyRobomasterOpenChassisSpeedUpdates  = newKey("KeyRobomasterOpenChassisSpeedUpdates", 33554474, AccessTypeAction)
	KeyRobomasterCloseChassisSpeedUpdates = newKey("KeyRobomasterCloseChassisSpeedUpdates", 33554475, AccessTypeAction)

	KeyRobomasterSystemConnection                       = newKey("KeyRobomasterSystemConnection", 83886081, AccessTypeRead)
	KeyRobomasterSystemFirmwareVersion                  = newKey("KeyRobomasterSystemFirmwareVersion", 83886082, AccessTypeRead)
	KeyRobomasterSystemCANFirmwareVersion               = newKey("KeyRobomasterSystemCANFirmwareVersion", 83886083, AccessTypeRead)
	KeyRobomasterSystemScratchFirmwareVersion           = newKey("KeyRobomasterSystemScratchFirmwareVersion", 83886084, AccessTypeRead)
	KeyRobomasterSystemSerialNumber                     = newKey("KeyRobomasterSystemSerialNumber", 83886085, AccessTypeRead)
	KeyRobomasterSystemAbilitiesAttack                  = newKey("KeyRobomasterSystemAbilitiesAttack", 83886086, AccessTypeAction)
	KeyRobomasterSystemUnderAbilitiesAttack             = newKey("KeyRobomasterSystemUnderAbilitiesAttack", 83886087, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemKill                             = newKey("KeyRobomasterSystemKill", 83886088, AccessTypeAction)
	KeyRobomasterSystemRevive                           = newKey("KeyRobomasterSystemRevive", 83886089, AccessTypeAction)
	KeyRobomasterSystemGet1860LinkAck                   = newKey("KeyRobomasterSystemGet1860LinkAck", 83886090, AccessTypeRead)
	KeyRobomasterSystemGameRoleConfig                   = newKey("KeyRobomasterSystemGameRoleConfig", 83886093, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemGameColorConfig                  = newKey("KeyRobomasterSystemGameColorConfig", 83886094, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemGameStart                        = newKey("KeyRobomasterSystemGameStart", 83886095, AccessTypeAction)
	KeyRobomasterSystemGameEnd                          = newKey("KeyRobomasterSystemGameEnd", 83886096, AccessTypeAction)
	KeyRobomasterSystemDebugLog                         = newKey("KeyRobomasterSystemDebugLog", 83886097, AccessTypeRead)
	KeyRobomasterSystemSoundEnabled                     = newKey("KeyRobomasterSystemSoundEnabled", 83886098, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemLeftHeadlightBrightness          = newKey("KeyRobomasterSystemLeftHeadlightBrightness", 83886099, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemRightHeadlightBrightness         = newKey("KeyRobomasterSystemRightHeadlightBrightness", 83886100, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemLEDColor                         = newKey("KeyRobomasterSystemLEDColor", 83886101, AccessTypeWrite)
	KeyRobomasterSystemUploadScratch                    = newKey("KeyRobomasterSystemUploadScratch", 83886102, AccessTypeWrite)
	KeyRobomasterSystemUploadScratchByFTP               = newKey("KeyRobomasterSystemUploadScratchByFTP", 83886103, AccessTypeWrite)
	KeyRobomasterSystemUninstallScratchSkill            = newKey("KeyRobomasterSystemUninstallScratchSkill", 83886104, AccessTypeAction)
	KeyRobomasterSystemInstallScratchSkill              = newKey("KeyRobomasterSystemInstallScratchSkill", 83886105, AccessTypeAction)
	KeyRobomasterSystemInquiryDspMd5                    = newKey("KeyRobomasterSystemInquiryDspMd5", 83886106, AccessTypeWrite)
	KeyRobomasterSystemInquiryDspMd5Ack                 = newKey("KeyRobomasterSystemInquiryDspMd5Ack", 83886107, AccessTypeWrite)
	KeyRobomasterSystemInquiryDspResourceMd5            = newKey("KeyRobomasterSystemInquiryDspResourceMd5", 83886108, AccessTypeWrite)
	KeyRobomasterSystemInquiryDspResourceMd5Ack         = newKey("KeyRobomasterSystemInquiryDspResourceMd5Ack", 83886109, AccessTypeWrite)
	KeyRobomasterSystemLaunchSinglePlayerCustomSkill    = newKey("KeyRobomasterSystemLaunchSinglePlayerCustomSkill", 83886110, AccessTypeAction)
	KeyRobomasterSystemStopSinglePlayerCustomSkill      = newKey("KeyRobomasterSystemStopSinglePlayerCustomSkill", 83886111, AccessTypeAction)
	KeyRobomasterSystemControlScratch                   = newKey("KeyRobomasterSystemControlScratch", 83886112, AccessTypeAction)
	KeyRobomasterSystemScratchState                     = newKey("KeyRobomasterSystemScratchState", 83886113, AccessTypeRead)
	KeyRobomasterSystemScratchCallback                  = newKey("KeyRobomasterSystemScratchCallback", 83886114, AccessTypeRead)
	KeyRobomasterSystemForesightPosition                = newKey("KeyRobomasterSystemForesightPosition", 83886115, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemPullLogFiles                     = newKey("KeyRobomasterSystemPullLogFiles", 83886116, AccessTypeRead)
	KeyRobomasterSystemCurrentHP                        = newKey("KeyRobomasterSystemCurrentHP", 83886117, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemTotalHP                          = newKey("KeyRobomasterSystemTotalHP", 83886118, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemCurrentBullets                   = newKey("KeyRobomasterSystemCurrentBullets", 83886119, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemTotalBullets                     = newKey("KeyRobomasterSystemTotalBullets", 83886120, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemEquipments                       = newKey("KeyRobomasterSystemEquipments", 83886121, AccessTypeRead)
	KeyRobomasterSystemBuffs                            = newKey("KeyRobomasterSystemBuffs", 83886122, AccessTypeRead)
	KeyRobomasterSystemSkillStatus                      = newKey("KeyRobomasterSystemSkillStatus", 83886123, AccessTypeRead)
	KeyRobomasterSystemGunCoolDown                      = newKey("KeyRobomasterSystemGunCoolDown", 83886124, AccessTypeRead)
	KeyRobomasterSystemGameConfigList                   = newKey("KeyRobomasterSystemGameConfigList", 83886125, AccessTypeWrite)
	KeyRobomasterSystemCarAndSkillID                    = newKey("KeyRobomasterSystemCarAndSkillID", 83886126, AccessTypeWrite)
	KeyRobomasterSystemAppStatus                        = newKey("KeyRobomasterSystemAppStatus", 83886127, AccessTypeWrite)
	KeyRobomasterSystemLaunchMultiPlayerSkill           = newKey("KeyRobomasterSystemLaunchMultiPlayerSkill", 83886128, AccessTypeAction)
	KeyRobomasterSystemStopMultiPlayerSkill             = newKey("KeyRobomasterSystemStopMultiPlayerSkill", 83886129, AccessTypeAction)
	KeyRobomasterSystemConfigSkillTable                 = newKey("KeyRobomasterSystemConfigSkillTable", 83886130, AccessTypeWrite)
	KeyRobomasterSystemWorkingDevices                   = newKey("KeyRobomasterSystemWorkingDevices", 83886131, AccessTypeRead)
	KeyRobomasterSystemExceptions                       = newKey("KeyRobomasterSystemExceptions", 83886132, AccessTypeRead)
	KeyRobomasterSystemTaskStatus                       = newKey("KeyRobomasterSystemTaskStatus", 83886133, AccessTypeRead)
	KeyRobomasterSystemReturnEnabled                    = newKey("KeyRobomasterSystemReturnEnabled", 83886134, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemSafeMode                         = newKey("KeyRobomasterSystemSafeMode", 83886135, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemScratchExecuteState              = newKey("KeyRobomasterSystemScratchExecuteState", 83886136, AccessTypeRead)
	KeyRobomasterSystemAttitudeInfo                     = newKey("KeyRobomasterSystemAttitudeInfo", 83886137, AccessTypeRead)
	KeyRobomasterSystemSightBeadPosition                = newKey("KeyRobomasterSystemSightBeadPosition", 83886138, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemSpeakerLanguage                  = newKey("KeyRobomasterSystemSpeakerLanguage", 83886139, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemSpeakerVolumn                    = newKey("KeyRobomasterSystemSpeakerVolumn", 83886140, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemChassisSpeedLevel                = newKey("KeyRobomasterSystemChassisSpeedLevel", 83886141, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemIsEncryptedFirmware              = newKey("KeyRobomasterSystemIsEncryptedFirmware", 83886142, AccessTypeRead)
	KeyRobomasterSystemScratchErrorInfo                 = newKey("KeyRobomasterSystemScratchErrorInfo", 83886143, AccessTypeRead)
	KeyRobomasterSystemScratchOutputInfo                = newKey("KeyRobomasterSystemScratchOutputInfo", 83886144, AccessTypeRead)
	KeyRobomasterSystemBarrelCoolDown                   = newKey("KeyRobomasterSystemBarrelCoolDown", 83886145, AccessTypeAction)
	KeyRobomasterSystemResetBarrelOverheat              = newKey("KeyRobomasterSystemResetBarrelOverheat", 83886146, AccessTypeAction)
	KeyRobomasterSystemMobileAccelerInfo                = newKey("KeyRobomasterSystemMobileAccelerInfo", 83886147, AccessTypeWrite)
	KeyRobomasterSystemMobileGyroAttitudeAngleInfo      = newKey("KeyRobomasterSystemMobileGyroAttitudeAngleInfo", 83886148, AccessTypeWrite)
	KeyRobomasterSystemMobileGyroRotationRateInfo       = newKey("KeyRobomasterSystemMobileGyroRotationRateInfo", 83886149, AccessTypeWrite)
	KeyRobomasterSystemEnableAcceleratorSubscribe       = newKey("KeyRobomasterSystemEnableAcceleratorSubscribe", 83886150, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemEnableGyroRotationRateSubscribe  = newKey("KeyRobomasterSystemEnableGyroRotationRateSubscribe", 83886151, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemEnableGyroAttitudeAngleSubscribe = newKey("KeyRobomasterSystemEnableGyroAttitudeAngleSubscribe", 83886152, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemDeactivate                       = newKey("KeyRobomasterSystemDeactivate", 83886153, AccessTypeAction)
	KeyRobomasterSystemFunctionEnable                   = newKey("KeyRobomasterSystemFunctionEnable", 83886154, AccessTypeAction)
	KeyRobomasterSystemIsGameRunning                    = newKey("KeyRobomasterSystemIsGameRunning", 83886155, AccessTypeRead)
	KeyRobomasterSystemIsActivated                      = newKey("KeyRobomasterSystemIsActivated", 83886156, AccessTypeRead)
	KeyRobomasterSystemLowPowerConsumption              = newKey("KeyRobomasterSystemLowPowerConsumption", 83886157, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterSystemEnterLowPowerConsumption         = newKey("KeyRobomasterSystemEnterLowPowerConsumption", 83886158, AccessTypeAction)
	KeyRobomasterSystemExitLowPowerConsumption          = newKey("KeyRobomasterSystemExitLowPowerConsumption", 83886159, AccessTypeAction)
	KeyRobomasterSystemIsLowPowerConsumption            = newKey("KeyRobomasterSystemIsLowPowerConsumption", 83886160, AccessTypeRead)
	KeyRobomasterSystemPushFile                         = newKey("KeyRobomasterSystemPushFile", 83886161, AccessTypeAction)
	KeyRobomasterSystemPlaySound                        = newKey("KeyRobomasterSystemPlaySound", 83886162, AccessTypeAction)
	KeyRobomasterSystemPlaySoundStatus                  = newKey("KeyRobomasterSystemPlaySoundStatus", 83886163, AccessTypeRead)
	KeyRobomasterSystemCustomUIAttribute                = newKey("KeyRobomasterSystemCustomUIAttribute", 83886164, AccessTypeRead)
	KeyRobomasterSystemCustomUIFunctionEvent            = newKey("KeyRobomasterSystemCustomUIFunctionEvent", 83886165, AccessTypeAction)
	KeyRobomasterSystemTotalMileage                     = newKey("KeyRobomasterSystemTotalMileage", 83886166, AccessTypeRead)
	KeyRobomasterSystemTotalDrivingTime                 = newKey("KeyRobomasterSystemTotalDrivingTime", 83886167, AccessTypeRead)
	KeyRobomasterSystemSetPlayMode                      = newKey("KeyRobomasterSystemSetPlayMode", 83886168, AccessTypeWrite)
	KeyRobomasterSystemCustomSkillInfo                  = newKey("KeyRobomasterSystemCustomSkillInfo", 83886169, AccessTypeRead)
	KeyRobomasterSystemAddressing                       = newKey("KeyRobomasterSystemAddressing", 83886170, AccessTypeAction)
	KeyRobomasterSystemLEDLightEffect                   = newKey("KeyRobomasterSystemLEDLightEffect", 83886171, AccessTypeAction)
	KeyRobomasterSystemOpenImageTransmission            = newKey("KeyRobomasterSystemOpenImageTransmission", 83886172, AccessTypeAction)
	KeyRobomasterSystemCloseImageTransmission           = newKey("KeyRobomasterSystemCloseImageTransmission", 83886173, AccessTypeAction)

	KeyRobomasterWaterGunFirmwareVersion       = newKey("KeyRobomasterWaterGunFirmwareVersion", 167772161, AccessTypeRead)
	KeyRobomasterWaterGunWaterGunFire          = newKey("KeyRobomasterWaterGunWaterGunFire", 167772162, AccessTypeAction)
	KeyRobomasterWaterGunWaterGunFireWithTimes = newKey("KeyRobomasterWaterGunWaterGunFireWithTimes", 167772163, AccessTypeAction)
	KeyRobomasterWaterGunShootSpeed            = newKey("KeyRobomasterWaterGunShootSpeed", 167772164, AccessTypeRead)
	KeyRobomasterWaterGunShootFrequency        = newKey("KeyRobomasterWaterGunShootFrequency", 167772165, AccessTypeRead)

	KeyRobomasterInfraredGunConnection      = newKey("KeyRobomasterInfraredGunConnection", 301989889, AccessTypeRead)
	KeyRobomasterInfraredGunFirmwareVersion = newKey("KeyRobomasterInfraredGunFirmwareVersion", 301989890, AccessTypeRead)
	KeyRobomasterInfraredGunInfraredGunFire = newKey("KeyRobomasterInfraredGunInfraredGunFire", 301989891, AccessTypeAction)
	KeyRobomasterInfraredGunShootFrequency  = newKey("KeyRobomasterInfraredGunShootFrequency", 301989892, AccessTypeRead)

	KeyRobomasterBatteryFirmwareVersion = newKey("KeyRobomasterBatteryFirmwareVersion", 218103809, AccessTypeRead)
	KeyRobomasterBatteryPowerPercent    = newKey("KeyRobomasterBatteryPowerPercent", 218103810, AccessTypeRead)
	KeyRobomasterBatteryVoltage         = newKey("KeyRobomasterBatteryVoltage", 218103811, AccessTypeRead)
	KeyRobomasterBatteryTemperature     = newKey("KeyRobomasterBatteryTemperature", 218103812, AccessTypeRead)
	KeyRobomasterBatteryCurrent         = newKey("KeyRobomasterBatteryCurrent", 218103813, AccessTypeRead)
	KeyRobomasterBatteryShutdown        = newKey("KeyRobomasterBatteryShutdown", 218103814, AccessTypeAction)
	KeyRobomasterBatteryReboot          = newKey("KeyRobomasterBatteryReboot", 218103815, AccessTypeAction)

	KeyRobomasterGamePadConnection                   = newKey("KeyRobomasterGamePadConnection", 234881025, AccessTypeRead)
	KeyRobomasterGamePadFirmwareVersion              = newKey("KeyRobomasterGamePadFirmwareVersion", 234881026, AccessTypeRead)
	KeyRobomasterGamePadHasMouse                     = newKey("KeyRobomasterGamePadHasMouse", 234881027, AccessTypeRead)
	KeyRobomasterGamePadHasKeyboard                  = newKey("KeyRobomasterGamePadHasKeyboard", 234881028, AccessTypeRead)
	KeyRobomasterGamePadCtrlSensitivityX             = newKey("KeyRobomasterGamePadCtrlSensitivityX", 234881029, AccessTypeWrite)
	KeyRobomasterGamePadCtrlSensitivityY             = newKey("KeyRobomasterGamePadCtrlSensitivityY", 234881030, AccessTypeWrite)
	KeyRobomasterGamePadCtrlSensitivityYaw           = newKey("KeyRobomasterGamePadCtrlSensitivityYaw", 234881031, AccessTypeWrite)
	KeyRobomasterGamePadCtrlSensitivityYawSlop       = newKey("KeyRobomasterGamePadCtrlSensitivityYawSlop", 234881032, AccessTypeWrite)
	KeyRobomasterGamePadCtrlSensitivityYawDeadZone   = newKey("KeyRobomasterGamePadCtrlSensitivityYawDeadZone", 234881033, AccessTypeWrite)
	KeyRobomasterGamePadCtrlSensitivityPitch         = newKey("KeyRobomasterGamePadCtrlSensitivityPitch", 234881034, AccessTypeWrite)
	KeyRobomasterGamePadCtrlSensitivityPitchSlop     = newKey("KeyRobomasterGamePadCtrlSensitivityPitchSlop", 234881035, AccessTypeWrite)
	KeyRobomasterGamePadCtrlSensitivityPitchDeadZone = newKey("KeyRobomasterGamePadCtrlSensitivityPitchDeadZone", 234881036, AccessTypeWrite)
	KeyRobomasterGamePadMouseLeftButton              = newKey("KeyRobomasterGamePadMouseLeftButton", 234881037, AccessTypeRead)
	KeyRobomasterGamePadMouseRightButton             = newKey("KeyRobomasterGamePadMouseRightButton", 234881038, AccessTypeRead)
	KeyRobomasterGamePadC1                           = newKey("KeyRobomasterGamePadC1", 234881039, AccessTypeRead)
	KeyRobomasterGamePadC2                           = newKey("KeyRobomasterGamePadC2", 234881040, AccessTypeRead)
	KeyRobomasterGamePadFire                         = newKey("KeyRobomasterGamePadFire", 234881041, AccessTypeRead)
	KeyRobomasterGamePadFn                           = newKey("KeyRobomasterGamePadFn", 234881042, AccessTypeRead)
	KeyRobomasterGamePadNoCalibrate                  = newKey("KeyRobomasterGamePadNoCalibrate", 234881043, AccessTypeRead)
	KeyRobomasterGamePadNotAtMiddle                  = newKey("KeyRobomasterGamePadNotAtMiddle", 234881044, AccessTypeRead)
	KeyRobomasterGamePadBatteryWarning               = newKey("KeyRobomasterGamePadBatteryWarning", 234881045, AccessTypeRead)
	KeyRobomasterGamePadBatteryPercent               = newKey("KeyRobomasterGamePadBatteryPercent", 234881046, AccessTypeRead)
	KeyRobomasterGamePadActivationSettings           = newKey("KeyRobomasterGamePadActivationSettings", 234881047, AccessTypeRead|AccessTypeWrite)
	KeyRobomasterGamePadControlEnabled               = newKey("KeyRobomasterGamePadControlEnabled", 234881048, AccessTypeWrite)

	KeyRobomasterClawConnection          = newKey("KeyRobomasterClawConnection", 251658241, AccessTypeRead)
	KeyRobomasterClawFirmwareVersion     = newKey("KeyRobomasterClawFirmwareVersion", 251658242, AccessTypeRead)
	KeyRobomasterClawCtrl                = newKey("KeyRobomasterClawCtrl", 251658243, AccessTypeAction)
	KeyRobomasterClawStatus              = newKey("KeyRobomasterClawStatus", 251658244, AccessTypeRead)
	KeyRobomasterClawInfoSubscribe       = newKey("KeyRobomasterClawInfoSubscribe", 251658245, AccessTypeRead)
	KeyRobomasterEnableClawInfoSubscribe = newKey("KeyRobomasterEnableClawInfoSubscribe", 251658246, AccessTypeAction)

	KeyRobomasterArmConnection          = newKey("KeyRobomasterArmConnection", 285212673, AccessTypeRead)
	KeyRobomasterArmCtrl                = newKey("KeyRobomasterArmCtrl", 285212674, AccessTypeAction)
	KeyRobomasterArmCtrlMode            = newKey("KeyRobomasterArmCtrlMode", 285212675, AccessTypeAction)
	KeyRobomasterArmCalibration         = newKey("KeyRobomasterArmCalibration", 285212676, AccessTypeAction)
	KeyRobomasterArmBlockedFlag         = newKey("KeyRobomasterArmBlockedFlag", 285212677, AccessTypeRead)
	KeyRobomasterArmPositionSubscribe   = newKey("KeyRobomasterArmPositionSubscribe", 285212678, AccessTypeRead)
	KeyRobomasterArmReachLimitX         = newKey("KeyRobomasterArmReachLimitX", 285212679, AccessTypeRead)
	KeyRobomasterArmReachLimitY         = newKey("KeyRobomasterArmReachLimitY", 285212680, AccessTypeRead)
	KeyRobomasterEnableArmInfoSubscribe = newKey("KeyRobomasterEnableArmInfoSubscribe", 285212681, AccessTypeAction)
	KeyRobomasterArmControlMode         = newKey("KeyRobomasterArmControlMode", 285212682, AccessTypeRead|AccessTypeWrite)

	KeyRobomasterTOFConnection          = newKey("KeyRobomasterTOFConnection", 318767105, AccessTypeRead)
	KeyRobomasterTOFLEDColor            = newKey("KeyRobomasterTOFLEDColor", 318767106, AccessTypeWrite)
	KeyRobomasterTOFOnlineModules       = newKey("KeyRobomasterTOFOnlineModules", 318767107, AccessTypeRead)
	KeyRobomasterTOFInfoSubscribe       = newKey("KeyRobomasterTOFInfoSubscribe", 318767108, AccessTypeRead)
	KeyRobomasterEnableTOFInfoSubscribe = newKey("KeyRobomasterEnableTOFInfoSubscribe", 318767109, AccessTypeAction)
	KeyRobomasterTOFFirmwareVersion1    = newKey("KeyRobomasterTOFFirmwareVersion1", 318767110, AccessTypeRead)
	KeyRobomasterTOFFirmwareVersion2    = newKey("KeyRobomasterTOFFirmwareVersion2", 318767111, AccessTypeRead)
	KeyRobomasterTOFFirmwareVersion3    = newKey("KeyRobomasterTOFFirmwareVersion3", 318767112, AccessTypeRead)
	KeyRobomasterTOFFirmwareVersion4    = newKey("KeyRobomasterTOFFirmwareVersion4", 318767113, AccessTypeRead)

	KeyRobomasterServoConnection          = newKey("KeyRobomasterServoConnection", 335544321, AccessTypeRead)
	KeyRobomasterServoLEDColor            = newKey("KeyRobomasterServoLEDColor", 335544322, AccessTypeWrite)
	KeyRobomasterServoSpeed               = newKey("KeyRobomasterServoSpeed", 335544323, AccessTypeWrite)
	KeyRobomasterServoOnlineModules       = newKey("KeyRobomasterServoOnlineModules", 335544324, AccessTypeRead)
	KeyRobomasterServoInfoSubscribe       = newKey("KeyRobomasterServoInfoSubscribe", 335544325, AccessTypeRead)
	KeyRobomasterEnableServoInfoSubscribe = newKey("KeyRobomasterEnableServoInfoSubscribe", 335544326, AccessTypeAction)
	KeyRobomasterServoFirmwareVersion1    = newKey("KeyRobomasterServoFirmwareVersion1", 335544327, AccessTypeRead)
	KeyRobomasterServoFirmwareVersion2    = newKey("KeyRobomasterServoFirmwareVersion2", 335544328, AccessTypeRead)
	KeyRobomasterServoFirmwareVersion3    = newKey("KeyRobomasterServoFirmwareVersion3", 335544329, AccessTypeRead)
	KeyRobomasterServoFirmwareVersion4    = newKey("KeyRobomasterServoFirmwareVersion4", 335544330, AccessTypeRead)

	KeyRobomasterSensorAdapterConnection          = newKey("KeyRobomasterSensorAdapterConnection", 352321537, AccessTypeRead)
	KeyRobomasterSensorAdapterOnlineModules       = newKey("KeyRobomasterSensorAdapterOnlineModules", 352321538, AccessTypeRead)
	KeyRobomasterSensorAdapterInfoSubscribe       = newKey("KeyRobomasterSensorAdapterInfoSubscribe", 352321539, AccessTypeRead)
	KeyRobomasterEnableSensorAdapterInfoSubscribe = newKey("KeyRobomasterEnableSensorAdapterInfoSubscribe", 352321540, AccessTypeAction)
	KeyRobomasterSensorAdapterFirmwareVersion1    = newKey("KeyRobomasterSensorAdapterFirmwareVersion1", 352321541, AccessTypeRead)
	KeyRobomasterSensorAdapterFirmwareVersion2    = newKey("KeyRobomasterSensorAdapterFirmwareVersion2", 352321542, AccessTypeRead)
	KeyRobomasterSensorAdapterFirmwareVersion3    = newKey("KeyRobomasterSensorAdapterFirmwareVersion3", 352321543, AccessTypeRead)
	KeyRobomasterSensorAdapterFirmwareVersion4    = newKey("KeyRobomasterSensorAdapterFirmwareVersion4", 352321544, AccessTypeRead)
	KeyRobomasterSensorAdapterFirmwareVersion5    = newKey("KeyRobomasterSensorAdapterFirmwareVersion5", 352321545, AccessTypeRead)
	KeyRobomasterSensorAdapterFirmwareVersion6    = newKey("KeyRobomasterSensorAdapterFirmwareVersion6", 352321546, AccessTypeRead)
	KeyRobomasterSensorAdapterLEDColor            = newKey("KeyRobomasterSensorAdapterLEDColor", 352321547, AccessTypeWrite)

	KeyRemoteControllerConnection = newKey("KeyRemoteControllerConnection", 50331649, AccessTypeRead)

	KeyGimbalConnection              = newKey("KeyGimbalConnection", 67108865, AccessTypeRead)
	KeyGimbalESCFirmwareVersion      = newKey("KeyGimbalESCFirmwareVersion", 67108866, AccessTypeRead)
	KeyGimbalFirmwareVersion         = newKey("KeyGimbalFirmwareVersion", 67108867, AccessTypeRead)
	KeyGimbalWorkMode                = newKey("KeyGimbalWorkMode", 67108868, AccessTypeRead|AccessTypeWrite)
	KeyGimbalControlMode             = newKey("KeyGimbalControlMode", 67108869, AccessTypeRead|AccessTypeWrite)
	KeyGimbalResetPosition           = newKey("KeyGimbalResetPosition", 67108870, AccessTypeAction)
	KeyGimbalResetPositionState      = newKey("KeyGimbalResetPositionState", 67108871, AccessTypeRead)
	KeyGimbalCalibration             = newKey("KeyGimbalCalibration", 67108872, AccessTypeAction)
	KeyGimbalSpeedRotation           = newKey("KeyGimbalSpeedRotation", 67108873, AccessTypeAction)
	KeyGimbalSpeedRotationEnabled    = newKey("KeyGimbalSpeedRotationEnabled", 67108874, AccessTypeWrite)
	KeyGimbalAngleIncrementRotation  = newKey("KeyGimbalAngleIncrementRotation", 67108875, AccessTypeAction)
	KeyGimbalAngleFrontYawRotation   = newKey("KeyGimbalAngleFrontYawRotation", 67108876, AccessTypeAction)
	KeyGimbalAngleFrontPitchRotation = newKey("KeyGimbalAngleFrontPitchRotation", 67108877, AccessTypeAction)
	KeyGimbalAttitude                = newKey("KeyGimbalAttitude", 67108878, AccessTypeRead)
	KeyGimbalAutoCalibrate           = newKey("KeyGimbalAutoCalibrate", 67108879, AccessTypeAction)
	KeyGimbalCalibrationStatus       = newKey("KeyGimbalCalibrationStatus", 67108880, AccessTypeRead)
	KeyGimbalCalibrationProgress     = newKey("KeyGimbalCalibrationProgress", 67108881, AccessTypeRead)
	KeyGimbalOpenAttitudeUpdates     = newKey("KeyGimbalOpenAttitudeUpdates", 67108882, AccessTypeAction)
	KeyGimbalCloseAttitudeUpdates    = newKey("KeyGimbalCloseAttitudeUpdates", 67108883, AccessTypeAction)
	KeyGimbalGetLinkAck              = newKey("KeyGimbalGetLinkAck", 83886092, AccessTypeRead)

	KeyVisionFirmwareVersion             = newKey("KeyVisionFirmwareVersion", 100663297, AccessTypeRead)
	KeyVisionTrackingAutoLockTarget      = newKey("KeyVisionTrackingAutoLockTarget", 100663298, AccessTypeRead|AccessTypeWrite)
	KeyVisionARParameters                = newKey("KeyVisionARParameters", 100663299, AccessTypeRead)
	KeyVisionARTagEnabled                = newKey("KeyVisionARTagEnabled", 100663300, AccessTypeRead)
	KeyVisionDebugRect                   = newKey("KeyVisionDebugRect", 100663301, AccessTypeRead)
	KeyVisionLaserPosition               = newKey("KeyVisionLaserPosition", 100663302, AccessTypeRead)
	KeyVisionDetectionEnable             = newKey("KeyVisionDetectionEnable", 100663303, AccessTypeRead|AccessTypeWrite)
	KeyVisionMarkerRunningStatus         = newKey("KeyVisionMarkerRunningStatus", 100663304, AccessTypeRead)
	KeyVisionTrackingRunningStatus       = newKey("KeyVisionTrackingRunningStatus", 100663305, AccessTypeRead)
	KeyVisionAimbotRunningStatus         = newKey("KeyVisionAimbotRunningStatus", 100663306, AccessTypeRead)
	KeyVisionHeadAndShoulderStatus       = newKey("KeyVisionHeadAndShoulderStatus", 100663307, AccessTypeRead)
	KeyVisionHumanDetectionRunningStatus = newKey("KeyVisionHumanDetectionRunningStatus", 100663308, AccessTypeRead)
	KeyVisionUserConfirm                 = newKey("KeyVisionUserConfirm", 100663309, AccessTypeAction)
	KeyVisionUserCancel                  = newKey("KeyVisionUserCancel", 100663310, AccessTypeAction)
	KeyVisionUserTrackingRect            = newKey("KeyVisionUserTrackingRect", 100663311, AccessTypeWrite)
	KeyVisionTrackingDistance            = newKey("KeyVisionTrackingDistance", 100663312, AccessTypeWrite)
	KeyVisionLineColor                   = newKey("KeyVisionLineColor", 100663313, AccessTypeWrite)
	KeyVisionMarkerColor                 = newKey("KeyVisionMarkerColor", 100663314, AccessTypeWrite)
	KeyVisionMarkerAdvanceStatus         = newKey("KeyVisionMarkerAdvanceStatus", 100663315, AccessTypeRead)

	KeyPerceptionFirmwareVersion = newKey("KeyPerceptionFirmwareVersion", 184549377, AccessTypeRead)
	KeyPerceptionMarkerEnable    = newKey("KeyPerceptionMarkerEnable", 184549378, AccessTypeRead|AccessTypeWrite)
	KeyPerceptionMarkerResult    = newKey("KeyPerceptionMarkerResult", 184549379, AccessTypeRead)

	KeyESCFirmwareVersion1 = newKey("KeyESCFirmwareVersion1", 201326593, AccessTypeRead)
	KeyESCFirmwareVersion2 = newKey("KeyESCFirmwareVersion2", 201326594, AccessTypeRead)
	KeyESCFirmwareVersion3 = newKey("KeyESCFirmwareVersion3", 201326595, AccessTypeRead)
	KeyESCFirmwareVersion4 = newKey("KeyESCFirmwareVersion4", 201326596, AccessTypeRead)
	KeyESCMotorInfomation1 = newKey("KeyESCMotorInfomation1", 201326597, AccessTypeRead)
	KeyESCMotorInfomation2 = newKey("KeyESCMotorInfomation2", 201326598, AccessTypeRead)
	KeyESCMotorInfomation3 = newKey("KeyESCMotorInfomation3", 201326599, AccessTypeRead)
	KeyESCMotorInfomation4 = newKey("KeyESCMotorInfomation4", 201326600, AccessTypeRead)

	KeyWiFiLinkFirmwareVersion         = newKey("KeyWiFiLinkFirmwareVersion", 134217729, AccessTypeRead)
	KeyWiFiLinkDebugInfo               = newKey("KeyWiFiLinkDebugInfo", 134217730, AccessTypeRead)
	KeyWiFiLinkMode                    = newKey("KeyWiFiLinkMode", 134217731, AccessTypeRead)
	KeyWiFiLinkSSID                    = newKey("KeyWiFiLinkSSID", 134217732, AccessTypeRead|AccessTypeWrite)
	KeyWiFiLinkPassword                = newKey("KeyWiFiLinkPassword", 134217733, AccessTypeRead|AccessTypeWrite)
	KeyWiFiLinkAvailableChannelNumbers = newKey("KeyWiFiLinkAvailableChannelNumbers", 134217734, AccessTypeRead)
	KeyWiFiLinkCurrentChannelNumber    = newKey("KeyWiFiLinkCurrentChannelNumber", 134217735, AccessTypeRead|AccessTypeWrite)
	KeyWiFiLinkSNR                     = newKey("KeyWiFiLinkSNR", 134217736, AccessTypeRead)
	KeyWiFiLinkSNRPushEnabled          = newKey("KeyWiFiLinkSNRPushEnabled", 134217737, AccessTypeWrite)
	KeyWiFiLinkReboot                  = newKey("KeyWiFiLinkReboot", 134217738, AccessTypeAction)
	KeyWiFiLinkChannelSelectionMode    = newKey("KeyWiFiLinkChannelSelectionMode", 134217739, AccessTypeRead|AccessTypeWrite)
	KeyWiFiLinkInterference            = newKey("KeyWiFiLinkInterference", 134217740, AccessTypeRead)
	KeyWiFiLinkDeleteNetworkConfig     = newKey("KeyWiFiLinkDeleteNetworkConfig", 134217741, AccessTypeAction)

	KeySDRLinkSNR                  = newKey("KeySDRLinkSNR", 268435457, AccessTypeRead)
	KeySDRLinkBandwidth            = newKey("KeySDRLinkBandwidth", 268435458, AccessTypeRead|AccessTypeWrite)
	KeySDRLinkChannelSelectionMode = newKey("KeySDRLinkChannelSelectionMode", 268435459, AccessTypeRead|AccessTypeWrite)
	KeySDRLinkCurrentFreqPoint     = newKey("KeySDRLinkCurrentFreqPoint", 268435460, AccessTypeRead|AccessTypeWrite)
	KeySDRLinkCurrentFreqBand      = newKey("KeySDRLinkCurrentFreqBand", 268435461, AccessTypeRead|AccessTypeWrite)
	KeySDRLinkIsDualFreqSupported  = newKey("KeySDRLinkIsDualFreqSupported", 268435462, AccessTypeRead)
	KeySDRLinkUpdateConfigs        = newKey("KeySDRLinkUpdateConfigs", 268435463, AccessTypeAction)

	KeyAirLinkConnection         = newKey("KeyAirLinkConnection", 117440513, AccessTypeRead)
	KeyAirLinkSignalQuality      = newKey("KeyAirLinkSignalQuality", 117440514, AccessTypeRead)
	KeyAirLinkCountryCode        = newKey("KeyAirLinkCountryCode", 117440515, AccessTypeWrite)
	KeyAirLinkCountryCodeUpdated = newKey("KeyAirLinkCountryCodeUpdated", 117440516, AccessTypeRead)

	KeyArmorFirmwareVersion1 = newKey("KeyArmorFirmwareVersion1", 150994945, AccessTypeRead)
	KeyArmorFirmwareVersion2 = newKey("KeyArmorFirmwareVersion2", 150994946, AccessTypeRead)
	KeyArmorFirmwareVersion3 = newKey("KeyArmorFirmwareVersion3", 150994947, AccessTypeRead)
	KeyArmorFirmwareVersion4 = newKey("KeyArmorFirmwareVersion4", 150994948, AccessTypeRead)
	KeyArmorFirmwareVersion5 = newKey("KeyArmorFirmwareVersion5", 150994949, AccessTypeRead)
	KeyArmorFirmwareVersion6 = newKey("KeyArmorFirmwareVersion6", 150994950, AccessTypeRead)
	KeyArmorUnderAttack      = newKey("KeyArmorUnderAttack", 150994951, AccessTypeRead)
	KeyArmorEnterResetID     = newKey("KeyArmorEnterResetID", 150994952, AccessTypeAction)
	KeyArmorCancelResetID    = newKey("KeyArmorCancelResetID", 150994953, AccessTypeAction)
	KeyArmorSkipCurrentID    = newKey("KeyArmorSkipCurrentID", 150994954, AccessTypeAction)
	KeyArmorResetStatus      = newKey("KeyArmorResetStatus", 150994955, AccessTypeRead)
)

// String returns the name of the
func (k *Key) String() string {
	return k.name
}

// SubType returns the sub-type associated wit hhis key. Used for events.
func (k *Key) SubType() uint32 {
	return k.subType
}

// AccessType returns the access type of the
func (k *Key) AccessType() AccessType {
	return k.accessType
}

// FromEvent returns a Key associated with the given event. It returns
// an error in case the key cna not be inferred.
func FromEvent(ev *event.Event) (*Key, error) {
	k, ok := keyBySubType[ev.SubType()]
	if !ok {
		return nil, fmt.Errorf("event sub-type does not match any key: %d",
			ev.SubType())
	}

	return k, nil
}

func newKey(name string, subType uint32, accessType AccessType) *Key {
	k := &Key{
		name:       name,
		subType:    subType,
		accessType: accessType,
	}

	keyBySubType[subType] = k

	return k
}
