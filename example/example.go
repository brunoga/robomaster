package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/wrapper"
)

// Simple example of connecting to Robomaster S1 or EP. This assumes the IP
// is known.
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <robot_ip>\n", os.Args[0])
		os.Exit(1)
	}

	ub := unitybridge.Get(wrapper.Get())

	// Start unity bridge.
	err := ub.Start()
	if err != nil {
		panic(err)
	}
	defer ub.Stop()

	// Listen for connection status changes.
	var wg sync.WaitGroup
	wg.Add(1)
	token, err := ub.AddKeyListener(key.KeyAirLinkConnection, func(data []byte) {
		// Just print whatever we get as result.
		fmt.Println(string(data))
		wg.Done()
	}, false)
	if err != nil {
		panic(err)
	}
	defer ub.RemoveKeyListener(key.KeyAirLinkConnection, token)

	// We will be sending connection events.
	ev := event.NewFromType(event.TypeConnection)

	// Set connection type (?).
	ev.ResetSubType(1)
	ub.SendEvent(ev)

	// Set robot IP.
	ev.ResetSubType(2)
	ub.SendEventWithString(ev, os.Args[1])

	// Set port.
	ev.ResetSubType(3)
	ub.SendEventWithUint64(ev, 10607)

	// Connect.
	ev.ResetSubType(0)
	ub.SendEvent(ev)

	wg.Wait()
}
