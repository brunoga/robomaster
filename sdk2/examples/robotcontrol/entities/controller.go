package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/components"
)

type Chassis struct {
	*ecs.BasicEntity
	*components.Chassis
}
