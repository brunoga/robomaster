package main

import (
	"github.com/brunoga/robomaster/legacy/unitybridge"
)

func main() {
	unitybridge.Create("Robomaster", true, "./log")
	defer unitybridge.Destroy()

	if !unitybridge.Initialize() {
		panic("UnityBridge failed to initialize")
	}
	defer unitybridge.Uninitialize()
}

