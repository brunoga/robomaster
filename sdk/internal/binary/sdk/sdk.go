package sdk

import (
	"net"

	"github.com/brunoga/robomaster/sdk/modules/chassis"
	"github.com/brunoga/robomaster/sdk/modules/gimbal"
	"github.com/brunoga/robomaster/sdk/modules/robot"
	"github.com/brunoga/robomaster/sdk/modules/video"
	"github.com/brunoga/robomaster/sdk/support/logger"
	"github.com/brunoga/robomaster/sdk/types"
)

type SDK struct {
}

func New(connProto types.ConnectionProtocol, ip net.IP,
	l *logger.Logger) (*SDK, error) {
	return &SDK{}, nil
}

func (s *SDK) Open() error {
	return nil
}

func (s *SDK) Close() error {
	return nil
}

func (s *SDK) Robot() robot.Robot {
	return nil
}

func (s *SDK) Gimbal() gimbal.Gimbal {
	return nil
}

func (s *SDK) Chassis() chassis.Chassis {
	return nil
}

func (s *SDK) Video() video.Video {
	return nil
}
