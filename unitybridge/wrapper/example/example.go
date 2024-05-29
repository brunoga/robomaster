package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge/wrapper"
)

var (
	allEventTypes = []uint64{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 100, 101, 200, 300, 301, 302, 303, 304, 305,
		306, 500}
)

func callbackHandler(eventCode uint64, data []byte, tag uint64) {
	dataType := tag >> 56
	tag = tag & 0x00FFFFFFFFFFFFFF

	fmt.Printf("Callback handler called for event code %d, type %d and tag "+
		"%d\n", eventCode, eventCode>>32, tag)

	if dataType == 0 {
		fmt.Printf("Data len=%d (string): %v\n", len(data), string(data))
	} else if dataType == 1 {
		fmt.Printf("Data len=%d (uint64): %v\n", len(data), binary.LittleEndian.Uint64(data))
	} else {
		fmt.Printf("Data len=%d (unrecognized type): %v\n", len(data), data)
	}
}

func main() {
	l := logger.New(slog.LevelDebug)

	ub := wrapper.Get(l)

	ub.Create("Robomaster", true, "./log")
	defer ub.Destroy()

	if !ub.Initialize() {
		panic("Could not initialize UnityBridge. Did you call Create() and " +
			"passed \"Robomaster\" as name?")
	}
	defer ub.Uninitialize()

	fmt.Println("Started listening for events. We should get callbacks.")
	for _, typ := range allEventTypes {
		ub.SetEventCallback(typ<<32, callbackHandler)
	}

	time.Sleep(5 * time.Second)

	for _, typ := range allEventTypes {
		ub.SetEventCallback(typ<<32, nil)
	}
	fmt.Println("Stopped listening for events. We should not get callbacks anymore.")

	time.Sleep(5 * time.Second)
}
