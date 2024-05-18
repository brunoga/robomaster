package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type Video struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}
