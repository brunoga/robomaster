package scenes

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/text/examples/robotcontrol/systems"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type RobotControl struct {
	Client        *sdk.Client
	MirrorClients []*sdk.Client
}

func (r RobotControl) Preload() {}

func (r RobotControl) Setup(updater engo.Updater) {
	if engo.CurrentBackEnd == engo.BackEndGLFW {
		// Use glfw directly to set the mouse input to what we need (basically
		// infinite movement). Otherwise movement would stop whenever the limits
		// of the window would be reached.
		glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else {
		// Better than nothing, I guess.
		engo.SetCursorVisibility(false)
	}

	engo.Input.RegisterAxis("Left/Right",
		engo.AxisKeyPair{Min: engo.KeyA, Max: engo.KeyD})
	engo.Input.RegisterAxis("Forward/Backward",
		engo.AxisKeyPair{Min: engo.KeyW, Max: engo.KeyS})

	engo.Input.RegisterAxis("MouseXAxis",
		engo.NewAxisMouse(engo.AxisMouseHori))
	engo.Input.RegisterAxis("MouseYAxis",
		engo.NewAxisMouse(engo.AxisMouseVert))

	w, _ := updater.(*ecs.World)
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(systems.NewVideo(r.Client))
	w.AddSystem(systems.NewGimbal(r.Client, r.MirrorClients))
	w.AddSystem(systems.NewChassis(r.Client, r.MirrorClients))
	w.AddSystem(systems.NewBlaster(r.Client, r.MirrorClients))
	w.AddSystem(&common.FPSSystem{
		Display: true,
	})
}

func (r RobotControl) Type() string {
	return "RobotControl"
}
