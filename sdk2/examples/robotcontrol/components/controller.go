package components

import (
	"github.com/brunoga/robomaster/sdk2/module/controller"
)

type Controller struct {
	Controller *controller.Controller

	PreviousLeftRight       float32
	PreviousForwardBackward float32
}

func NewController(r *Robomaster) (*Controller, error) {
	return &Controller{
		r.Client().Controller(),
		0.0,
		0.0,
	}, nil
}
