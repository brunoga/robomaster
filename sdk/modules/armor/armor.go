package armor

import (
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules"
	"github.com/brunoga/robomaster/sdk/modules/notification"
	"strconv"
)

type Armor struct {
	control *modules.Control
	event   *notification.Event
}

func New(control *modules.Control, event *notification.Event) *Armor {
	return &Armor{
		control,
		event,
	}
}

func (a *Armor) SetSensitivity(sensitivity int) error {
	return a.control.SendDataExpectOk(fmt.Sprintf(
		"armor sensitivity %d;", sensitivity))
}

func (a *Armor) GetSensitivity() (int, error) {
	data, err := a.control.SendAndReceiveData("armor sensitivity ?;")
	if err != nil {
		return -1, fmt.Errorf("error sending sdk command: %w", err)
	}

	sensitivity, err := strconv.Atoi(data)
	if err != nil {
		return -1, fmt.Errorf("error parsing data: %w", err)
	}

	return sensitivity, nil
}
