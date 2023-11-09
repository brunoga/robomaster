package internal

import (
	"fmt"

	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/wrapper"
)

type UnityBridgeImpl struct {
	wu wrapper.UnityBridge
}

func NewUnityBridgeImpl() *UnityBridgeImpl {
	return &UnityBridgeImpl{
		wu: wrapper.Get(),
	}
}

func (u *UnityBridgeImpl) Start() error {
	u.wu.Create("Robomaster", true, "./log")
	if !u.wu.Initialize() {
		return fmt.Errorf("failed to initialize Unity Bridge library")
	}

	for _, eventType := range event.AllTypes() {
		eventTypeCode := event.NewFromType(eventType).Code()
		u.wu.SetEventCallback(eventTypeCode, u.eventCallback)
	}

	return nil
}

func (u *UnityBridgeImpl) Stop() error {
	for _, eventType := range event.AllTypes() {
		eventTypeCode := event.NewFromType(eventType).Code()
		u.wu.SetEventCallback(eventTypeCode, nil)
	}

	u.wu.Uninitialize()
	u.wu.Destroy()

	return nil
}

func (u *UnityBridgeImpl) eventCallback(eventCode uint64, data []byte, tag uint64) {
	// TODO(bga): Dispatch events to any interested parties.
}
