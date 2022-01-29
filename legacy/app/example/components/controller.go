package components

import (
	"github.com/brunoga/robomaster/legacy/app/controller"
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
