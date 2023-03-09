package main

import (
	"fmt"

	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/types"
)

func main() {
	sdk, err := sdk.New(types.SDKProtocolText, nil)
	if err != nil {
		panic(err)
	}

	err = sdk.Open(types.ConnectionModeInfrastructure, types.ConnectionProtocolTCP, nil)
	if err != nil {
		panic(err)
	}
	defer sdk.Close()

	version, err := sdk.Robot().GetSDKVersion()
	if err != nil {
		panic(err)
	}

	fmt.Println("Text Mode SDK version:", version)
}
