package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/brunoga/robomaster/unitybridge/support/logger"
	"github.com/brunoga/robomaster/unitybridge/support/token"
	"github.com/brunoga/robomaster/unitybridge/unity/event"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
	"github.com/brunoga/robomaster/unitybridge/wrapper"
)

var (
	voidType = reflect.TypeOf(&value.Void{})
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

	l = l.WithGroup("unity_bridge")

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

func (u *UnityBridgeImpl) Start() (err error) {
	endTrace := u.l.Trace("Start")
	defer func() {
		endTrace("error", err)
	}()

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
	immediate bool) (t token.Token, err error) {
	endTrace := u.l.Trace("AddKeyListener", "key", k, "callback", c,
		"immediate", immediate)
	defer func() {
		endTrace("token", t, "error", err)
	}()

	if k.AccessType()&key.AccessTypeRead == 0 {
		return 0, fmt.Errorf("key %s is not readable", k)
	}

	if c == nil {
		return 0, fmt.Errorf("callback cannot be nil")
	}

	t = u.tg.Next()

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
	if err != nil {
		return 0, err
	}

	if r == nil {
		// Nil error with nil result means no cached value. Just return the
		// token.
		return t, nil
	}

	go c(r)

	return t, nil
}

func (u *UnityBridgeImpl) RemoveKeyListener(k *key.Key,
	token token.Token) (err error) {
	endTrace := u.l.Trace("RemoveKeyListener", "key", k, "token", token)
	defer func() {
		endTrace("error", err)
	}()

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

func (u *UnityBridgeImpl) GetKeyValue(k *key.Key, c result.Callback) (err error) {
	endTrace := u.l.Trace("GetKeyValue", "key", k, "callback", c)
	defer func() {
		endTrace("error", err)
	}()

	defer u.l.Trace("GetKeyValue", []any{
		"error", err,
	}, "key", k, "callback", c)

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

func (u *UnityBridgeImpl) GetKeyValueSync(k *key.Key,
	useCache bool) (r *result.Result, err error) {
	endTrace := u.l.Trace("GetKeyValueSync", "key", k, "useCache", useCache)
	defer func() {
		endTrace("result", r, "error", err)
	}()

	if useCache {
		r, err = u.GetCachedKeyValue(k)
		if err == nil && r != nil && r.Succeeded() {
			// Have a valid cached result.
			return r, err
		}
	}

	done := make(chan struct{})

	err = u.GetKeyValue(k, func(r2 *result.Result) {
		if r2.ErrorCode() != 0 {
			err = fmt.Errorf("error getting value for key %s: %s", k,
				r2.ErrorDesc())
		} else {
			r = r2
		}

		close(done)
	})

	if err != nil {
		return nil, err
	}

	select {
	case <-done:
		return r, err
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout getting value for key %s", k)
	}
}

func (u *UnityBridgeImpl) GetCachedKeyValue(k *key.Key) (r *result.Result, err error) {
	endTrace := u.l.Trace("GetCachedKeyValue", "key", k)
	defer func() {
		endTrace("result", r, "error", err)
	}()

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

	if len(output) == 0 {
		// No cached value.
		return nil, nil
	}

	return result.NewFromJSON(output), nil
}

func (u *UnityBridgeImpl) SetKeyValue(k *key.Key, value any,
	c result.Callback) (err error) {
	endTrace := u.l.Trace("SetKeyValue", "key", k, "value", value, "callback", c)
	defer func() {
		endTrace("error", err)
	}()

	if k.AccessType()&key.AccessTypeWrite == 0 {
		return fmt.Errorf("key %s is not writable", k)
	}

	expectedKeyValue := k.ResultValue()

	if reflect.TypeOf(value) != reflect.TypeOf(expectedKeyValue) {
		return fmt.Errorf("value type %s does not match expected key %s type "+
			"%s", reflect.TypeOf(value), k, reflect.TypeOf(expectedKeyValue))
	}

	data, err := json.Marshal(value)
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

func (u *UnityBridgeImpl) SetKeyValueSync(k *key.Key, value any) (err error) {
	endTrace := u.l.Trace("SetKeyValueSync", "key", k, "value", value)
	defer func() {
		endTrace("error", err)
	}()

	u.l.Trace("SetKeyValueSync", "key", k, "value", value)

	done := make(chan struct{})

	err = u.SetKeyValue(k, value, func(r *result.Result) {
		if r.ErrorCode() != 0 {
			err = fmt.Errorf("error setting value for key %s: %s", k, r.ErrorDesc())
		}

		close(done)
	})

	if err != nil {
		return err
	}

	select {
	case <-done:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout setting value for key %s", k)
	}
}

func (u *UnityBridgeImpl) PerformActionForKey(k *key.Key, value any,
	c result.Callback) (err error) {
	endTrace := u.l.Trace("PerformActionForKey", "key", k, "value", value,
		"callback", c)
	defer func() {
		endTrace("error", err)
	}()

	if k.AccessType()&key.AccessTypeAction == 0 {
		return fmt.Errorf("key %s is not an action", k)
	}

	expectedKeyValue := k.ResultValue()
	expectedType := reflect.TypeOf(expectedKeyValue)
	actualType := reflect.TypeOf(value)

	if expectedType == voidType {
		if value != nil {
			return fmt.Errorf("key %s is void type but value is not nil", k)
		}
	} else if actualType != expectedType {
		return fmt.Errorf("value type %s does not match expected key %s type "+
			"%s", actualType, k, expectedType)
	}

	var data []byte

	if value != nil {
		data, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	ev := event.NewFromTypeAndSubType(event.TypePerformAction, k.SubType())

	tag := u.tg.Next()

	u.m.Lock()

	u.callbackListener[tag] = c

	u.m.Unlock()

	if value != nil {
		u.uw.SendEventWithString(ev.Code(), string(data), uint64(tag))
	} else {
		u.uw.SendEvent(ev.Code(), nil, uint64(tag))
	}

	return nil
}

func (u *UnityBridgeImpl) PerformActionForKeySync(k *key.Key,
	value any) (err error) {
	endTrace := u.l.Trace("PerformActionForKeySync", "key", k, "value", value)
	defer func() {
		endTrace("error", err)
	}()

	u.l.Trace("PerformActionForKeySync", "key", k, "value", value)

	done := make(chan struct{})

	err = u.PerformActionForKey(k, value, func(r *result.Result) {
		if r.ErrorCode() != 0 {
			err = fmt.Errorf("error performing action for key %s: %s", k,
				r.ErrorDesc())
		}

		close(done)
	})

	if err != nil {
		return err
	}

	select {
	case <-done:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout performing action for key %s", k)
	}
}

func (u *UnityBridgeImpl) DirectSendKeyValue(k *key.Key,
	value uint64) (err error) {
	endTrace := u.l.Trace("DirectSendKeyValue", "key", k, "value", value)
	defer func() {
		endTrace("error", err)
	}()

	// This always sends the event as an action but it never checks for the
	// action access type. This is intentional as we want to be able to do this
	// for any key regardless of access type.
	ev := event.NewFromTypeAndSubType(event.TypePerformAction, k.SubType())

	u.uw.SendEventWithNumber(ev.Code(), value, 0)

	return nil
}

func (u *UnityBridgeImpl) SendEvent(ev *event.Event) (err error) {
	endTrace := u.l.Trace("SendEvent", "event", ev)
	defer func() {
		endTrace("error", err)
	}()

	u.uw.SendEvent(ev.Code(), nil, 0)

	return nil
}

func (u *UnityBridgeImpl) SendEventWithString(ev *event.Event,
	data string) (err error) {
	endTrace := u.l.Trace("SendEventWithString", "event", ev, "data", data)
	defer func() {
		endTrace("error", err)
	}()

	u.uw.SendEventWithString(ev.Code(), data, 0)

	return nil
}

func (u *UnityBridgeImpl) SendEventWithUint64(ev *event.Event,
	data uint64) (err error) {
	endTrace := u.l.Trace("SendEventWithUint64", "event", ev, "data", data)
	defer func() {
		endTrace("error", err)
	}()

	u.uw.SendEventWithNumber(ev.Code(), data, 0)

	return nil
}

func (u *UnityBridgeImpl) AddEventTypeListener(et event.Type,
	c event.TypeCallback) (t token.Token, err error) {
	endTrace := u.l.Trace("AddEventTypeListener", "eventType", et, "callback",
		c)
	defer func() {
		endTrace("token", t, "error", err)
	}()

	if c == nil {
		return 0, fmt.Errorf("callback cannot be nil")
	}

	tk := u.tg.Next()

	u.m.Lock()

	if _, ok := u.eventTypeListeners[et]; !ok {
		u.eventTypeListeners[et] = make(map[token.Token]event.TypeCallback)
	}

	u.eventTypeListeners[et][tk] = c

	u.m.Unlock()

	return tk, nil
}

func (u *UnityBridgeImpl) RemoveEventTypeListener(t event.Type,
	tk token.Token) (err error) {
	endTrace := u.l.Trace("RemoveEventTypeListener", "eventType", t, "token", tk)
	defer func() {
		endTrace("error", err)
	}()

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

func (u *UnityBridgeImpl) Stop() (err error) {
	endTrace := u.l.Trace("Stop")
	defer func() {
		endTrace("error", err)
	}()

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
	tag uint64, dataType event.DataType) (err error) {
	endTrace := u.l.Trace("handleOwnedEvents", "event", e, "data", data,
		"tag", tag, "dataType", dataType)
	defer func() {
		endTrace("error", err)
	}()

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
	endTrace := u.l.Trace("eventCallback", "eventCode", eventCode, "data", data,
		"tag", tag)
	defer endTrace()

	e := event.NewFromCode(eventCode)

	dataType, tag := event.DataTypeFromTag(tag)

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
	endTrace := u.l.Trace("notifyEventTypeListeners", "event", e, "data", data,
		"dataType", dataType)
	defer endTrace()

	u.m.RLock()

	if _, ok := u.eventTypeListeners[e.Type()]; ok {
		for _, c := range u.eventTypeListeners[e.Type()] {
			go c(data, dataType)
		}
	} else {
		u.l.Warn("No listeners registered for event type", "eventType",
			e.Type(), "event", e, "data", data)
	}

	u.m.RUnlock()
}

func (u *UnityBridgeImpl) notifyKeyListeners(k *key.Key, data []byte) {
	endTrace := u.l.Trace("notifyKeyListeners", "key", k, "data", data)
	defer endTrace()

	u.m.RLock()

	r := result.NewFromJSON(data)

	if _, ok := u.keyListeners[k]; ok {
		for _, c := range u.keyListeners[k] {
			go c(r)
		}
	} else {
		u.l.Warn("No listeners registered for key", "key", k, "data", string(data))
	}

	u.m.RUnlock()
}

func (u *UnityBridgeImpl) notifyCallbacks(data []byte, tag uint64) {
	endTrace := u.l.Trace("notifyCallbacks", "data", data, "tag", tag)
	defer endTrace()

	u.m.Lock()
	if c, ok := u.callbackListener[token.Token(tag)]; ok {
		if c != nil {
			go c(result.NewFromJSON(data))
		}
		delete(u.callbackListener, token.Token(tag))
	} else {
		u.l.Error("No callback registered for tag", "tag", tag)
	}

	u.m.Unlock()
}
