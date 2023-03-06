package gimbal

import (
	"fmt"

	notification2 "github.com/brunoga/robomaster/sdk/internal/text/modules/notification"
	push2 "github.com/brunoga/robomaster/sdk/internal/text/modules/push"

	"github.com/brunoga/robomaster/sdk/internal/text/modules/control"
)

type PushAttribute int

// Supported gimbal push attributes.
const (
	// Enables gimbal attitude push notification. The events will be in the
	// format "attitude [pitch] [yaw]" where [pitch] and [yaw] are float64s.
	PushAttributeAttitude PushAttribute = iota
)

// Gimbal allows sending commands to control the robot's gimbal.
type Gimbal struct {
	control *control.Control
	push    *push2.Push
}

// New returns a new Gimbal instance associated with the given control.
func New(control *control.Control, push *push2.Push) *Gimbal {
	return &Gimbal{
		control,
		push,
	}
}

// SetSpeed sets the gimbal pitch and yaw rotation speeds in degrees/second. It
// will continue moving until it is stopped or it hits a physical limit. Returns
// a nil error on success and a non-nil error on failure.
func (g *Gimbal) SetSpeed(speed *Speed, async bool) error {
	if async {
		return g.control.SendDataExpectOkAsync(fmt.Sprintf(
			"gimbal speed p %f y %f;", speed.Pitch(),
			speed.Yaw()))
	}

	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal speed p %f y %f;", speed.Pitch(),
		speed.Yaw()))
}

// MoveRelative moves the gimbal pitch and yaw position by the given degrees
// relative to its current position and with the given speeds. Returns a nil
// error on success and a non-nil error on failure.
func (g *Gimbal) MoveRelative(position *Position, speed *Speed,
	async bool) error {
	if async {
		return g.control.SendDataExpectOkAsync(fmt.Sprintf(
			"gimbal move p %f y %f vp %f vy %f;", position.Pitch(),
			position.Yaw(), speed.Pitch(), speed.Yaw()))
	}

	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal move p %f y %f vp %f vy %f;", position.Pitch(),
		position.Yaw(), speed.Pitch(), speed.Yaw()))
}

// MoveAbsolute moves the gimbal pitch and yaw position to the given absolute
// degrees (i.e. from its origin position, not the current position) and with
// the given speeds. Returns a nil error on success and a non-nil error on
// failure.
func (g *Gimbal) MoveAbsolute(position *Position, speed *Speed,
	async bool) error {
	if async {
		return g.control.SendDataExpectOkAsync(fmt.Sprintf(
			"gimbal moveto p %f y %f vp %f vy %f;", position.Pitch(),
			position.Yaw(), speed.Pitch(), speed.Yaw()))
	}

	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal moveto p %f y %f vp %f vy %f;", position.Pitch(),
		position.Yaw(), speed.Pitch(), speed.Yaw()))
}

// Suspend puts the gimbal in power saving mode.
func (g *Gimbal) Suspend() error {
	return g.control.SendDataExpectOk("gimbal suspend;")
}

// Resume disables the gimbal's power saving mode.
func (g *Gimbal) Resume() error {
	return g.control.SendDataExpectOk("gimbal resume;")
}

// Recenter moves the gimbal to its origin position at a very low speed.
func (g *Gimbal) Recenter() error {
	return g.control.SendDataExpectOk("gimbal recenter;")
}

// GetAttitude returns the current gimbal attitude. Returns pitch attitude, yaw
// attitude and a nil error on success and a non-nil error on failure.
func (g *Gimbal) GetAttitude() (*Attitude, error) {
	data, err := g.control.SendAndReceiveData("gimbal attitude ?;")
	if err != nil {
		return nil, fmt.Errorf("error sending sdk command: %w", err)
	}

	var pitch, yaw float64
	n, err := fmt.Sscanf(data, "%f, %f", &pitch, &yaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing data: %w", err)
	}
	if n != 2 {
		return nil, fmt.Errorf("unexpected number of entries in data: %w",
			err)
	}

	return NewAttitude(pitch, yaw), nil
}

// StartPush starts listening to updates to the given attr. Returns
// a token (used to stop receiving events) and a nil error on success and a
// non-nil error on failure.
//
// TODO(bga): Add parsing of data and use a specific handler that takes parsed
//
//	attributes instead of the generic EventHandler.
func (g *Gimbal) StartPush(attr PushAttribute,
	pushHandler notification2.Handler) (int, error) {
	var token int

	var pushAttributeStr string
	switch attr {
	case PushAttributeAttitude:
		pushAttributeStr = "attitude"
	default:
		return -1, fmt.Errorf("invalid gimbal event push attribute")
	}

	token, err := g.push.StartListening("gimbal push",
		pushAttributeStr, "", pushHandler)
	if err != nil {
		return -1, fmt.Errorf(
			"error listening to gimbal push event: %w", err)
	}

	return token, nil
}

// StopPush stops sending the given event attr to the handler
// associate with the given token. Returns a nil error on success and a non-nil
// error on failure.
func (g *Gimbal) StopPush(attr PushAttribute,
	token int) error {
	var pushAttributeStr string
	switch attr {
	case PushAttributeAttitude:
		pushAttributeStr = "attitude"
	default:
		return fmt.Errorf("invalid gimbal event push attribute")
	}

	err := g.push.StopListening("gimbal push", pushAttributeStr, token)
	if err != nil {
		return fmt.Errorf(
			"error stopping listening to gimbal push notifivation: %w", err)
	}

	return nil
}
