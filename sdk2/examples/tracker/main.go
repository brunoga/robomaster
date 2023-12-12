package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/robomaster/sdk2/examples/tracker/mode"
	"github.com/brunoga/robomaster/sdk2/module/camera"
	"github.com/brunoga/robomaster/sdk2/module/chassis"
	"github.com/brunoga/robomaster/sdk2/module/gimbal"
	"github.com/brunoga/robomaster/sdk2/support/pid"
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
	tracker      *mode.ColorObject
	gimbalModule *gimbal.Gimbal
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
	gimbalModule *gimbal.Gimbal) (*exampleVideoHandler, error) {
	window := gocv.NewWindow("Robomaster")
	window.ResizeWindow(camera.HorizontalResolutionPoints/2,
		camera.VerticalResolutionPoints/2)

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
		mode.NewColorObject(hl, sl, vl, hu, su, vu, 10),
		gimbalModule,
		pid.NewPIDController(400, 0, 0, -400, 400),
		pid.NewPIDController(800, 0, 20, -400, 400),
		//gimbal.NewSpeed(0.0, 0.0),
		make(chan struct{}),
	}, nil
}

func (e *exampleVideoHandler) QuitChan() <-chan struct{} {
	return e.quitChan
}

func (e *exampleVideoHandler) HandleFrame(frame *camera.RGB) {
	// Do not explicitly close inFrameRGBA as the underlying pixel data is
	// managed by Go itself.
	inFrameRGBA, err := gocv.NewMatFromBytes(camera.VerticalResolutionPoints,
		camera.HorizontalResolutionPoints, gocv.MatTypeCV8UC4, frame.Pix)
	if err != nil {
		return
	}

	// A Mat underlying pixel format is BRG, so we convert to it. Note we use
	// ColorBGRToRGB as the color conversion code as there is no ColorRGBToBGR
	// (which is fine as the conversion works both ways anyway).
	inFrame := gocv.NewMatWithSize(720, 1280, gocv.MatTypeCV8SC3)
	gocv.CvtColor(inFrameRGBA, &inFrame, gocv.ColorBGRToRGB)
	defer inFrame.Close()

	// Clone frame as we are going to modify it. We could potentially call
	// wg.Done() right after this but without implementing a queue, this is not
	// a good idea (frames will be racing against each other).
	outputFrame := inFrame.Clone()
	defer outputFrame.Close()

	x, y, radius, err := e.tracker.FindLargestObject(&inFrame)
	if err == nil {
		// Found something. Draw a circle around it into our modified frame.
		gocv.Circle(&outputFrame, image.Point{X: int(x), Y: int(y)},
			int(radius), color.RGBA{R: 0, G: 255, B: 255, A: 255}, 2)

		// Get errors in the x and y axis normalized to [-0.5, 0.5]
		errX := float64(x-(camera.HorizontalResolutionPoints/2)) /
			camera.HorizontalResolutionPoints
		errY := float64((camera.VerticalResolutionPoints/2)-y) /
			camera.VerticalResolutionPoints

		// If there is some error (object in not in center of image), move the
		// gimbal to minimize it.
		if math.Abs(errX) > 0.0 || math.Abs(errY) > 0.0 {
			// Move the gimbal with a speed determined by the pitch and yaw PID
			// controllers.
			yawSpeed := int16(e.pidYaw.Output(errX))
			pitchSpeed := int16(e.pidPitch.Output(errY))
			err = e.gimbalModule.SetSpeed(pitchSpeed, yawSpeed)
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

	client, err := sdk2.New(nil, 0)
	if err != nil {
		panic(err)
	}

	err = client.Start()
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	// Obtain references to the modules we are interested in.
	gimbalModule := client.Gimbal()
	chassisModule := client.Chassis()
	cameraModule := client.Camera()

	err = chassisModule.SetMode(chassis.ModeFPV)
	if err != nil {
		panic(err)
	}

	// Reset gimbal.
	gimbalModule.ResetPosition()

	videoHandler, err := newExampleVideoHandler(gimbalModule)
	if err != nil {
		panic(err)
	}

	token, err := cameraModule.AddVideoCallback(videoHandler.HandleFrame)
	if err != nil {
		panic(err)
	}
	defer cameraModule.RemoveVideoCallback(token)

	<-videoHandler.QuitChan()
}
