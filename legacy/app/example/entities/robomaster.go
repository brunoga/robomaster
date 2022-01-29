package entities

import (
	"git.bug-br.org.br/bga/robomasters1/app/example/components"
	"github.com/EngoEngine/ecs"
)

type Robomaster struct {
	*ecs.BasicEntity
	*components.Robomaster
}
