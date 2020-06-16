package modules

import (
	"fmt"
)

type GimbalPushAttribute int

// Supported gimbal push attributes.
const (
	// Enables gimbal attitude push notifications. The events will be in the
	// format "attitude [pitch] [yaw]" where [pitch] and [yaw] are float64s.
	GimbalPushAttributeAttitude GimbalPushAttribute = iota
	GimbalPushAttributeInvalid
)

// Gimbal allows sending commands to control the robot's gimbal.
type Gimbal struct {
	control *Control
	push    *Push
}

// NewGimbal returns a new Gimbal instance associated with the given control.
func NewGimbal(control *Control, push *Push) *Gimbal {
	return &Gimbal{
		control,
		push,
	}
}

// SetSpeed sets the gimbal pitch and yaw rotation speeds in degrees/second. It
// will continue moving until it is stopped or it hits a physical limit. Returns
// a nil error on success and a non-nil error on failure.
func (g *Gimbal) SetSpeed(pitchDegreesPerSecond,
	yawDegreesPerSecond float64) error {
	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal speed p %f y %f;", pitchDegreesPerSecond,
		yawDegreesPerSecond))
}

// MoveRelative moves the gimbal pitch and yaw position by the given degrees
// relative to its current position and with the given speeds. Returns a nil
// error on success and a non-nil error on failure.
func (g *Gimbal) MoveRelative(pitchDegreesPos, yawDegreesPos,
	pitchDegreesPerSecond, yawDegreesPerSecond float64) error {
	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal move p %f y %f vp %f vy %f;", pitchDegreesPos,
		yawDegreesPos, pitchDegreesPerSecond, yawDegreesPerSecond))
}

// MoveAbsolute moves the gimbal pitch and yaw position to the given absolute
// degrees (i.e. from its origin position, not the current position) and with
// the given speeds. Returns a nil error on success and a non-nil error on
// failure.
func (g *Gimbal) MoveAbsolute(pitchDegreesPos, yawDegreesPos,
	pitchDegreesPerSecond, yawDegreesPerSecond float64) error {
	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal moveto p %f y %f vp %f vy %f;", pitchDegreesPos,
		yawDegreesPos, pitchDegreesPerSecond, yawDegreesPerSecond))
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
func (g *Gimbal) GetAttitude() (float64, float64, error) {
	data, err := g.control.SendAndReceiveData("gimbal attitude ?;")
	if err != nil {
		return 0, 0, fmt.Errorf("error sending sdk command: %w", err)
	}

	var pitch, yaw float64
	n, err := fmt.Sscanf(data, "%f, %f", &pitch, &yaw)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing data: %w", err)
	}
	if n != 2 {
		return 0, 0, fmt.Errorf("unexpected number of entries in data: %w",
			err)
	}

	return pitch, yaw, nil
}

// StartGimbalPush starts listening to updates to the given attr. Returns
// a token (used to stop receiving events) and a nil error on success and a
// non-nil error on failure.
//
// TODO(bga): Add parsing of data and use a specific handler that takes parsed
//  attributes instead of the generic EventHandler.
func (g *Gimbal) StartGimbalPush(attr GimbalPushAttribute,
	pushHandler PushHandler) (int, error) {
	var token int

	switch attr {
	case GimbalPushAttributeAttitude:
		var err error
		token, err = g.push.StartListening("gimbal push",
			"attitude", "", pushHandler)
		if err != nil {
			return -1, fmt.Errorf(
				"error listening to gimbal push event: %w", err)
		}
	default:
		return -1, fmt.Errorf("invalid gimbal event push attribute")
	}

	return token, nil
}

// StopGimbalPush stops sending the given event attr to the handler
// associate with the given token. Returns a nil error on success and a non-nil
// error on failure.
func (g *Gimbal) StopGimbalPush(attr GimbalPushAttribute,
	token int) error {
	switch attr {
	case GimbalPushAttributeAttitude:
		err := g.push.StopListening("gimbal push",
			"attitude", token)
		if err != nil {
			return fmt.Errorf(
				"error stopping listening to gimbal push notifivation: %w", err)
		}
	default:
		return fmt.Errorf("invalid gimbal push notification attribute")
	}

	return nil
}
