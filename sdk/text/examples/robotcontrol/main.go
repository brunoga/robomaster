package main

import (
	"flag"

	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/text/examples/robotcontrol/scenes"
	"github.com/brunoga/robomaster/sdk/text/modules/robot"
)

// Flags
var (
	mirrors = flag.Uint("mirrors", 0, "number of mirror robots")
)

func main() {
	flag.Parse()

	client, err := sdk.NewClient(nil)
	if err != nil {
		panic(err)
	}

	err = client.Open()
	if err != nil {
		panic(err)
	}

	client.RobotModule().SetMotionMode(robot.MotionModeGimbalLead)

	var mirrorClients []*sdk.Client = nil
	if *mirrors > 0 {
		for i := 0; i < int(*mirrors); i++ {
			mirrorClient, err := sdk.NewClient(nil)
			if err != nil {
				// Just ignore. We do not care that much about the
				// mirror robots.
				continue
			}

			err = mirrorClient.Open()
			if err != nil {
				// Ditto.
				continue
			}

			mirrorClient.RobotModule().SetMotionMode(robot.MotionModeGimbalLead)

			mirrorClients = append(mirrorClients, mirrorClient)
		}
	}

	opts := engo.RunOptions{
		Title:         "Robot Control",
		Width:         1280,
		Height:        720,
		VSync:         true,
		ScaleOnResize: true,
		FPSLimit:      60,
		NotResizable:  true,
	}

	engo.Run(opts, &scenes.RobotControl{
		Client:        client,
		MirrorClients: mirrorClients,
	})
}
