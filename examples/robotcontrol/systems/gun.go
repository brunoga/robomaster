package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/examples/robotcontrol/components"
	"github.com/brunoga/robomaster/examples/robotcontrol/entities"
	"github.com/brunoga/robomaster/module/gun"
)

type Gun struct {
	gunEntityMap map[uint64]*entities.Gun
}

func (c *Gun) New(w *ecs.World) {
	c.gunEntityMap = make(map[uint64]*entities.Gun)
}

func (c *Gun) Add(basicEntity *ecs.BasicEntity,
	gunComponent *components.Gun) {
	_, ok := c.gunEntityMap[basicEntity.ID()]
	if ok {
		return
	}

	c.gunEntityMap[basicEntity.ID()] = &entities.Gun{
		BasicEntity: basicEntity,
		Gun:         gunComponent,
	}
}

func (c *Gun) Remove(basicEntity ecs.BasicEntity) {
	delete(c.gunEntityMap, basicEntity.ID())
}

func (c *Gun) Update(dt float32) {
	fire := false
	if engo.Input.Mouse.Action == engo.Press &&
		engo.Input.Mouse.Button == engo.MouseButtonLeft {
		fire = true
	} else if engo.Input.Mouse.Action == engo.Release &&
		engo.Input.Mouse.Button == engo.MouseButtonLeft {
		fire = false
	}

	for _, gunEntity := range c.gunEntityMap {
		gug := gunEntity.Gun
		if fire {
			gug.Gun.Fire(gun.TypeInfrared)
		}
	}
}
