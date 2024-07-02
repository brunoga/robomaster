package gimbal

import (
	"fmt"
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/support/token"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

// Gimbal is the module that allows controlling the gimbal.
type Gimbal struct {
	*internal.BaseModule

	gaToken token.Token

	controlMode ControlMode
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
// degrees per second. The gimbal will move until either it can not move anymore
// due to pgysical constraints or StopRotation() is called.
func (g *Gimbal) SetRotationSpeed(pitch, yaw int16) error {
	if pitch < -360 || pitch > 360 {
		return fmt.Errorf("invalid pitch value %d", pitch)
	}
	if yaw < -360 || yaw > 360 {
		return fmt.Errorf("invalid yaw value %d", yaw)
	}

	// TODO(bga): Check if this is needed all the time.
	err := g.UB().PerformActionForKey(key.KeyGimbalSpeedRotationEnabled, &value.Uint64{Value: 1}, nil)
	if err != nil {
		return err
	}

	return g.UB().PerformActionForKey(key.KeyGimbalSpeedRotation,
		&value.GimbalSpeedRotation{Pitch: pitch * 10, Yaw: yaw * 10, Roll: 0}, nil)
}

// SetRelativeAngleRotation sets the gimbal rotation relative to the current
// position. This is executed asynchronously.
//
// TODO(bga): Figure out units.
func (g *Gimbal) SetRelativeAngleRotation(angle int16, axis Axis,
	duration time.Duration) error {
	gimbalIncrementRotation := value.GimbalAngleRotation{
		Time: int16(duration / time.Millisecond),
	}

	if axis == AxisPitch {
		if angle < -60 || angle > 60 {
			return fmt.Errorf("invalid pitch angle %d", angle)
		}

		gimbalIncrementRotation.Pitch = angle * 10
		gimbalIncrementRotation.Yaw = 0
	} else {
		// TODO(bga): Fix this. It might be just something that needs to be set
		//            before this is called, like the the chassis or gimbal
		//            modes.
		gimbalIncrementRotation.Pitch = 0
		gimbalIncrementRotation.Yaw = angle * 10
	}

	return g.UB().PerformActionForKeySync(key.KeyGimbalAngleIncrementRotation,
		&gimbalIncrementRotation)
}

// SetAbsoluteAngleRotation sets the absolute gimbal rotation relative to its
// default position. This is executed asynchronously.
func (g *Gimbal) SetAbsoluteAngleRotation(angle int16, axis Axis,
	duration time.Duration) error {
	gimbalAngleRotation := value.GimbalAngleRotation{
		Time: int16(duration / time.Millisecond),
	}

	var k *key.Key

	if axis == AxisPitch {
		if angle < -25 || angle > 35 {
			return fmt.Errorf("invalid pitch angle %d", angle)
		}

		gimbalAngleRotation.Pitch = angle * 10
		gimbalAngleRotation.Yaw = 0
		k = key.KeyGimbalAngleFrontPitchRotation
	} else {
		gimbalAngleRotation.Pitch = 0
		gimbalAngleRotation.Yaw = angle * 10
		k = key.KeyGimbalAngleFrontYawRotation
	}

	return g.UB().PerformActionForKeySync(k, &gimbalAngleRotation)
}

// StopRotation stops any ongoing gimbal rotation.
func (g *Gimbal) StopRotation() error {
	err := g.SetRotationSpeed(0, 0)
	if err != nil {
		return err
	}

	return g.UB().PerformActionForKey(key.KeyGimbalSpeedRotationEnabled, &value.Uint64{Value: 0}, nil)
}

// ResetPosition resets the gimbal position.
func (g *Gimbal) ResetPosition() error {
	err := g.UB().PerformActionForKeySync(key.KeyGimbalResetPosition, nil)
	if err != nil {
		return err
	}

	gimbalReset := 0
	c := make(chan struct{})

	t, err := g.UB().AddKeyListener(key.KeyGimbalResetPositionState, func(r *result.Result) {
		g.Logger().Debug("Reset position state", "result", r)
		if !r.Succeeded() {
			g.Logger().Error("Error resetting gimbal position", "error", r.ErrorDesc())
			return
		}

		value, ok := r.Value().(*value.Uint64)
		if !ok {
			g.Logger().Error("Unexpected value", "key", r.Key(), "value", r.Value())
			return
		}

		if value.Value == 0 && gimbalReset == 1 {
			g.Logger().Debug("Reset position done")
			close(c)
			return
		}

		gimbalReset = int(value.Value)
	}, true)
	defer g.UB().RemoveKeyListener(key.KeyGimbalResetPositionState, t)

	<-c

	return nil
}

func (g *Gimbal) ControlMode() ControlMode {
	return g.controlMode
}

// SetControlMode sets the gimbal control mode.
func (g *Gimbal) SetControlMode(cm ControlMode) error {
	if !cm.Valid() {
		return fmt.Errorf("invalid control mode %d", cm)
	}

	if cm == ControlMode3 {
		cm = ControlMode1
	}

	err := g.UB().DirectSendKeyValue(key.KeyGimbalControlMode, uint64(cm))
	if err != nil {
		return err
	}

	g.controlMode = cm

	return nil
}

func (g *Gimbal) WorkMode() (uint64, error) {
	r, err := g.UB().GetKeyValueSync(key.KeyGimbalWorkMode, true)
	if err != nil {
		return 0, err
	}

	if !r.Succeeded() {
		return 0, fmt.Errorf("error getting work mode: %s", r.ErrorDesc())
	}

	wm, ok := r.Value().(*value.Uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected value: %v", r.Value())
	}

	return wm.Value, nil
}

// SetWorkMode sets the gimbal work mode.
func (g *Gimbal) SetWorkMode(wm uint64) error {
	return g.UB().SetKeyValueSync(key.KeyGimbalWorkMode, &value.Uint64{Value: wm})
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
