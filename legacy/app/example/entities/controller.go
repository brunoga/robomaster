package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/brunoga/robomaster/legacy/app/example/components"
)

type Controller struct {
	*ecs.BasicEntity
	*components.Controller
}
