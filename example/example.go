package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/finder"
	"github.com/brunoga/unitybridge/support/qrcode"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
	"github.com/brunoga/unitybridge/wrapper"
)

var (
	ssid     = flag.String("ssid", "", "SSID of the network to connect to.")
	password = flag.String("password", "", "Password of the network to connect "+
		"to.")
	appID = flag.Uint64("app_id", 0, "App ID to use. If 0, a random one will be "+
		"generated.")
)

// Simple example of connecting to Robomaster S1 or EP. This *REQUIRES* a
// robot broadcasting in the network. It will find the robot and connect to
// it. It will then wait for the connection to be stablished and print the
// connection status. Then it will wait for the connection to be lost and
// print the connection status again.
func main() {
	flag.Parse()

	if strings.TrimSpace(*ssid) == "" || strings.TrimSpace(*password) == "" {
		panic("SSID and password must be provided.")
	}

	ub := unitybridge.Get(wrapper.Get(), true)

	// Start unity bridge.
	err := ub.Start()
	if err != nil {
		panic(err)
	}
	defer ub.Stop()

	// We are an app, so generate our app ID.
	if *appID == 0 {
		*appID, err = support.GenerateAppID()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Using App ID:", *appID)

	// And generate a QRCode to pair a Robomaster.
	qrCode, err := qrcode.New(*appID, "CN", "Discworld", "zwergschnauzer", "")
	if err != nil {
		panic(err)
	}

	// Print the QRCode.
	fmt.Println(qrCode.Text())

	f := finder.New(*appID)
	var robotIP net.IP

	// Listen for VideoTransferSpeed type events as these starting comming right away.
	ub.AddEventTypeListener(event.TypeVideoTransferSpeed,
		func(eventCode uint64, data []byte, tag uint64) {
			fmt.Println("Video Transfer Speed:", binary.LittleEndian.Uint64(data))
		})

	// Listen for connection status changes.
	var wg sync.WaitGroup
	wg.Add(2) // Connection status should change twice.
	token, err := ub.AddKeyListener(key.KeyAirLinkConnection, func(r *result.Result) {
		// Just print whatever we get as result.
		fmt.Println(r)

		if !r.Succeeded() {
			fmt.Println("Result error:", r.ErrorDesc())
		} else {
			// Expected value is a bool.
			value := r.Value().(bool)
			if value {
				fmt.Println("Connected to robot.")
			} else {
				fmt.Println("Disconnected from robot.")
			}
		}

		wg.Done()
	}, false)
	if err != nil {
		panic(err)
	}
	defer ub.RemoveKeyListener(key.KeyAirLinkConnection, token)

	// Find a robot. Wait for up to 1 minute.
	broadcast, err := f.Find(1 * time.Minute)
	if err != nil {
		panic(err)
	}

	robotIP = broadcast.SourceIp()
	fmt.Println("Found robot at", robotIP)

	// Setup connection and connect to robot.
	resetRobotConnection(ub, robotIP)

	time.Sleep(5 * time.Second)

	closeRobotConnection(ub)

	wg.Wait()
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
