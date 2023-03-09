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

func (r *Robot) GetProductVersion() (string, error) {
	m := message.New(
		r.control.HostByte(),
		protocol.HostToByte(8, 1),
		command.NewGetProductVersionRequest(),
	)

	resp, err := r.control.SendSync(m)
	if err != nil {
		return "", err
	}

	return resp.Command().(*command.GetProductVersionResponse).Version(), nil
}

func (r *Robot) SetMode(robotMode robot.Mode) error {
	cmd := command.NewSetRobotModeRequest()
	cmd.SetMode(robotMode)

	m := message.New(
		r.control.HostByte(),
		protocol.HostToByte(9, 0),
		cmd,
	)

	resp, err := r.control.SendSync(m)
	if err != nil {
		return err
	}

	if !resp.Command().(command.Response).Ok() {
		return fmt.Errorf("client set robot mode: not ok")
	}

	return nil
}
