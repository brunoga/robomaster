package main

import (
	"fmt"
	"time"

	"github.com/brunoga/robomaster/sdk"
)

func main() {
	binarySdk, err := sdk.New(sdk.Binary)
	if err != nil {
		panic(err)
	}

	textSdk, err := sdk.New(sdk.Text)
	if err != nil {
		panic(err)
	}

	err = binarySdk.Find(nil, 10*time.Second)
	if err != nil {
		panic(err)
	}

	err = textSdk.Find(nil, 10*time.Second)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println("")
		fmt.Printf("Detected text mode SDK robots   : %d\n", textSdk.NumRobots())
		fmt.Printf("Detected binary mode SDK robots : %d\n", binarySdk.NumRobots())
		time.Sleep((1 * time.Second))
	}
}
