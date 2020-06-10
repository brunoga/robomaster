package main

import (
	"fmt"
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

	// Setup yaw and pitch PID controllers.
	pidYaw := pid.NewPIDController(250, 10, 30, -400, 400)
	pidPitch := pid.NewPIDController(150, 10, 20, -400, 400)

	window := gocv.NewWindow("Robomaster S1")
	window.ResizeWindow(640, 360)

	ballTracker := support.NewColorObjectTracker(35, 219, 90, 119, 255,
		255, 10)

	_, err = videoModule.StartStream(func(frame *gocv.Mat, wg *sync.WaitGroup) {
		x, y, radius, err := ballTracker.FindLargestObject(frame)
		if err == nil {
			// Found something. Draw a circle around it.
			gocv.Circle(frame, image.Point{X: int(x), Y: int(y)}, int(radius),
				color.RGBA{R: 0, G: 255, B: 255, A: 255}, 2)

			// Get errors in the x and y axis normalized to [-0.5, 0.5]
			errX := float64(x-(sdk.CameraHorizontalResolutionPoints/2)) /
				sdk.CameraHorizontalResolutionPoints
			errY := float64((sdk.CameraVerticalResolutionPoints/2)-y) /
				sdk.CameraVerticalResolutionPoints

			if math.Abs(errX) > 0.0 || math.Abs(errY) > 0.0 {
				err = gimbalModule.SetSpeed(pidPitch.Output(errY),
						pidYaw.Output(errX))
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		window.IMShow(*frame)
		window.WaitKey(1)
		if window.GetWindowProperty(gocv.WindowPropertyAspectRatio) == -1.0 {
			panic("window closed")
		}

		wg.Done()
	})
	if err != nil {
		panic(err)
	}

	select {}
}
