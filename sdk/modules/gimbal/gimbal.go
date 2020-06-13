package gimbal

import (
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules"
	push2 "github.com/brunoga/robomaster/sdk/modules/push"
	"strconv"
	"strings"
)

type PushAttribute int

// Supported gimbal push attributes.
const (
	// Enables gimbal attitude push notifications. The events will be in the
	// format "attitude [pitch] [yaw]" where [pitch] and [yaw] are float64s.
	PushAttributeAttitude PushAttribute = iota
	PushAttributeInvalid
)

// Gimbal allows sending commands to control the robot's gimbal.
type Gimbal struct {
	control *modules.Control
	push    *push2.Push
}

// NewGimbal returns a new Gimbal instance associated with the given control.
func NewGimbal(control *modules.Control, push *push2.Push) *Gimbal {
	return &Gimbal{
		control,
		push,
	}
}

// SetSpeed sets the gimbal pitch and yaw rotation speeds in degrees/second. It
// will continue moving until it is stopped or it hits a physical limit. Returns
// a nil error on success and a non-nil error on failure.
func (g *Gimbal) SetSpeed(req SetSpeedRequest) error {
	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal speed p %f y %f;", req.PitchSpeedDegreesPerSecond,
		req.YawSpeedDegreesPerSecond))
}

// MoveRelative moves the gimbal pitch and yaw position by the given degrees
// relative to its current position and with the given speeds. Returns a nil
// error on success and a non-nil error on failure.
func (g *Gimbal) MoveRelative(req MoveRelativeRequest) error {
	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal move p %f y %f vp %f vy %f;", req.PitchAngleDegrees,
		req.YawAngleDegrees, req.PitchSpeedDegreesPerSecond,
		req.YawSpeedDegreesPerSecond))
}

// MoveAbsolute moves the gimbal pitch and yaw position to the given absolute
// degrees (i.e. from its origin position, not the current position) and with
// the given speeds. Returns a nil error on success and a non-nil error on
// failure.
func (g *Gimbal) MoveAbsolute(req MoveAbsoluteRequest) error {
	return g.control.SendDataExpectOk(fmt.Sprintf(
		"gimbal moveto p %f y %f vp %f vy %f;", req.PitchAngleDegrees,
		req.YawAngleDegrees, req.PitchSpeedDegreesPerSecond,
		req.YawSpeedDegreesPerSecond))
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
func (g *Gimbal) GetAttitude() (GetAttitudeResponse, error) {
	data, err := g.control.SendAndReceiveData("gimbal attitude ?;")
	if err != nil {
		return GetAttitudeResponse{}, fmt.Errorf(
			"error sending sdk command: %w", err)
	}

	fields := strings.Fields(data)
	if len(fields) != 2 {
		return GetAttitudeResponse{}, fmt.Errorf("unexpected response received")
	}

	pitch, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return GetAttitudeResponse{}, fmt.Errorf(
			"error decoding pitch angle: %w", err)
	}

	yaw, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return GetAttitudeResponse{}, fmt.Errorf(
			"error decoding yaw angle: %w", err)
	}

	return GetAttitudeResponse{pitch, yaw}, nil
}

// StartPush starts listening to updates to the given attr. Returns
// a token (used to stop receiving events) and a nil error on success and a
// non-nil error on failure.
//
// TODO(bga): Add parsing of data and use a specific handler that takes parsed
//  attributes instead of the generic EventHandler.
func (g *Gimbal) StartPush(req StartPushRequest) (StartPushResponse, error) {
	var token int

	switch req.PushAttribute {
	case PushAttributeAttitude:
		var err error
		token, err = g.push.StartListening("gimbal push",
			"attitude on", req.PushHandler)
		if err != nil {
			return StartPushResponse{
					-1,
				}, fmt.Errorf(
					"error listening to gimbal push event: %w", err)
		}
	default:
		return StartPushResponse{
			-1,
		}, fmt.Errorf("invalid gimbal event push attribute")
	}

	return StartPushResponse{
		token,
	}, nil
}

// StopPush stops sending the given event attr to the handler
// associate with the given token. Returns a nil error on success and a non-nil
// error on failure.
func (g *Gimbal) StopPush(req StopPushRequest) error {
	switch req.PushAttribute {
	case PushAttributeAttitude:
		err := g.push.StopListening("gimbal push",
			"attitude off", req.Token)
		if err != nil {
			return fmt.Errorf(
				"error stopping listening to gimbal push notifivation: %w", err)
		}
	default:
		return fmt.Errorf("invalid gimbal push notification attribute")
	}

	return nil
}
