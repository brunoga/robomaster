package components

import (
	"github.com/brunoga/robomaster/sdk2/module/controller"
)

type Controller struct {
	Controller *controller.Controller
}

func NewController(r *Robomaster) (*Controller, error) {
	return &Controller{
		r.Client().Controller(),
	}, nil
}
