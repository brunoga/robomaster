package main

import (
	"fmt"
	"os"

	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/support/logger"
	"github.com/brunoga/robomaster/sdk/types"
)

func main() {
	// Set a logger that we allow us seeing all debug messages.
	l := logger.New(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	sdk, err := sdk.New(types.SDKProtocolBinary, l)
	if err != nil {
		panic(err)
	}

	err = sdk.Open(types.ConnectionModeInfrastructure, types.ConnectionProtocolUDP, nil)
	if err != nil {
		panic(err)
	}
	defer sdk.Close()

	version, err := sdk.Robot().GetProductVersion()
	if err != nil {
		panic(err)
	}

	fmt.Println("Binary Mode SDK version", version)
}
