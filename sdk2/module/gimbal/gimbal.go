package gimbal

import (
	"fmt"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
	"github.com/brunoga/unitybridge/unity/result/value"
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
				err := g.UB().PerformActionForKey(key.KeyGimbalOpenAttitudeUpdates, nil, nil)
				if err != nil {
					g.Logger().Error("Error opening attitude updates", "error", err)
				}
			} else {
				err := g.UB().PerformActionForKey(key.KeyGimbalCloseAttitudeUpdates, nil, nil)
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

// ResetPosition resets the gimbal position.
func (g *Gimbal) ResetPosition() {
	g.UB().PerformActionForKeySync(key.KeyGimbalResetPosition, nil)
}

type gimbalSpeedRotation struct {
	Pitch int16 `json:"pitch"`
	Roll  int16 `json:"roll"` // unused
	Yaw   int16 `json:"yaw"`
}

// SetSpeed sets the gimbal speed.
//
// TODO(bga): Figure out units.
func (g *Gimbal) SetSpeed(pitch, yaw int16) error {
	// TODO(bga): Check if this is needed all the time.
	err := g.UB().PerformActionForKey(key.KeyGimbalSpeedRotationEnabled, 1, nil)
	if err != nil {
		return err
	}

	return g.UB().PerformActionForKey(key.KeyGimbalSpeedRotation,
		gimbalSpeedRotation{Pitch: pitch, Yaw: yaw}, nil)
}

type gimbalAngleRotation struct {
	Pitch int16 `json:"pitch"`
	Yaw   int16 `json:"yaw"`
	Time  int16 `json:"time"`
}

// SetRelativePosition sets the gimbal position relative to the current
// position. This is executed asynchronously.
//
// TODO(bga): Figure out units.
func (g *Gimbal) SetRelativePosition(pitch, yaw int16,
	duration time.Duration) error {
	return g.UB().PerformActionForKey(key.KeyGimbalAngleIncrementRotation,
		gimbalAngleRotation{Pitch: pitch, Yaw: yaw,
			Time: int16(duration * time.Second)}, nil)
}

// SetAbsolutePosition sets the absolute gimbal position in relation to its
// default position. This is executed asynchronously.
func (g *Gimbal) SetAbsolutePosition(pitch, yaw int16,
	duration time.Duration) error {
	// Unfortunatelly it seems that there is no way to set the absolute position
	// for both pitch and yaw at the same time. So we need to do it in two
	// steps.

	// Set pitch.
	err := g.UB().PerformActionForKey(key.KeyGimbalAngleFrontPitchRotation,
		gimbalAngleRotation{Pitch: pitch, Yaw: yaw,
			Time: int16(duration * time.Second)}, nil)
	if err != nil {
		return err
	}

	// Set yaw
	return g.UB().PerformActionForKey(key.KeyGimbalAngleFrontYawRotation,
		gimbalAngleRotation{Pitch: pitch, Yaw: yaw,
			Time: int16(duration * time.Second)}, nil)
}

// StopMovement stops any ongoing gimbal movement.
func (g *Gimbal) StopMovement() error {
	err := g.SetSpeed(0, 0)
	if err != nil {
		return err
	}

	return g.UB().PerformActionForKey(key.KeyGimbalSpeedRotationEnabled, 0, nil)
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
	err := g.StopMovement()
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

	value, ok := r.Value().(map[string]interface{})
	if !ok {
		g.Logger().Error("Unexpected result value", "key", r.Key(), "value", value)
		return
	}

	pitch := g.tryGetFloat64Field(value, "pitch")
	roll := g.tryGetFloat64Field(value, "roll")
	yaw := g.tryGetFloat64Field(value, "yaw")
	yawOpposite := g.tryGetFloat64Field(value, "yawOpposite")
	pitchSpeed := g.tryGetFloat64Field(value, "pitchSpeed")
	rollSpeed := g.tryGetFloat64Field(value, "rollSpeed")
	yawSpeed := g.tryGetFloat64Field(value, "yawSpeed")

	g.Logger().Info("Gimbal attitude", "pitch", pitch, "roll", roll, "yaw", yaw,
		"yawOpposite", yawOpposite, "pitchSpeed", pitchSpeed, "rollSpeed", rollSpeed,
		"yawSpeed", yawSpeed)
}

func (g *Gimbal) tryGetFloat64Field(value map[string]interface{}, field string) float64 {
	v, ok := value[field]
	if !ok {
		g.Logger().Error("Missing field", "field", field)
		return 0.0
	}

	f, ok := v.(float64)
	if !ok {
		g.Logger().Error("Unexpected value", "field", field, "value", v)
		return 0.0
	}

	return f
}
