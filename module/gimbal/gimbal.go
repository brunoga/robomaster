package gimbal

import (
	"fmt"
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
	"github.com/brunoga/robomaster/unitybridge/support/token"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

// Gimbal is the module that allows controlling the gimbal.
type Gimbal struct {
	*internal.BaseModule

	gaToken token.Token
}

var _ module.Module = (*Gimbal)(nil)

// New creates a new Gimbal instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection) (*Gimbal, error) {
	g := &Gimbal{}

	g.BaseModule = internal.NewBaseModule(ub, l, "Gimbal",
		key.KeyGimbalConnection, func(r *result.Result) {
			if r == nil || !r.Succeeded() {
				g.Logger().Error("Error connecting to gimbal", "error", r.ErrorDesc())
				return
			}

			connected, ok := r.Value().(*value.Bool)
			if !ok {
				g.Logger().Error("Unexpected value", "key", r.Key(), "value", r.Value())
				return
			}

			if connected.Value {
				err := g.UB().PerformActionForKeySync(key.KeyGimbalOpenAttitudeUpdates, nil)
				if err != nil {
					g.Logger().Error("Error opening attitude updates", "error", err)
				}
			} else {
				err := g.UB().PerformActionForKeySync(key.KeyGimbalCloseAttitudeUpdates, nil)
				if err != nil {
					g.Logger().Error("Error closing attitude updates", "error", err)
				}
			}
		}, cm)

	return g, nil
}

func (g *Gimbal) Start() error {
	var err error

	g.gaToken, err = g.UB().AddKeyListener(key.KeyGimbalAttitude,
		g.onAttitudeUpdates, false)
	if err != nil {
		return err
	}

	return g.BaseModule.Start()
}

// SetRotationSpeed sets the gimbal rotation speed for the pitch and yaw axis in
// degrees per second.
func (g *Gimbal) SetRotationSpeed(pitch, yaw int16) error {
	if pitch < -360 || pitch > 360 {
		return fmt.Errorf("invalid pitch value %d", pitch)
	}
	if yaw < -360 || yaw > 360 {
		return fmt.Errorf("invalid yaw value %d", yaw)
	}

	// TODO(bga): Check if this is needed all the time.
	//err := g.UB().PerformActionForKeySync(key.KeyGimbalSpeedRotationEnabled, &value.Uint64{Value: 1})
	//if err != nil {
	//	return err
	//}

	return g.UB().PerformActionForKeySync(key.KeyGimbalSpeedRotation,
		&value.GimbalSpeedRotation{Pitch: pitch, Yaw: yaw})
}

// SetRelativeAngleRotation sets the gimbal rotation relative to the current
// position. This is executed asynchronously.
//
// TODO(bga): Figure out units.
func (g *Gimbal) SetRelativeAngleRotation(pitch, yaw int16,
	duration time.Duration) error {
	//if pitch < -60 || pitch > 60 {
	//	return fmt.Errorf("invalid pitch value %d", pitch)
	//}

	return g.UB().PerformActionForKeySync(key.KeyGimbalAngleIncrementRotation,
		&value.GimbalAngleRotation{Pitch: pitch, Yaw: yaw,
			Time: int16(duration * time.Second)})
}

// SetAbsoluteAngleRotation sets the absolute gimbal rotation relative to its
// default position. This is executed asynchronously.
func (g *Gimbal) SetAbsoluteAngleRotation(pitch, yaw int16,
	duration time.Duration) error {
	//if pitch < -25 || pitch > 35 {
	//	return fmt.Errorf("invalid pitch value %d", pitch)
	//}

	// Unfortunatelly it seems that there is no way to set the absolute position
	// for both pitch and yaw at the same time. So we need to do it in two
	// steps.

	// Set pitch.
	err := g.UB().PerformActionForKeySync(key.KeyGimbalAngleFrontPitchRotation,
		&value.GimbalAngleRotation{Pitch: pitch, Yaw: yaw,
			Time: int16(duration * time.Second)})
	if err != nil {
		return err
	}

	// Set yaw,
	return g.UB().PerformActionForKeySync(key.KeyGimbalAngleFrontYawRotation,
		&value.GimbalAngleRotation{Pitch: pitch, Yaw: yaw,
			Time: int16(duration * time.Second)})
}

// StopRotation stops any ongoing gimbal rotation.
func (g *Gimbal) StopRotation() error {
	err := g.SetRotationSpeed(0, 0)
	if err != nil {
		return err
	}

	return g.UB().PerformActionForKeySync(key.KeyGimbalSpeedRotationEnabled, 0)
}

// ResetPosition resets the gimbal position.
func (g *Gimbal) ResetPosition() {
	g.UB().PerformActionForKeySync(key.KeyGimbalResetPosition, nil)
}

func (g *Gimbal) ControlMode() (uint64, error) {
	r, err := g.UB().GetKeyValueSync(key.KeyGimbalControlMode, true)
	if err != nil {
		return 0, err
	}

	if !r.Succeeded() {
		return 0, fmt.Errorf("error getting control mode: %s", r.ErrorDesc())
	}

	cm, ok := r.Value().(*value.Uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return cm.Value, nil
}

// SetControlMode sets the gimbal control mode.
func (g *Gimbal) SetControlMode(cm ControlMode) error {
	if !cm.Valid() {
		return fmt.Errorf("invalid control mode %d", cm)
	}

	if cm == ControlMode3 {
		cm = ControlMode1
	}

	return g.UB().DirectSendKeyValue(key.KeyGimbalControlMode, uint64(cm))
}

func (g *Gimbal) Stop() error {
	err := g.StopRotation()
	if err != nil {
		return err
	}

	err = g.UB().RemoveKeyListener(key.KeyGimbalAttitude, g.gaToken)
	if err != nil {
		return err
	}

	return g.BaseModule.Stop()
}

func (g *Gimbal) onAttitudeUpdates(r *result.Result) {
	if r == nil || !r.Succeeded() {
		g.Logger().Error("Error getting gimbal attitude", "error", r.ErrorDesc())
		return
	}

	value, ok := r.Value().(*value.GimbalAttitude)
	if !ok {
		g.Logger().Error("Unexpected result value", "key", r.Key(), "value", value)
		return
	}

	g.Logger().Info("Gimbal attitude", "pitch", value.Pitch, "roll", value.Roll,
		"yaw", value.Yaw, "yawOpposite", value.YawOpposite, "pitchSpeed",
		value.PitchSpeed, "rollSpeed", value.RollSpeed, "yawSpeed", value.YawSpeed)
}
