package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster/examples/robotcontrol/components"
)

type Chassis struct {
	*ecs.BasicEntity
	*components.Chassis
}
