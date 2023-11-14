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

	m                  sync.RWMutex
	started            bool
	keyListeners       map[*key.Key]map[token.Token]result.Callback
	eventTypeListeners map[event.Type]map[token.Token]event.TypeCallback
	callbackListener   map[token.Token]result.Callback
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

	if u.started {
		u.m.Unlock()
		return fmt.Errorf("unity bridge already started")
	}

	u.started = true

	u.m.Unlock()

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

	t := u.tg.Next()

	u.m.Lock()

	if _, ok := u.keyListeners[k]; !ok {
		u.keyListeners[k] = make(map[token.Token]result.Callback)
	}

	if len(u.keyListeners[k]) == 0 {
		ev := event.NewFromTypeAndSubType(event.TypeStartListening, k.SubType())
		u.uw.SendEvent(ev.Code(), nil, 0)
	}

	u.keyListeners[k][t] = c

	u.m.Unlock()

	if !immediate {
		return t, nil
	}

	r, err := u.GetCachedKeyValue(k)
	if err != nil || r.ErrorCode() != 0 {
		// Basically ignore the error and return the token anyway.
		return t, nil
	}

	c(r)

	return t, nil
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

	tag := u.tg.Next()

	u.m.Lock()

	u.callbackListener[tag] = c

	u.m.Unlock()

	u.uw.SendEvent(ev.Code(), nil, uint64(tag))

	return nil
}

func (u *UnityBridgeImpl) GetCachedKeyValue(k *key.Key) (*result.Result, error) {
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

	return result.NewFromJSON(output), nil
}

func (u *UnityBridgeImpl) SetKeyValue(k *key.Key, value any,
	c result.Callback) error {
	if k.AccessType()&key.AccessTypeWrite == 0 {
		return fmt.Errorf("key %s is not writable", k)
	}

	v := struct {
		Value interface{} `json:"value"`
	}{
		value,
	}
	data, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	ev := event.NewFromTypeAndSubType(event.TypeSetValue, k.SubType())

	tag := u.tg.Next()

	u.m.Lock()

	u.callbackListener[tag] = c

	u.m.Unlock()

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

	u.m.Lock()

	u.callbackListener[tag] = c

	u.m.Unlock()

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

	tk := u.tg.Next()

	u.m.Lock()

	if _, ok := u.eventTypeListeners[t]; !ok {
		u.eventTypeListeners[t] = make(map[token.Token]event.TypeCallback)
	}

	u.eventTypeListeners[t][tk] = c

	u.m.Unlock()

	return tk, nil
}

func (u *UnityBridgeImpl) RemoveEventTypeListener(t event.Type,
	tk token.Token) error {
	if tk == 0 {
		return fmt.Errorf("token cannot be 0")
	}

	u.m.Lock()
	defer u.m.Unlock()

	if _, ok := u.eventTypeListeners[t]; !ok {
		return fmt.Errorf("no listeners registered for event type %s", t)
	}

	if _, ok := u.eventTypeListeners[t][tk]; !ok {
		return fmt.Errorf("no listener registered with token %d for event type %s",
			tk, t)
	}

	delete(u.eventTypeListeners[t], tk)

	if len(u.eventTypeListeners[t]) == 0 {
		delete(u.eventTypeListeners, t)
	}

	return nil
}

func (u *UnityBridgeImpl) Stop() error {
	u.m.Lock()

	if !u.started {
		u.m.Unlock()
		return fmt.Errorf("unity bridge not started")
	}

	u.started = false

	u.m.Unlock()

	for _, eventType := range event.AllTypes() {
		eventTypeCode := event.NewFromType(eventType).Code()
		u.uw.SetEventCallback(eventTypeCode, nil)
	}

	u.uw.Uninitialize()
	u.uw.Destroy()

	return nil
}

func (u *UnityBridgeImpl) handleOwnedEvents(e *event.Event, data []byte,
	tag uint64, dataType event.DataType) error {
	switch e.Type() {
	case event.TypeSetValue, event.TypePerformAction, event.TypeGetValue:
		u.notifyCallbacks(data, tag)
	case event.TypeStartListening:
		k, err := key.FromSubType(e.SubType())
		if err != nil {
			return err
		}
		u.notifyKeyListeners(k, data)
	}

	return nil
}

func (u *UnityBridgeImpl) eventCallback(eventCode uint64, data []byte, tag uint64) {
	e := event.NewFromCode(eventCode)

	dataType, tag := event.DataTypeFromTag(tag)

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	if e.Type() == event.TypeGetValue || e.Type() == event.TypeSetValue ||
		e.Type() == event.TypePerformAction || e.Type() == event.TypeStartListening {
		if err := u.handleOwnedEvents(e, data, tag, dataType); err != nil {
			u.l.Error("Error handling owned event", "event", e, "err", err)
		}

		// TODO(bga): We might want to not return here to allow other listeners
		// to also get notified for the event we own.
		return
	}

	// Call all registered event type listeners.
	u.notifyEventTypeListeners(e, data, dataType)
}

func (u *UnityBridgeImpl) notifyEventTypeListeners(e *event.Event,
	data []byte, dataType event.DataType) {
	u.m.RLock()

	if _, ok := u.eventTypeListeners[e.Type()]; ok {
		for _, c := range u.eventTypeListeners[e.Type()] {
			go c(data, dataType)
		}
	} else {
		u.l.Warn("No listeners registered for event type", "eventType",
			e.Type(), "event", e, "len(data)", len(data))
	}

	u.m.RUnlock()
}

func (u *UnityBridgeImpl) notifyKeyListeners(k *key.Key, data []byte) {
	u.m.RLock()

	if _, ok := u.keyListeners[k]; ok {
		for _, c := range u.keyListeners[k] {
			go c(result.NewFromJSON(data))
		}
	} else {
		u.l.Warn("No listeners registered for key", "key", k, "data", string(data))
	}

	u.m.RUnlock()
}

func (u *UnityBridgeImpl) notifyCallbacks(data []byte, tag uint64) {
	u.m.Lock()
	if c, ok := u.callbackListener[token.Token(tag)]; ok {
		go c(result.NewFromJSON(data))
		delete(u.callbackListener, token.Token(tag))
	} else {
		u.l.Error("No callback registered for tag", "tag", tag)
	}

	u.m.Unlock()
}
