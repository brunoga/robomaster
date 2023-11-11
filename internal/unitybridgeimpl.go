package internal

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
	"github.com/brunoga/unitybridge/wrapper"
)

type UnityBridgeImpl struct {
	uw               wrapper.UnityBridge
	unityBridgeDebug bool

	m            sync.Mutex
	started      bool
	listeners    map[*key.Key]map[uint64]result.Callback
	callbacks    map[uint64]result.Callback
	currentToken uint64
}

func NewUnityBridgeImpl(uw wrapper.UnityBridge,
	unityBridgeDebug bool) *UnityBridgeImpl {
	return &UnityBridgeImpl{
		uw:               uw,
		unityBridgeDebug: unityBridgeDebug,
		listeners:        make(map[*key.Key]map[uint64]result.Callback),
		callbacks:        make(map[uint64]result.Callback),
	}
}

func (u *UnityBridgeImpl) Start() error {
	u.m.Lock()
	defer u.m.Unlock()

	if u.started {
		return fmt.Errorf("unity bridge already started")
	}

	var logPath string
	if u.unityBridgeDebug {
		logPath = "./log"
	}

	u.uw.Create("Robomaster", u.unityBridgeDebug, logPath)
	if !u.uw.Initialize() {
		return fmt.Errorf("failed to initialize Unity Bridge library")
	}

	for _, eventType := range event.AllTypes() {
		eventTypeCode := event.NewFromType(eventType).Code()
		u.uw.SetEventCallback(eventTypeCode, u.eventCallback)
	}

	u.started = true

	return nil
}

func (u *UnityBridgeImpl) AddKeyListener(k *key.Key, c result.Callback,
	immediate bool) (uint64, error) {
	if k.AccessType()&key.AccessTypeRead == 0 {
		return 0, fmt.Errorf("key %s is not readable", k)
	}

	if c == nil {
		return 0, fmt.Errorf("callback cannot be nil")
	}

	u.m.Lock()
	defer u.m.Unlock()

	if _, ok := u.listeners[k]; !ok {
		u.listeners[k] = make(map[uint64]result.Callback)
	}

	if len(u.listeners[k]) == 0 {
		ev := event.NewFromTypeAndSubType(event.TypeStartListening, k.SubType())
		u.uw.SendEvent(ev.Code(), nil, 0)
	}

	token := u.getAndUpdateTokenLocked()

	u.listeners[k][token] = c

	if !immediate {
		return token, nil
	}

	output, err := u.GetCachedKeyValue(k)
	if err != nil {
		// Basically ignore the error and return the token anyway.
		return token, nil
	}

	c(output)

	return token, nil
}

func (u *UnityBridgeImpl) RemoveKeyListener(k *key.Key, token uint64) error {
	if token == 0 {
		return fmt.Errorf("token cannot be 0")
	}

	u.m.Lock()
	defer u.m.Unlock()

	if _, ok := u.listeners[k]; !ok {
		return fmt.Errorf("no listeners registered for key %s", k)
	}

	if _, ok := u.listeners[k][token]; !ok {
		return fmt.Errorf("no listener registered with token %d for key %s",
			token, k)
	}

	delete(u.listeners[k], token)

	if len(u.listeners[k]) == 0 {
		ev := event.NewFromTypeAndSubType(event.TypeStopListening, k.SubType())
		u.uw.SendEvent(ev.Code(), nil, 0)
	}

	return nil
}

func (u *UnityBridgeImpl) GetKeyValue(k *key.Key, c result.Callback) error {
	if k.AccessType()&key.AccessTypeRead == 0 {
		return fmt.Errorf("key %s is not readable", k)
	}

	ev := event.NewFromTypeAndSubType(event.TypeGetValue, k.SubType())

	u.m.Lock()
	defer u.m.Unlock()

	tag := u.getAndUpdateTokenLocked()

	u.callbacks[tag] = c

	u.uw.SendEvent(ev.Code(), nil, uint64(tag))

	return nil
}

