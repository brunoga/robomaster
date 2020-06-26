package chassis

import (
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules"
	"github.com/brunoga/robomaster/sdk/modules/notification"
)

// PushAttribute represents chassis attributes that can be monitored through
// push notification.
type PushAttribute uint8

// Supported chassis push attributes.
const (
	PushAttributePosition PushAttribute = iota
	PushAttributeAttitude
	PushAttributeStatus
)

// Chassis allows sending commands to control the robot's chassis.
type Chassis struct {
	control *modules.Control
	push    *notification.Push
}

// New returns a new Chassis instance associated with the given control.
func New(control *modules.Control, push *notification.Push) *Chassis {
	return &Chassis{
		control,
		push,
	}
}

// SetSpeed sets the chassis speed to the given speed. Returns a nil error on
// success and a non-nil error on failure.
func (c *Chassis) SetSpeed(speed *Speed) error {
	return c.control.SendDataExpectOk(fmt.Sprintf(
		"chassis speed x %f y %f z %f;", speed.X(), speed.Y(), speed.Z()))
}

// SetWheelSpeed sets the chassis individual wheels speed to the given
// wheelSpeed. Returns a nil error on success and a non-nil error on
// failure.
func (c *Chassis) SetWheelSpeed(wheelSpeed *WheelSpeed) error {
	return c.control.SendDataExpectOk(fmt.Sprintf(
		"chassis wheel w1 %f w2 %f w3 %f w4 %f;", wheelSpeed.W1(),
		wheelSpeed.W2(), wheelSpeed.W3(), wheelSpeed.W4()))
}

// MoveRelative moves the Chassis tro the given position relative to the current
// one at the given speed. Returns a nil error on success and a non-nil error on
// failure.
func (c *Chassis) MoveRelative(position *Position, speed *Speed) error {
	return c.control.SendDataExpectOk(fmt.Sprintf(
		"chassis move x %f y %f z %f vxy %f vz %f;", position.X(),
		position.Y(), position.Z(), speed.X(), speed.Z()))
}

// GetSpeed returns the current chassis and wheel speed and a nil error on
// success and a non-nil error on failure.
func (c *Chassis) GetSpeed() (*Speed, *WheelSpeed, error) {
	data, err := c.control.SendAndReceiveData("chassis speed ?;")
	if err != nil {
		return nil, nil, fmt.Errorf("error sending sdk command: %w",
			err)
	}

	speed, err := NewSpeedFromData(data)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing speed: %w", err)
	}

	wheelSpeed, err := NewWheelSpeedFromData(findWheelSpeedInData(data))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing wheel speed: %w", err)
	}

	return speed, wheelSpeed, nil
}

// GetPosition returns the current chassis position relative to the position at
// power on and a nil error on success and a non-nil error on failure.
func (c *Chassis) GetPosition() (*Position, error) {
	data, err := c.control.SendAndReceiveData("chassis position ?;")
	if err != nil {
		return nil, fmt.Errorf("error sending sdk command: %w",
			err)
	}

	position, err := NewPositionFromData(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing position: %w", err)
	}

	return position, nil
}

// GetAttitude returns the current chassis attitude relative to the position at
// power on and a nil error on success and a non-nil error on failure.
func (c *Chassis) GetAttitude() (*Attitude, error) {
	data, err := c.control.SendAndReceiveData("chassis attitude ?;")
	if err != nil {
		return nil, fmt.Errorf("error sending sdk command: %w",
			err)
	}

	attitude, err := NewAttitudeFromData(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing attitude: %w", err)
	}

	return attitude, nil
}

// GetStatus returns the current chassis status and a nil error on success and
// a non-nil error on failure.
func (c *Chassis) GetStatus() (*Status, error) {
	data, err := c.control.SendAndReceiveData("chassis status ?;")
	if err != nil {
		return nil, fmt.Errorf("error sending sdk command: %w", err)
	}

	status, err := NewStatusFromData(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing status: %w", err)
	}

	return status, nil
}

// StartPush starts the event push for the given pushType and PushAttribute.
// Updates will be sent to the given pushHandler at the specified frequency.
// Returns a token (used to stop pushes for the given eventHandler) and a nil
// error on success and a non-nil error on failure.
func (c *Chassis) StartPush(pushAttribute PushAttribute,
	pushHandler notification.Handler, frequency int) (int, error) {
	var pushAttributeStr string
	switch pushAttribute {
	case PushAttributePosition:
		pushAttributeStr = "position"
	case PushAttributeAttitude:
		pushAttributeStr = "attitude"
	case PushAttributeStatus:
		pushAttributeStr = "status"
	default:
		return -1, fmt.Errorf("invalid chassis push attribute")
	}

	pushParameters := fmt.Sprintf("freq %d", frequency)

	token, err := c.push.StartListening("chassis push",
		pushAttributeStr, pushParameters, pushHandler)
	if err != nil {
		return -1, fmt.Errorf(
			"error starting to listen to gimbal push event: %w", err)
	}

	return token, nil
}

// StopPush stops push event to the push handler represented by the given
// pushAttribute and token pair. Returns a nil error on success and a non-nil
// error on failure.
func (c *Chassis) StopPush(pushAttribute PushAttribute, token int) error {
	var pushAttributeStr string
	switch pushAttribute {
	case PushAttributePosition:
		pushAttributeStr = "position"
	case PushAttributeAttitude:
		pushAttributeStr = "attitude"
	case PushAttributeStatus:
		pushAttributeStr = "status"
	default:
		return fmt.Errorf("invalid chassis push attribute")
	}

	err := c.push.StopListening("chassis push", pushAttributeStr, token)
	if err != nil {
		return fmt.Errorf(
			"error starting to listen to gimbal push event: %w", err)
	}

	return nil
}

func findWheelSpeedInData(data string) string {
	count := 0
	for i := 0; i < len(data); i++ {
		if data[i] == ' ' {
			count++
		}

		if count == 3 {
			return data[count+1:]
		}
	}

	return ""
}
