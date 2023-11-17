package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/components"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/entities"
	"github.com/brunoga/robomaster/sdk2/module/controller"
)

type Controller struct {
	controllerEntityMap map[uint64]*entities.Controller
}

func (c *Controller) New(w *ecs.World) {
	c.controllerEntityMap = make(map[uint64]*entities.Controller)
}

func (c *Controller) Add(basicEntity *ecs.BasicEntity,
	controllerComponent *components.Controller) {
	_, ok := c.controllerEntityMap[basicEntity.ID()]
	if ok {
		return
	}

	c.controllerEntityMap[basicEntity.ID()] = &entities.Controller{
		BasicEntity: basicEntity,
		Controller:  controllerComponent,
	}
}

func (c *Controller) Remove(basicEntity ecs.BasicEntity) {
	delete(c.controllerEntityMap, basicEntity.ID())
}

func (c *Controller) Update(dt float32) {
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

	//keyPressed := currentLeftRight != 0.0 || currentForwardBackward != 0.0

	for _, controllerEntity := range c.controllerEntityMap {
		cec := controllerEntity.Controller

		//previousKeyPressed := cec.PreviousLeftRight != 0 ||
		//	cec.PreviousForwardBackward != 0
		//if !previousKeyPressed && !keyPressed {
		//	continue
		//}

		chassisY :=
			float64(((currentLeftRight * 0.3) + 1.0) / 2.0)
		chassisX :=
			float64((-(currentForwardBackward * 0.3) + 1.0) / 2.0)

		gimbalY := ((float64(-mouseYDelta) / float64(100)) + 1.0) / 2.0
		gimbalX := ((float64(mouseXDelta) / float64(100)) + 1.0) / 2.0

		cec.Controller.Move(
			&controller.StickPosition{
				X: chassisX,
				Y: chassisY,
			},
			&controller.StickPosition{
				X: gimbalX,
				Y: gimbalY,
			},
			controller.ControlModeDefault)

		cec.PreviousLeftRight = currentLeftRight
		cec.PreviousForwardBackward = currentForwardBackward
	}
}
