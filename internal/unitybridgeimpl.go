package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
	"github.com/brunoga/unitybridge/wrapper"
)

type UnityBridgeImpl struct {
	uw               wrapper.UnityBridge
	unityBridgeDebug bool
	l                *logger.Logger
	tg               *token.Generator

	m                  sync.Mutex
	started            bool
	keyListeners       map[*key.Key]map[token.Token]result.Callback
	eventTypeListeners map[event.Type]map[token.Token]event.TypeCallback
	callbackListener   map[token.Token]result.Callback
	currentToken       uint64
}

func NewUnityBridgeImpl(uw wrapper.UnityBridge,
	unityBridgeDebug bool, l *logger.Logger) *UnityBridgeImpl {
	if l == nil {
		// Create a logger that only log errors.
		l = logger.New(slog.LevelError)
	}

	return &UnityBridgeImpl{
		uw:                 uw,
		unityBridgeDebug:   unityBridgeDebug,
		l:                  l,
		tg:                 token.NewGenerator(),
		keyListeners:       make(map[*key.Key]map[token.Token]result.Callback),
		eventTypeListeners: make(map[event.Type]map[token.Token]event.TypeCallback),
		callbackListener:   make(map[token.Token]result.Callback),
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
		u.l.Info("Unity Bridge debug mode enabled.")
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
	immediate bool) (token.Token, error) {
	if k.AccessType()&key.AccessTypeRead == 0 {
		return 0, fmt.Errorf("key %s is not readable", k)
	}

	if c == nil {
		return 0, fmt.Errorf("callback cannot be nil")
	}

	u.m.Lock()
	defer u.m.Unlock()

	if _, ok := u.keyListeners[k]; !ok {
		u.keyListeners[k] = make(map[token.Token]result.Callback)
	}

	if len(u.keyListeners[k]) == 0 {
		ev := event.NewFromTypeAndSubType(event.TypeStartListening, k.SubType())
		u.uw.SendEvent(ev.Code(), nil, 0)
	}

	token := u.tg.Next()

	u.keyListeners[k][token] = c

	if !immediate {
		return token, nil
	}

	output, err := u.GetCachedKeyValue(k)
	if err != nil {
		// Basically ignore the error and return the token anyway.
		return token, nil
	}

	c(result.NewFromJSON(output))

	return token, nil
}

func (u *UnityBridgeImpl) RemoveKeyListener(k *key.Key, token token.Token) error {
	if token == 0 {
		return fmt.Errorf("token cannot be 0")
	}

	u.m.Lock()
	defer u.m.Unlock()

	if _, ok := u.keyListeners[k]; !ok {
		return fmt.Errorf("no listeners registered for key %s", k)
	}

	if _, ok := u.keyListeners[k][token]; !ok {
		return fmt.Errorf("no listener registered with token %d for key %s",
			token, k)
	}

	delete(u.keyListeners[k], token)

	if len(u.keyListeners[k]) == 0 {
		ev := event.NewFromTypeAndSubType(event.TypeStopListening, k.SubType())
		u.uw.SendEvent(ev.Code(), nil, 0)
		delete(u.keyListeners, k)
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

	tag := u.tg.Next()

	u.callbackListener[tag] = c

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

	tag := u.tg.Next()

	u.callbackListener[tag] = c

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

	tag := u.tg.Next()

	u.callbackListener[tag] = c

	u.uw.SendEventWithString(ev.Code(), string(data), uint64(tag))

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

func (u *UnityBridgeImpl) AddEventTypeListener(t event.Type,
	c event.TypeCallback) (token.Token, error) {
	if c == nil {
		return 0, fmt.Errorf("callback cannot be nil")
	}

	u.m.Lock()
	defer u.m.Unlock()

	if _, ok := u.eventTypeListeners[t]; !ok {
		u.eventTypeListeners[t] = make(map[token.Token]event.TypeCallback)
	}

	token := u.tg.Next()

	u.eventTypeListeners[t][token] = c

	return token, nil
}

func (u *UnityBridgeImpl) RemoveEventTypeListener(t event.Type,
	token token.Token) error {
	if token == 0 {
		return fmt.Errorf("token cannot be 0")
	}

	u.m.Lock()
	defer u.m.Unlock()

	if _, ok := u.eventTypeListeners[t]; !ok {
		return fmt.Errorf("no listeners registered for event type %s", t)
	}

	if _, ok := u.eventTypeListeners[t][token]; !ok {
		return fmt.Errorf("no listener registered with token %d for event type %s",
			token, t)
	}

	delete(u.eventTypeListeners[t], token)

	if len(u.eventTypeListeners[t]) == 0 {
		delete(u.eventTypeListeners, t)
	}

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

func (u *UnityBridgeImpl) eventCallback(eventCode uint64, data []byte, tag uint64) {
	e := event.NewFromCode(eventCode)

	var dataType event.DataType
	dataType, tag = event.DataTypeFromTag(tag)

	u.m.Lock()
	defer u.m.Unlock()

	// Call all registered event type listeners.
	u.notifyEventTypeListeners(e, data, dataType)

	if tag != 0 {
		// This should be associated with a callback. Call it.
		u.notifyCallbacks(data, tag)
	}

	if e.SubType() == 0 {
		// This is not a key event.
		return
	}

	k, err := key.FromEvent(e)
	if err != nil {
		// TODO(bga): This is actually expected. Consider removing this
		//            eventually.
		u.l.Warn("Error creating key from event", "event", e, "err", err)
	} else {
		u.notifyKeyListeners(k, data)
	}
}

func (u *UnityBridgeImpl) notifyEventTypeListeners(e *event.Event,
	data []byte, dataType event.DataType) {
	if _, ok := u.eventTypeListeners[e.Type()]; ok {
		for _, c := range u.eventTypeListeners[e.Type()] {
			c(data, dataType)
		}
	} else {
		u.l.Warn("No listeners registered for event type", "eventType", e.Type())
	}
}

func (u *UnityBridgeImpl) notifyKeyListeners(k *key.Key, data []byte) {
	if _, ok := u.keyListeners[k]; ok {
		for _, c := range u.keyListeners[k] {
			c(result.NewFromJSON(data))
		}
	} else {
		u.l.Warn("No listeners registered for key", "key", k)
	}
}

func (u *UnityBridgeImpl) notifyCallbacks(data []byte, tag uint64) {
	if c, ok := u.callbackListener[token.Token(tag)]; ok {
		c(result.NewFromJSON(data))
		delete(u.callbackListener, token.Token(tag))
	} else {
		u.l.Error("No callback registered for tag", "tag", tag)
	}
}
