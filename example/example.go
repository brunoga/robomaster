package main

import (
	"time"

	"github.com/brunoga/unitybridge"
)

func main() {
	ub := unitybridge.Get()

	ub.Create("Robomaster", true, "./log")
	defer ub.Destroy()

	if !ub.Initialize() {
		panic("Could not initialize UnityBridge.")
	}
	defer ub.Uninitialize()

	time.Sleep(5 * time.Second)
}
