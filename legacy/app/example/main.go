package main

import (
	"flag"

	"git.bug-br.org.br/bga/robomasters1/app/example/scenes"
	"github.com/EngoEngine/engo"
)

var (
	ssID = flag.String("ssid", "testssid",
		"wifi network to connect to")
	password = flag.String("password", "testpassword", "wifi password")
	appID    = flag.Uint64("appid", 0, "if provided, use this app ID "+
		"instead of creating a new one")
)

func main() {
	flag.Parse()

	opts := engo.RunOptions{
		Title:         "Robomaster",
		Width:         1280,
		Height:        720,
		VSync:         true,
		ScaleOnResize: true,
		FPSLimit:      60,
	}

	engo.Run(opts, &scenes.Robomaster{
		*ssID,
		*password,
		*appID,
	})
}
