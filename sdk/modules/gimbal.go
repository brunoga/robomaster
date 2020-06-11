package modules

import (
	"fmt"
	"strconv"
	"strings"
)

type GimbalEventPushAttribute int

// Supported gimbal push attributes.
const (
	// Enables gimbal attitude push events. The events will be in the format
	// "attitude [pitch] [yaw]" where [pitch] and [yaw] are float64s.
	GimbalEventPushAttributeAttitude GimbalEventPushAttribute = iota
	GimbalEventPushAttributeInvalid
)

// Gimbal allows sending commands to control the robot's gimbal.
type Gimbal struct {
	control *Control
	event   *Event
}

// NewGimbal returns a new Gimbal instance associated with the given control.
func NewGimbal(control *Control, event *Event) *Gimbal {
	return &Gimbal{
		control,
		event,
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
func (g *Gimbal) GetAttitude() (int, int, error) {
	data, err := g.control.SendAndReceiveData("gimbal attitude ?;")
	if err != nil {
		return 0, 0, fmt.Errorf("error sending sdk command: %w", err)
	}

	fields := strings.Fields(data)
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected response received")
	}

	pitch, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("error decoding pitch angle: %w", err)
	}

	yaw, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("error decoding yaw angle: %w", err)
	}

	return pitch, yaw, nil
}

// StartGimbalEventPush starts listening to updates to the given attr. Returns
// a token (used to stop receiving events) and a nil error on success and a
// non-nil error on failure.
//
// TODO(bga): Add parsing of data and use a specific handler that takes parsed
//  attributes instead of the generic EventHandler.
func (g *Gimbal) StartGimbalEventPush(attr GimbalEventPushAttribute,
	eventHandler EventHandler) (int, error) {
	var token int

	switch attr {
	case GimbalEventPushAttributeAttitude:
		var err error
		token, err = g.event.StartListening("gimbal push",
			"attitude on", eventHandler)
		if err != nil {
			return -1, fmt.Errorf(
				"error listening to gimbal push event: %w", err)
		}
	default:
		return -1, fmt.Errorf("invalid gimbal event push attribute")
	}

	return token, nil
}

// StopGimbalEventPush stops sending the given event attr to the handler
// associate with the given token. Returns a nil error on success and a non-nil
// error on failure.
func (g *Gimbal) StopGimbalEventPush(attr GimbalEventPushAttribute,
	token int) error {
	switch attr {
	case GimbalEventPushAttributeAttitude:
		err := g.event.StopListening("gimbal push",
			"attitude off", token)
		if err != nil {
			return fmt.Errorf(
				"error stopping listening to gimbal push event: %w", err)
		}
	default:
		return fmt.Errorf("invalid gimbal event push attribute")
	}

	return nil
}
