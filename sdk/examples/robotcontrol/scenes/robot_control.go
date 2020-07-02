package scenes

import (
	"flag"

	"github.com/brunoga/robomaster/sdk/modules/robot"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/examples/robotcontrol/systems"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var mirrors = flag.Uint("mirrors", 0, "number of mirror robots")

type RobotControl struct{}

func (r RobotControl) Preload() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

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

	client, err := sdk.NewClient(nil)
	if err != nil {
		panic(err)
	}

	err = client.Open()
	if err != nil {
		panic(err)
	}

	client.RobotModule().SetMotionMode(robot.MotionModeGimbalLead)

	var mirrorClients []*sdk.Client
	if *mirrors > 0 {
		mirrorClients = make([]*sdk.Client, *mirrors)
		for i := 0; i < len(mirrorClients); i++ {
			mirrorClients[i], err = sdk.NewClient(nil)
			if err != nil {
				panic(err)
			}

			err = mirrorClients[i].Open()
			if err != nil {
				panic(err)
			}

			mirrorClients[i].RobotModule().SetMotionMode(robot.MotionModeGimbalLead)
		}
	}

	w, _ := updater.(*ecs.World)
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(systems.NewVideo(client))
	w.AddSystem(systems.NewGimbal(client, mirrorClients))
	w.AddSystem(systems.NewChassis(client, mirrorClients))
	w.AddSystem(systems.NewBlaster(client, mirrorClients))
	w.AddSystem(&common.FPSSystem{
		Display: true,
	})
}

func (r RobotControl) Type() string {
	return "RobotControl"
}
