package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster/examples/robotcontrol/components"
)

type Gun struct {
	*ecs.BasicEntity
	*components.Gun
}
