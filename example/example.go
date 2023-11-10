package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/finder"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/wrapper"
)

// Simple example of connecting to Robomaster S1 or EP. This *REQUIRES* a
// robot broadcasting in the network. It will find the robot and connect to
// it. It will then wait for the connection to be stablished and print the
// connection status. It will then wait for the connection to be lost and
// print the connection status again. It will then exit.
func main() {
	ub := unitybridge.Get(wrapper.Get())

	// Start unity bridge.
	err := ub.Start()
	if err != nil {
		panic(err)
	}
	defer ub.Stop()

	// Listen for connection status changes.
	var wg sync.WaitGroup
	wg.Add(2) // Connection status should change twice.
	token, err := ub.AddKeyListener(key.KeyAirLinkConnection, func(data []byte) {
		// Just print whatever we get as result.
		fmt.Println(string(data))

		wg.Done()
	}, false)
	if err != nil {
		panic(err)
	}
	defer ub.RemoveKeyListener(key.KeyAirLinkConnection, token)

	ip, err := findRobot()
	if err != nil {
		panic(err)
	}

	fmt.Println("Found robot with IP:", ip)

	resetRobotConnection(ub, ip)

	time.Sleep(5 * time.Second)

	closeRobotConnection(ub)

	wg.Wait()
}

func findRobot() (net.IP, error) {
	f := finder.New(0) // Any robot in any state.
	broadcast, err := f.Find(5 * time.Second)
	if err != nil {
		return nil, err
	}

	return broadcast.SourceIp(), nil
}

// resetRobotConnection should be called whenever the IP for the robot
// changes. It is safe to call it whenever a connection needs to be
// stablished anyway.
func resetRobotConnection(ub unitybridge.UnityBridge, ip net.IP) {
	closeRobotConnection(ub)
	setRobotIPAndPort(ub, ip, 10607)
	openRobotConnection(ub)
}

func openRobotConnection(ub unitybridge.UnityBridge) {
	ev := event.NewFromTypeAndSubType(event.TypeConnection, 0)
	ub.SendEvent(ev)
}

func closeRobotConnection(ub unitybridge.UnityBridge) {
	ev := event.NewFromTypeAndSubType(event.TypeConnection, 1)
	ub.SendEvent(ev)
}

func setRobotIPAndPort(ub unitybridge.UnityBridge, ip net.IP, port uint64) {
	ev := event.NewFromTypeAndSubType(event.TypeConnection, 2)
	ub.SendEventWithString(ev, ip.String())

	ev.ResetSubType(3)
	ub.SendEventWithUint64(ev, port)
}
