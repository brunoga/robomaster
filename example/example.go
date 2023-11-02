package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/event"
)

func callbackHandler(eventCode uint64, data []byte, tag uint64) {
	ev := event.NewFromCode(eventCode)

	dataType := (tag >> 56) & 0xff

	dataTypeStr := "unknown"
	switch dataType {
	case 0:
		dataTypeStr = "string"
	case 1:
		dataTypeStr = "uint64"
	}

	tag = tag & 0x0000ffffffffffff

	fmt.Printf("Callback handler called for event with type %s, sub-type %d, data type %s and tag %d\n",
		ev.Type(), ev.SubType(), dataTypeStr, tag)

	if dataType == 0 {
		fmt.Printf("Data: %s\n", string(data))
	} else {
		fmt.Printf("Data: %d\n", binary.NativeEndian.Uint64(data))
	}
}

func main() {
	ub := unitybridge.Get()

	ub.Create("Robomaster", true, "./log")
	defer ub.Destroy()

	if !ub.Initialize() {
		panic("Could not initialize UnityBridge.")
	}
	defer ub.Uninitialize()

	for _, typ := range event.TypeValues() {
		ev := event.NewFromType(typ)
		ub.SetEventCallback(ev.Code(), callbackHandler)
	}

	time.Sleep(5 * time.Second)

	for _, typ := range event.TypeValues() {
		ev := event.NewFromType(typ)
		ub.SetEventCallback(ev.Code(), nil)
	}
}
