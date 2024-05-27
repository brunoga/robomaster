package chassis

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
	"github.com/brunoga/robomaster/unitybridge/unity/task"
)

// Chassis allows controlling the robot chassis. It also works as the robot main
// controller interface.
type Chassis struct {
	*internal.BaseModule
}

var _ module.Module = (*Chassis)(nil)

// New creates a new Chassis instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection, rm *robot.Robot) (*Chassis, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("chassis_module")

	c := &Chassis{}

	c.BaseModule = internal.NewBaseModule(ub, l, "Chassis", nil, func(r *result.Result) {
		if !r.Succeeded() {
			c.Logger().Error("Connection: Unsuccessfull result.", "result", r)
			return
		}

		connectedValue, ok := r.Value().(*value.Bool)
		if !ok {
			c.Logger().Error("Connection: Unexpected value.", "value", r.Value())
			return
		}

		if connectedValue.Value {
			c.Logger().Debug("Connected.")
		} else {
			c.Logger().Debug("Disconnected.")
		}
	})

	return c, nil
}

// SetMode sets the chassis mode for the robot.
func (c *Chassis) SetMode(m Mode) error {
	if !m.Valid() {
		return fmt.Errorf("invalid mode: %d", m)
	}

	// TODO(bga): Figure out this value.
	value := uint64(1) | uint64(140) | uint64(17920) | uint64(235929600)

	defer func() {
		// Stop movement after 0.3 seconds. Most likelly to give enough
		// time for the mode setting to stick.
		time.Sleep(333 * time.Millisecond)
		c.StopMovement(m)
	}()

	return c.control(m, value)
}

// StopMovement stops the chassis movement.
func (c *Chassis) StopMovement(m Mode) error {
	if !m.Valid() {
		return fmt.Errorf("invalid mode: %d", m)
	}

	// TODO(bga): Figure out this value.
	value := uint64(0) | uint64(140) | uint64(17920) | uint64(235929600)

	return c.control(m, value)
}

// SetSpeed sets the chassis speed. Limits are [-3.5, 3.5] (m/s) for x and y and
// [-360, 360] (degrees/s) for z.
func (c *Chassis) SetSpeed(m Mode, x, y, z float64) error {
	if x > 3.5 || x < -3.5 || y > 3.5 || y < -3.5 || z > 360 || z < -360 {
		return fmt.Errorf("invalid speed values: x=%f, y=%f, z=%f", x, y, z)
	}

	xComponent := (int64(x*10) + 35) << 2
	yComponent := (int64(y*10) + 35) << 9
	zComponent := (int64(z*10) + 3600) << 16

	value := uint64(1 | xComponent | yComponent | zComponent)

	return c.control(m, value)
}

// SetPosition sets the chassis position.
func (c *Chassis) SetPosition(m Mode, x, y, z float64) error {
	// TODO(bga): We need to implement task id handling for this.

	var controlMode uint8
	if m == ModeYawFollow {
		controlMode = 1
	}

	return c.UB().PerformActionForKeySync(key.KeyMainControllerChassisPosition, &value.ChassisPosition{
		TaskType:    task.TypeChassisPosition,
		IsCancel:    0,
		ControlMode: controlMode,
		X:           float32(x),
		Y:           float32(y),
		Z:           float32(z),
	})
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
