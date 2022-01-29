package scenes

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/brunoga/robomaster/legacy/app/example/components"
	"github.com/brunoga/robomaster/legacy/app/example/entities"
	"github.com/brunoga/robomaster/legacy/app/example/systems"
)

type Robomaster struct {
	SSID     string
	Password string
	AppID    uint64
}

func (r *Robomaster) Preload() {
	// Do nothing.
}

func (r *Robomaster) Setup(u engo.Updater) {
	engo.Input.RegisterAxis("Left/Right",
		engo.AxisKeyPair{engo.KeyA, engo.KeyD})
	engo.Input.RegisterAxis("Forward/Backward",
		engo.AxisKeyPair{engo.KeyW, engo.KeyS})

	engo.Input.RegisterAxis("MouseXAxis",
		engo.NewAxisMouse(engo.AxisMouseHori))
	engo.Input.RegisterAxis("MouseYAxis",
		engo.NewAxisMouse(engo.AxisMouseVert))

	robomasterComponent, err := components.NewRobomaster(r.SSID, r.Password,
		r.AppID)
	if err != nil {
		panic(err)
	}

	err = robomasterComponent.App().Start(false)
	if err != nil {
		panic(err)
	}

	// Hack!
	// TODO(bga): Add waiting for full connection to be ablt to remove this.
	time.Sleep(5 * time.Second)

	controllerComponent, err := components.NewController(
		robomasterComponent)

	basicEntity := ecs.NewBasic()

	controllerEntity := entities.Controller{
		&basicEntity,
		controllerComponent,
	}

	// Disable cursor.
	engo.SetCursorVisibility(false)

	w, _ := u.(*ecs.World)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&systems.Video{})
	w.AddSystem(&systems.Controller{})
	w.AddSystem(&common.FPSSystem{
		Display: true,
	})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *systems.Controller:
			sys.Add(controllerEntity.BasicEntity,
				controllerComponent)
		}
	}
}

func (r *Robomaster) Type() string {
	return "Robomaster"
}
