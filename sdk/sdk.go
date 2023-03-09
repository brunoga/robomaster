package sdk

import (
	"fmt"
	"io"
	"net"

	"github.com/brunoga/robomaster/sdk/modules/chassis"
	"github.com/brunoga/robomaster/sdk/modules/gimbal"
	"github.com/brunoga/robomaster/sdk/modules/robot"
	"github.com/brunoga/robomaster/sdk/modules/video"
	"github.com/brunoga/robomaster/sdk/support/logger"
	"github.com/brunoga/robomaster/sdk/types"

	binarysdk "github.com/brunoga/robomaster/sdk/internal/binary/sdk"
	textsdk "github.com/brunoga/robomaster/sdk/internal/text/sdk"
)

const (
	defaultUSBConnectionIP = "192.168.42.2"
	defaultAPConnectionIP  = "192.168.2.1"
)

type SDK interface {
	// Open opens the SDK connection to the robot.
	Open(types.ConnectionMode, types.ConnectionProtocol, net.IP) error

	// Close closes the SDK connection to the robot.
	Close() error

	// Robot returns the robot module. This allows setting/getting general robot
	// operating parameters.
	Robot() robot.Robot

	// Gimbal returns the gimbal module. This allows controlling the gimbal.
	Gimbal() gimbal.Gimbal

	// Chassis returns the chassis module. This allows controlling the chassis.
	Chassis() chassis.Chassis

	// Video returns the video module. This allows getting video frames from the
	// robot's camera.
	Video() video.Video
}

// New creates a new SDK instance using the given SDK protocol and logger.
func New(sdkProto types.SDKProtocol, l *logger.Logger) (SDK, error) {
	if l == nil {
		// If no logger is provided, use a logger that discards all output.
		l = logger.New(io.Discard, io.Discard, io.Discard, io.Discard)
	}

	switch sdkProto {
	case types.SDKProtocolBinary:
		return binarysdk.New(l)
	case types.SDKProtocolText:
		return textsdk.New(l)
	}

	return nil, fmt.Errorf("unexpected protocol mode: %v", sdkProto)
}
