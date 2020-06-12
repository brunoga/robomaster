package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/modules"
	"github.com/brunoga/robomaster/sdk/support"
	"github.com/brunoga/robomaster/sdk/support/pid"
	"gocv.io/x/gocv"
)

// Flags.
var (
	hsvLower = flag.String("hsvlower", "35,219,90",
		"lower bound for color filtering (h,s,v)")
	hsvUpper = flag.String("hsvupper", "119,255,255",
		"lower bound for color filtering (h,s,v)")
)

type exampleVideoHandler struct {
	window       *gocv.Window
	tracker      *support.ColorObjectTracker
	gimbalModule *modules.Gimbal
	pidPitch     pid.Controller
	pidYaw       pid.Controller
	quitChan     chan struct{}
}

func parseHSVValues(hsvString string) (float64, float64, float64, error) {
	components := strings.Split(hsvString, ",")
	if len(components) != 3 {
		return -1, -1, -1, fmt.Errorf("invalid hsv values")
	}

	h, err := strconv.Atoi(components[0])
	if err != nil {
		return -1, -1, -1, fmt.Errorf("invalid h value")
	}

	s, err := strconv.Atoi(components[1])
	if err != nil {
		return -1, -1, -1, fmt.Errorf("invalid s value")
	}

	v, err := strconv.Atoi(components[2])
	if err != nil {
		return -1, -1, -1, fmt.Errorf("invalid v value")
	}

	return float64(h), float64(s), float64(v), nil
}

func newExampleVideoHandler(
	gimbalModule *modules.Gimbal) (*exampleVideoHandler, error) {
	window := gocv.NewWindow("Robomaster")
	window.ResizeWindow(sdk.CameraHorizontalResolutionPoints/2,
		sdk.CameraVerticalResolutionPoints/2)

	hl, sl, vl, err := parseHSVValues(*hsvLower)
	if err != nil {
		return nil, err
	}

	hu, su, vu, err := parseHSVValues(*hsvUpper)
	if err != nil {
		return nil, err
	}

	return &exampleVideoHandler{
		window,
		support.NewColorObjectTracker(hl, sl, vl, hu, su, vu, 10),
		gimbalModule,
		pid.NewPIDController(150, 10, 20, -400, 400),
		pid.NewPIDController(200, 9, 25, -400, 400),
		make(chan struct{}),
	}, nil
}

func (e *exampleVideoHandler) QuitChan() <-chan struct{} {
	return e.quitChan
}

func (e *exampleVideoHandler) HandleFrame(frame *gocv.Mat, wg *sync.WaitGroup) {
	// Automatically notify we processed the frame on return.
	defer wg.Done()

	// Clone frame as we are going to modify it. We could potentially call
	// wg.Done() right after this but without implementing a queue, this is not
	// a good idea (frames will be racing against each other).
	outputFrame := frame.Clone()

	x, y, radius, err := e.tracker.FindLargestObject(frame)
	if err == nil {
		// Found something. Draw a circle around it into our modified frame.
		gocv.Circle(&outputFrame, image.Point{X: int(x), Y: int(y)},
			int(radius), color.RGBA{R: 0, G: 255, B: 255, A: 255}, 2)

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

	// Show output frame and wait for an event.
	e.window.IMShow(outputFrame)
	e.window.WaitKey(1)

	// Hack to detect that the window was closed.
	if e.window.GetWindowProperty(gocv.WindowPropertyAspectRatio) == -1.0 {
		// Window closed. Notify listeners.
		close(e.quitChan)
	}
}

func main() {
	flag.Parse()

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

	// Enable gimbal attitude push events.
	token, err := gimbalModule.StartGimbalPush(
		modules.GimbalPushAttributeAttitude, func(data string) {
			// Just print all events we get.
			fmt.Println("Push :", data)
		})
	if err != nil {
		panic(err)
	}
	defer gimbalModule.StopGimbalPush(
		modules.GimbalPushAttributeAttitude, token)

	// Reset gimbal.
	err = gimbalModule.Recenter()
	if err != nil {
		panic(err)
	}

	videoHandler, err := newExampleVideoHandler(gimbalModule)
	if err != nil {
		panic(err)
	}

	token, err = videoModule.StartStream(videoHandler.HandleFrame)
	if err != nil {
		panic(err)
	}
	defer videoModule.StopStream(token)

	<-videoHandler.QuitChan()
}
