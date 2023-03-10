package main

import (
	"flag"

	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/examples/text/robotcontrol/scenes"
	"github.com/brunoga/robomaster/sdk/modules/robot"
	"github.com/brunoga/robomaster/sdk/types"
)

// Flags
var (
	mirrors = flag.Uint("mirrors", 0, "number of mirror robots")
)

func main() {
	flag.Parse()

	s, err := sdk.New(types.SDKProtocolText, nil)
	if err != nil {
		panic(err)
	}

	err = s.Open(types.ConnectionModeInfrastructure,
		types.ConnectionProtocolTCP, nil)
	if err != nil {
		panic(err)
	}

	s.Robot().SetMode(robot.ModeGimbalLead)

	var mirrorClients []sdk.SDK = nil
	if *mirrors > 0 {
		for i := 0; i < int(*mirrors); i++ {
			mirrorClient, err := sdk.New(types.SDKProtocolText, nil)
			if err != nil {
				// Just ignore. We do not care that much about the
				// mirror robots.
				continue
			}

			err = mirrorClient.Open(types.ConnectionModeInfrastructure,
				types.ConnectionProtocolTCP, nil)
			if err != nil {
				// Ditto.
				continue
			}

			mirrorClient.Robot().SetMode(robot.ModeGimbalLead)

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
		Sdk:           s,
		MirrorClients: mirrorClients,
	})
}
