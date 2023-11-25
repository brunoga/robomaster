package scenes

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/components"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/entities"
	"github.com/brunoga/robomaster/sdk2/examples/robotcontrol/systems"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Robomaster struct {
	Client *sdk2.Client
}

func (r *Robomaster) Preload() {
	// Do nothing.
}

func (r *Robomaster) Setup(u engo.Updater) {
	engo.Input.RegisterAxis("Left/Right",
		engo.AxisKeyPair{Min: engo.KeyA, Max: engo.KeyD})
	engo.Input.RegisterAxis("Forward/Backward",
		engo.AxisKeyPair{Min: engo.KeyW, Max: engo.KeyS})

	engo.Input.RegisterAxis("MouseXAxis",
		engo.NewAxisMouse(engo.AxisMouseHori))
	engo.Input.RegisterAxis("MouseYAxis",
		engo.NewAxisMouse(engo.AxisMouseVert))

	engo.Input.RegisterButton("exit", engo.KeyEscape)

	controllerComponent := &components.Chassis{
		Chassis: r.Client.Chassis(),
	}

	basicEntity := ecs.NewBasic()

	controllerEntity := entities.Controller{
		BasicEntity: &basicEntity,
		Chassis:     controllerComponent,
	}

	// Disable cursor.
	if engo.CurrentBackEnd == engo.BackEndGLFW ||
		engo.CurrentBackEnd == engo.BackEndVulkan {
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode,
			glfw.CursorDisabled)
	} else {
		panic("Backend does not seem to support mouse capture.")
	}

	w, _ := u.(*ecs.World)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&systems.Video{
		Camera: r.Client.Camera(),
	})

	w.AddSystem(&systems.Chassis{})
	w.AddSystem(&common.FPSSystem{
		Display: true,
	})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *systems.Chassis:
			sys.Add(controllerEntity.BasicEntity,
				controllerComponent)
		}
	}
}

func (r *Robomaster) Type() string {
	return "Robomaster"
}
