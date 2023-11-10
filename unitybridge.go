package unitybridge

import (
	"github.com/brunoga/unitybridge/internal"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/wrapper"
)

type UnityBridge interface {
	// Start configures and starts the Unity Bridge.
	Start() error

	// AddKeyListener adds a listener for events on the given key. If
	// immediate is true, the callback will be called immediatelly with any
	// cached value associated with the key. Returns a token that can be used
	// to remove the listener later.
	AddKeyListener(k *key.Key, c event.Callback, immediate bool) (uint64, error)

	// RemoveKeyListener removes the listener associated with the given token
	// for events on the given key.
	RemoveKeyListener(key *key.Key, token uint64) error

	// GetKeyValue returns the Unity Bridge value associated with the given
	// key.
	GetKeyValue(k *key.Key, c event.Callback) error

	// GetCachedKeyValue returns the Unity Bridge cached value associated
	// with the given key.
	GetCachedKeyValue(k *key.Key) ([]byte, error)

	// SetKeyValue sets the Unity Bridge value associated with the given key.
	SetKeyValue(k *key.Key, value any, c event.Callback) error

	// PerformActionForKey performs the Unity Bridge action associated with the
	// given key with the given value as parameter.
	PerformActionForKey(k *key.Key, value any, c event.Callback) error

	// DirectSendKeyValue sends the given value to the Unity Bridge for the
	// given key. This is a low level function that should be used with care.
	DirectSendKeyValue(k *key.Key, value uint64) error

	// SendEvent sends the given event to the Unity Bridge. This is a low level
	// function that should be used with care.
	SendEvent(ev *event.Event) error

	// DirectSendKeyValue sends the given event associated with the given string
	// data to the Unity Bridge. This is a low level function that should be
	// used with care.
	SendEventWithString(ev *event.Event, data string) error

	// SendEventWithUint64 sends the given event associated with the given
	// uint64 data to the Unity Bridge. This is a low level function that should
	// be used with care.
	SendEventWithUint64(ev *event.Event, data uint64) error

	// Stop cleans up and stops the Unity Bridge.
	Stop() error
}

// Get returns an instance of the high level Unity Bridge API using the given
// low-level Unity Bridge library wrapper (mostly so irt can be mocked for
// tests).
func Get(wu wrapper.UnityBridge, unityBridgeDebug bool) UnityBridge {
	return internal.NewUnityBridgeImpl(wu, unityBridgeDebug)
}
