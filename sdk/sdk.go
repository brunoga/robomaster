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
	Open() error

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

// New returns a new SDK instance that uses the given protocolMode and
// connectionMode to connect to the given robot ip.
func New(sdkProto types.SDKProtocol, connProto types.ConnectionProtocol,
	ip net.IP, l *logger.Logger) (SDK, error) {
	if l == nil {
		// If no logger is provided, use a logger that discards all output.
		l = logger.New(io.Discard, io.Discard, io.Discard, io.Discard)
	}

	switch sdkProto {
	case types.SDKProtocolBinary:
		return binarysdk.New(connProto, ip, l)
	case types.SDKProtocolText:
		return textsdk.New(connProto, ip, l)
	}

	return nil, fmt.Errorf("unexpected protocol mode: %v", sdkProto)
}

// NewUSB returns a new SDK instance that uses the given protocolMode and
// connectionMode to connect to the default USB (RNDIS) robot ip.
func NewUSB(sdkProto types.SDKProtocol, connProto types.ConnectionProtocol,
	l *logger.Logger) (SDK, error) {
	return New(sdkProto, connProto,
		net.ParseIP(defaultUSBConnectionIP), l)
}

// NewAP returns a new SDK instance that uses the given protocolMode and
// connectionMode to connect to the default AP robot ip.
func NewAP(sdkProto types.SDKProtocol, connProto types.ConnectionProtocol,
	l *logger.Logger) (SDK, error) {
	return New(sdkProto, connProto,
		net.ParseIP(defaultAPConnectionIP), l)
}
