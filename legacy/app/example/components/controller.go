package components

import (
	"git.bug-br.org.br/bga/robomasters1/app/controller"
)

type Controller struct {
	Controller *controller.Controller

	PreviousLeftRight       float32
	PreviousForwardBackward float32
}

func NewController(r *Robomaster) (*Controller, error) {
	return &Controller{
		controller.New(
			r.App().CommandController(),
		),
		0.0,
		0.0,
	}, nil
}
