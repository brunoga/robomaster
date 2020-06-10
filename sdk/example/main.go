package main

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/modules"
	"github.com/brunoga/robomaster/sdk/support"
	"github.com/brunoga/robomaster/sdk/support/pid"
	"gocv.io/x/gocv"
)

type exampleVideoHandler struct {
	window       *gocv.Window
	tracker      *support.ColorObjectTracker
	gimbalModule *modules.Gimbal
	pidPitch     pid.Controller
	pidYaw       pid.Controller
	quitChan     chan struct{}
}

func newExampleVideoHandler(gimbalModule *modules.Gimbal) *exampleVideoHandler {
	window := gocv.NewWindow("Robomaster")
	window.ResizeWindow(sdk.CameraHorizontalResolutionPoints/2,
		sdk.CameraVerticalResolutionPoints/2)

	return &exampleVideoHandler{
		window,
		support.NewColorObjectTracker(35, 219, 90, 119, 255,
			255, 10),
		gimbalModule,
		pid.NewPIDController(250, 10, 30, -400, 400),
		pid.NewPIDController(150, 10, 20, -400, 400),
		make(chan struct{}),
	}
}

func (e *exampleVideoHandler) QuitChan() <-chan struct{} {
	return e.quitChan
}

func (e *exampleVideoHandler) HandleFrame(frame *gocv.Mat, wg *sync.WaitGroup) {
	defer wg.Done()

	x, y, radius, err := e.tracker.FindLargestObject(frame)
	if err == nil {
		// Found something. Draw a circle around it.
		gocv.Circle(frame, image.Point{X: int(x), Y: int(y)}, int(radius),
			color.RGBA{R: 0, G: 255, B: 255, A: 255}, 2)

		// Get errors in the x and y axis normalized to [-0.5, 0.5]
		errX := float64(x-(sdk.CameraHorizontalResolutionPoints/2)) /
			sdk.CameraHorizontalResolutionPoints
		errY := float64((sdk.CameraVerticalResolutionPoints/2)-y) /
			sdk.CameraVerticalResolutionPoints

		// If there is some error (object in not in center of image), move the
		// gimbal to minimize it.
		if math.Abs(errX) > 0.0 || math.Abs(errY) > 0.0 {
			// Move the gimbal with a speed determined by the pitch and yaw PID
			// controllers.
			err = e.gimbalModule.SetSpeed(e.pidPitch.Output(errY),
				e.pidYaw.Output(errX))
			if err != nil {
				// TODO(bga): Log this.
			}
		}
	}

	// Show modified frame and wait for an event.
	e.window.IMShow(*frame)
	e.window.WaitKey(1)

	if e.window.GetWindowProperty(gocv.WindowPropertyAspectRatio) == -1.0 {
		// Window closed. Notify listeners.
		close(e.quitChan)
	}
}

func main() {
	client := sdk.NewClient(nil)

	err := client.Open()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Obtain references to the modules we are interested in.
	robotModule := client.RobotModule()
	gimbalModule := client.GimbalModule()
	videoModule := client.VideoModule()

	// Control gimbal/chassis independently.
	err = robotModule.SetMotionMode(modules.RobotMotionModeFree)
	if err != nil {
		panic(err)
	}

	// Reset gimbal.
	err = gimbalModule.Recenter()
	if err != nil {
		panic(err)
	}

	videoHandler := newExampleVideoHandler(gimbalModule)

	token, err := videoModule.StartStream(videoHandler.HandleFrame)
	if err != nil {
		panic(err)
	}
	defer videoModule.StopStream(token)

	<-videoHandler.QuitChan()
}
