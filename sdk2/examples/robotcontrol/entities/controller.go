package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/components"
)

type Controller struct {
	*ecs.BasicEntity
	*components.Chassis
}
