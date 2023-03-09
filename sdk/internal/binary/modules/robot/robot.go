package robot

import (
	"fmt"

	"github.com/brunoga/robomaster/sdk/internal/binary/modules/control"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/command"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/message"
	"github.com/brunoga/robomaster/sdk/modules/robot"
	"github.com/brunoga/robomaster/sdk/support/logger"
)

type Robot struct {
	l       *logger.Logger
	control *control.Control
}

var _ robot.Robot = (*Robot)(nil)

func New(control *control.Control, l *logger.Logger) *Robot {
	return &Robot{
		l:       l,
		control: control,
	}
}

func (r *Robot) GetSDKVersion() (string, error) {
	m := message.New(
		r.control.HostByte(),
		protocol.HostToByte(8, 1),
		command.NewGetVersionRequest(),
	)

	resp, err := r.control.SendSync(m)
	if err != nil {
		return "", err
	}

	return resp.Command().(*command.GetVersionResponse).Version(), nil
}

func (r *Robot) SetMotionMode(motionMode robot.MotionMode) error {
	return fmt.Errorf("not implemented")
}
