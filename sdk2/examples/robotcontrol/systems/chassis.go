package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/components"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/entities"
	"github.com/brunoga/robomaster/sdk2/module/chassis/controller"
)

type Chassis struct {
	controllerEntityMap map[uint64]*entities.Controller
}

func (c *Chassis) New(w *ecs.World) {
	c.controllerEntityMap = make(map[uint64]*entities.Controller)
}

func (c *Chassis) Add(basicEntity *ecs.BasicEntity,
	controllerComponent *components.Chassis) {
	_, ok := c.controllerEntityMap[basicEntity.ID()]
	if ok {
		return
	}

	c.controllerEntityMap[basicEntity.ID()] = &entities.Controller{
		BasicEntity: basicEntity,
		Chassis:     controllerComponent,
	}
}

func (c *Chassis) Remove(basicEntity ecs.BasicEntity) {
	delete(c.controllerEntityMap, basicEntity.ID())
}

func (c *Chassis) Update(dt float32) {
	if btn := engo.Input.Button("exit"); btn.JustPressed() {
		engo.Exit()
	}

	currentLeftRight := engo.Input.Axis("Left/Right").Value()
	currentForwardBackward := engo.Input.Axis("Forward/Backward").Value()

	mouseXDelta := engo.Input.Axis("MouseXAxis").Value()
	if mouseXDelta > 100 {
		mouseXDelta = 100
	} else if mouseXDelta < -100 {
		mouseXDelta = -100
	}

	mouseYDelta := engo.Input.Axis("MouseYAxis").Value()
	if mouseYDelta > 100 {
		mouseYDelta = 100
	} else if mouseYDelta < -100 {
		mouseYDelta = -100
	}

	for _, controllerEntity := range c.controllerEntityMap {
		cec := controllerEntity.Chassis

		chassisStickPosition := &controller.StickPosition{
			X: float64(currentLeftRight),
			Y: float64(currentForwardBackward),
		}

		gimbalStickPosition := &controller.StickPosition{
			X: float64(mouseXDelta) / float64(100),
			Y: float64(mouseYDelta) / float64(100),
		}

		cec.Chassis.Move(chassisStickPosition, gimbalStickPosition,
			controller.ControlModeDefault)
	}
}
