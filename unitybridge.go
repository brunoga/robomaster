package unitybridge

import (
	"github.com/brunoga/unitybridge/internal"
)

type UnityBridge interface {
	// Start configures and starts the Unity Bridge.
	Start() error

	// Stop clean up and stops the Unity Bridge.
	Stop() error
}

func Get() UnityBridge {
	return internal.NewUnityBridgeImpl()
}