func (u *UnityBridgeImpl) GetCachedKeyValue(k *key.Key) ([]byte, error) {
	if k.AccessType()&key.AccessTypeRead == 0 {
		return nil, fmt.Errorf("key %s is not readable", k)
	}

	ev := event.NewFromTypeAndSubType(event.TypeGetAvailableValue, k.SubType())

	output := make([]byte, 2048)
	u.uw.SendEvent(ev.Code(), output, 0)

	n := bytes.IndexByte(output, 0)

	if n != -1 {
		output = output[:n]
	}

	return output, nil
}

func (u *UnityBridgeImpl) SetKeyValue(k *key.Key, value any,
	c result.Callback) error {
	if k.AccessType()&key.AccessTypeWrite == 0 {
		return fmt.Errorf("key %s is not writable", k)
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ev := event.NewFromTypeAndSubType(event.TypeSetValue, k.SubType())

	u.m.Lock()
	defer u.m.Unlock()

	tag := u.getAndUpdateTokenLocked()

	u.callbacks[tag] = c

	u.uw.SendEventWithString(ev.Code(), string(data), uint64(tag))

	return nil
}

func (u *UnityBridgeImpl) PerformActionForKey(k *key.Key, value any,
	c result.Callback) error {
	if k.AccessType()&key.AccessTypeAction == 0 {
		return fmt.Errorf("key %s is not an action", k)
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ev := event.NewFromTypeAndSubType(event.TypePerformAction, k.SubType())

	tag := u.getAndUpdateTokenLocked()

	u.callbacks[tag] = c

	u.uw.SendEventWithString(ev.Code(), string(data), tag)

	return nil
}

func (u *UnityBridgeImpl) DirectSendKeyValue(k *key.Key,
	value uint64) error {
	// This always sends the event as an action but it never checks for the
	// action access type. This is intentional as we want to be able to do this
	// for any key regardless of access type.
	ev := event.NewFromTypeAndSubType(event.TypePerformAction, k.SubType())

	u.uw.SendEventWithNumber(ev.Code(), value, 0)

	return nil
}

func (u *UnityBridgeImpl) Stop() error {
	u.m.Lock()
	defer u.m.Unlock()

	if !u.started {
		return fmt.Errorf("unity bridge not started")
	}

	for _, eventType := range event.AllTypes() {
		eventTypeCode := event.NewFromType(eventType).Code()
		u.uw.SetEventCallback(eventTypeCode, nil)
	}

	u.uw.Uninitialize()
	u.uw.Destroy()

	u.started = false

	return nil
}

func (u *UnityBridgeImpl) SendEvent(ev *event.Event) error {
	u.uw.SendEvent(ev.Code(), nil, 0)

	return nil
}

func (u *UnityBridgeImpl) SendEventWithString(ev *event.Event,
	data string) error {
	u.uw.SendEventWithString(ev.Code(), data, 0)

	return nil
}

func (u *UnityBridgeImpl) SendEventWithUint64(ev *event.Event,
	data uint64) error {
	u.uw.SendEventWithNumber(ev.Code(), data, 0)

	return nil
}

func (u *UnityBridgeImpl) eventCallback(eventCode uint64, data []byte, tag uint64) {
	e := event.NewFromCode(eventCode)

	var dataType event.DataType
	dataType, tag = event.DataTypeFromTag(tag)

	k, err := key.FromEvent(e)
	if err != nil {
		// Do nothing as this just mean we got a non-key event.
		//
		// TODO(bga): Log now while we debug things but remove this later.
		fmt.Printf("Got non-key event with type %s, sub-type %d, dataType %s "+
			"and tag %d.\n", e.Type(), e.SubType(), dataType, tag)
		switch dataType {
		case event.StringDataType:
			fmt.Printf("Data: %s\n", string(data))
		case event.NumberDataType:
			fmt.Printf("Data: %d\n", binary.LittleEndian.Uint64(data))
		}
	}

	u.m.Lock()

	if _, ok := u.listeners[k]; ok {
		for _, c := range u.listeners[k] {
			c(data)
		}
	}

	if c, ok := u.callbacks[tag]; ok {
		c(data)
		delete(u.callbacks, tag)
	}

	u.m.Unlock()
}

func (u *UnityBridgeImpl) getAndUpdateTokenLocked() uint64 {
	if u.currentToken == 0 {
		u.currentToken = 1 // Never use 0.
	}

	next := u.currentToken

	u.currentToken++

	return next
}
