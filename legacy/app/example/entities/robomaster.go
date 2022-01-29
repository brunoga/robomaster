package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster/legacy/app/example/components"
)

type Robomaster struct {
	*ecs.BasicEntity
	*components.Robomaster
}
