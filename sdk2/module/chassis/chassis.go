package chassis

import (
	"fmt"
	"log/slog"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/chassis/controller"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/robomaster/sdk2/module/robot"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

// Chassis allows controlling the robot chassis. It also works as the robot main
// controller interface.
type Chassis struct {
	*internal.BaseModule

	r *robot.Robot
}

var _ module.Module = (*Chassis)(nil)

// New creates a new Chassis instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	r *robot.Robot) (*Chassis, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("chassis_module")

	c := &Chassis{
		r: r,
	}

	c.BaseModule = internal.NewBaseModule(ub, l, "Chassis",
		key.KeyMainControllerConnection, func(r *result.Result) {
			if r == nil || r.ErrorCode() != 0 {
				return
			}

			if connected, ok := r.Value().(bool); !ok || !connected {
				return
			}

			// TODO(bga): Maybe disable the function if we receive an actual
			//            false here?

			c.r.EnableFunction(robot.FunctionTypeMovementControl, true)
			c.SetControllerMode(controller.ModeFPV) // Seems to be the default mode.
		})

	return c, nil
}

// SetControllerMode sets the controller mode for the robot.
func (c *Chassis) SetControllerMode(m controller.Mode) error {
	if !m.Valid() {
		return fmt.Errorf("invalid controller mode: %d", m)
	}

	return c.UB().SetKeyValueSync(key.KeyMainControllerChassisCarControlMode, uint64(m))
}

// SetMode sets the chassis mode for the robot.
func (c *Chassis) SetMode(m Mode) error {
	// TODO(bga): Figure out this value.
	value := uint64(1) | uint64(140) | uint64(17920) | uint64(235929600)

	return c.control(m, value)

	// TODO(bga): Apparently we need to stop the chassis mode before returning.
	//            Check if that is indeed the case.
}

// StopMovement stops the chassis movement.
func (c *Chassis) StopMovement(m Mode) error {
	// TODO(bga): Figure out this value.
	value := uint64(0) | uint64(140) | uint64(17920) | uint64(235929600)

	return c.control(m, value)
}

// SetSpeed sets the chassis speed.
func (c *Chassis) SetSpeed(m Mode, x, y, z float64) error {
	value := uint64(1) |
		uint64(x*10.0) + 35<<2 |
		uint64(y*10.0) + 35<<9 |
		uint64(z*10.0) + 3600<<16

	return c.control(m, value)
}

type chassisPosition struct {
	TaskID      uint8   `json:"taskId"`
	IsCancel    uint8   `json:"isCancel"`
	ControlMode uint8   `json:"controlMode"`
	X           float32 `json:"positionX"`
	Y           float32 `json:"positionY"`
	Z           float32 `json:"positionYaw"`
}

// SetPosition sets the chassis position.
func (c *Chassis) SetPosition(m Mode, x, y, z float64) error {
	// TODO(bga): We need to implement task id handling for this.

	var controlMode uint8
	if m == ModeYawFollow {
		controlMode = 1
	}

	return c.UB().PerformActionForKey(key.KeyMainControllerChassisPosition, chassisPosition{
		TaskID:      1,
		IsCancel:    0,
		ControlMode: controlMode,
		X:           float32(x),
		Y:           float32(y),
		Z:           float32(z),
	}, nil)
}

// Move moves the robot using the given stick positions and control mode.
func (c *Chassis) Move(leftStick *controller.StickPosition,
	rightStick *controller.StickPosition, m controller.Mode) error {
	if !m.Valid() {
		return fmt.Errorf("invalid control mode: %d", m)
	}

	var leftStickEnabled uint64
	if leftStick != nil {
		leftStickEnabled = 1
	}

	var rightStickEnabled uint64
	if rightStick != nil {
		rightStickEnabled = 1
	}

	v := leftStick.InterpolatedY() |
		leftStick.InterpolatedX()<<11 |
		rightStick.InterpolatedY()<<22 |
		rightStick.InterpolatedX()<<33 |
		leftStickEnabled<<44 |
		rightStickEnabled<<45 |
		uint64(m)<<46

	return c.UB().DirectSendKeyValue(key.KeyMainControllerVirtualStick, v)
}

func (c *Chassis) control(m Mode, value uint64) error {
	if !m.Valid() {
		return fmt.Errorf("invalid mode: %d", m)
	}

	var k *key.Key
	if m == ModeYawFollow {
		k = key.KeyMainControllerChassisFollowMode
	} else {
		k = key.KeyMainControllerChassisSpeedMode
	}

	return c.UB().DirectSendKeyValue(k, value)
}
