package sdk

import (
	"net"

	"github.com/brunoga/robomaster/sdk/modules/chassis"
	"github.com/brunoga/robomaster/sdk/modules/gimbal"
	"github.com/brunoga/robomaster/sdk/modules/robot"
	"github.com/brunoga/robomaster/sdk/modules/video"
	"github.com/brunoga/robomaster/sdk/support/logger"
	"github.com/brunoga/robomaster/sdk/types"

	binarycontrol "github.com/brunoga/robomaster/sdk/internal/binary/modules/control"
	binaryfinder "github.com/brunoga/robomaster/sdk/internal/binary/modules/finder"
	binaryrobot "github.com/brunoga/robomaster/sdk/internal/binary/modules/robot"
)

type SDK struct {
	l       *logger.Logger
	finder  *binaryfinder.Finder
	control *binarycontrol.Control
	robot   *binaryrobot.Robot
}

func New(l *logger.Logger) (*SDK, error) {
	finder := binaryfinder.New(l)
	control, err := binarycontrol.New(9, 6, finder, l)
	if err != nil {
		return nil, err
	}
	robot := binaryrobot.New(control, l)

	return &SDK{
		l:       l,
		finder:  finder,
		control: control,
		robot:   robot,
	}, nil
}

func (s *SDK) Open(connMode types.ConnectionMode,
	connProto types.ConnectionProtocol, ip net.IP) error {
	return s.control.Open(connMode, connProto, ip)
}

func (s *SDK) Close() error {
	return s.control.Close()
}

func (s *SDK) Robot() robot.Robot {
	return s.robot
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
