package internal

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge"
	"git.bug-br.org.br/bga/robomasters1/app/internal/support/callbacks"
)

type ResultHandler func(result *dji.Result, wg *sync.WaitGroup)

type CommandController struct {
	*GenericController

	startListeningCallbacks *callbacks.Callbacks
	performActionCallbacks  *callbacks.Callbacks
	setValueForKeyCallbacks *callbacks.Callbacks
}

func NewCommandController() (*CommandController, error) {
	startListeningCallbacks := callbacks.New(
		"CommandController/StartListening",
		func(key callbacks.Key) error {
			return bridge.Instance().SendEvent(
				unity.NewEventWithSubType(
					unity.EventTypeStartListening,
					uint64(key)))
		},
		func(key callbacks.Key) error {
			return bridge.Instance().SendEvent(
				unity.NewEventWithSubType(
					unity.EventTypeStopListening,
					uint64(key)))
		},
	)

	performActionCallbacks := callbacks.New(
		"CommandController/PerformAction", nil, nil)

	setValueForKeyCallbacks := callbacks.New(
		"CommandController/SetValueForKey", nil, nil)

	cc := &CommandController{
		startListeningCallbacks: startListeningCallbacks,
		performActionCallbacks:  performActionCallbacks,
		setValueForKeyCallbacks: setValueForKeyCallbacks,
	}

	cc.GenericController = NewGenericController(cc.HandleEvent)

	var err error
	err = cc.StartControllingEvent(unity.EventTypeGetValue)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(unity.EventTypeSetValue)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(unity.EventTypePerformAction)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(unity.EventTypeStartListening)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *CommandController) StartListening(key dji.Key,
	resultHandler ResultHandler) (uint64, error) {
	if key < 1 || key >= dji.KeysCount {
		return 0, fmt.Errorf("invalid key")
	}
	if resultHandler == nil {
		return 0, fmt.Errorf("eventHandler must not be nil")
	}
	if (key.AccessType() & dji.KeyAccessTypeRead) == 0 {
		return 0, fmt.Errorf("key is not readable")
	}

	tag, err := c.startListeningCallbacks.AddContinuous(callbacks.Key(
		key.Value()), resultHandler)

	return uint64(tag), err
}

func (c *CommandController) StopListening(key dji.Key, tag uint64) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}

	return c.startListeningCallbacks.Remove(callbacks.Key(key),
		callbacks.Tag(tag))
}

func (c *CommandController) PerformAction(key dji.Key, param interface{},
	resultHandler ResultHandler) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}
	if (key.AccessType() & dji.KeyAccessTypeAction) == 0 {
		return fmt.Errorf("key can not be acted upon")
	}

	var err error
	tag := callbacks.Tag(0)

	if resultHandler != nil {
		tag, err = c.performActionCallbacks.AddSingleShot(
			callbacks.Key(key.Value()), resultHandler)
		if err != nil {
			return err
		}
	}

	var data []byte
	if param != nil {
		data, err = json.Marshal(param)
		if err != nil {
			return err
		}
	}

	bridge.Instance().SendEvent(unity.NewEventWithSubType(
		unity.EventTypePerformAction, uint64(key.Value())), data,
		uint64(tag))

	return nil
}

func (c *CommandController) SetValueForKey(key dji.Key, param interface{},
	resultHandler ResultHandler) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}
	if (key.AccessType() & dji.KeyAccessTypeWrite) == 0 {
		return fmt.Errorf("key can not be written to")
	}

	var err error
	tag := callbacks.Tag(0)

	if resultHandler != nil {
		tag, err = c.setValueForKeyCallbacks.AddSingleShot(
			callbacks.Key(key.Value()), resultHandler)
		if err != nil {
			return err
		}
	}

	var data []byte
	if param != nil {
		data, err = json.Marshal(param)
		if err != nil {
			return err
		}
	}

	bridge.Instance().SendEvent(unity.NewEventWithSubType(
		unity.EventTypeSetValue, uint64(key.Value())), string(data),
		uint64(tag))

	return nil
}

func (c *CommandController) DirectSendValue(key dji.Key, value uint64) {
	bridge.Instance().SendEvent(unity.NewEventWithSubType(
		unity.EventTypePerformAction, uint64(key.Value())), value)
}

func (c *CommandController) HandleEvent(event *unity.Event, data []byte,
	tag uint64, wg *sync.WaitGroup) {
	var value interface{}

	// TODO(bga): Apparently, the unity bridge reserves the upper 8 bits
	//  for reporting back type information. Double check this.
	dataType := tag >> 56
	switch dataType {
	case 0:
		value = string(data)
	case 1:
		value = binary.LittleEndian.Uint64(data)
	default:
		// Apparently only string and uint64 types are supported
		// currently.
		panic(fmt.Sprintf("Unexpected data type: %d.\n", dataType))
	}

	// See above.
	adjustedTag := tag & 0xffffffffffffff

	switch event.Type() {
	case unity.EventTypeStartListening:
		c.handleStartListening(event.SubType(), value)
	case unity.EventTypePerformAction:
		c.handlePerformAction(event.SubType(), value, adjustedTag)
	case unity.EventTypeSetValue:
		c.handleSetValue(event.SubType(), value, adjustedTag)
	default:
		log.Printf("Unsupported event %s. Value:%v. Tag:%d\n",
			unity.EventTypeName(event.Type()), value, tag)
	}

	wg.Done()
}

func (c *CommandController) handleStartListening(key uint64,
	value interface{}) {
	stringValue, ok := value.(string)
	if !ok {
		panic("unexpected non-string value")
	}

	cbs, err := c.startListeningCallbacks.CallbacksForKey(
		callbacks.Key(key))
	if err != nil {
		log.Printf("Error looking up callbacks for key %d: %s\n", key,
			err)
		return
	}

	result := dji.NewResultFromJSON([]byte(stringValue))

	var wg sync.WaitGroup
	for _, cb := range cbs {
		wg.Add(1)
		go cb.(ResultHandler)(result, &wg)
	}

	wg.Wait()
}

func (c *CommandController) handlePerformAction(key uint64,
	value interface{}, tag uint64) {
	stringValue, ok := value.(string)
	if !ok {
		panic("unexpected non-string value")
	}

	cb, err := c.startListeningCallbacks.Callback(callbacks.Key(key),
		callbacks.Tag(tag))
	if err != nil {
		// No callback set. That is fine.
		return
	}

	result := dji.NewResultFromJSON([]byte(stringValue))

	cb.(ResultHandler)(result, nil)
}

func (c *CommandController) handleSetValue(key uint64,
	value interface{}, tag uint64) {
	stringValue, ok := value.(string)
	if !ok {
		panic("unexpected non-string value")
	}

	cb, err := c.startListeningCallbacks.Callback(callbacks.Key(key),
		callbacks.Tag(tag))
	if err != nil {
		// No callback set. That is fine.
		return
	}

	result := dji.NewResultFromJSON([]byte(stringValue))

	cb.(ResultHandler)(result, nil)
}

func (c *CommandController) Teardown() error {
	var err error
	err = c.StopControllingEvent(unity.EventTypeGetValue)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(unity.EventTypeSetValue)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(unity.EventTypePerformAction)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(unity.EventTypeStartListening)
	if err != nil {
		return err
	}

	// TODO(bga): Also clean up ResultHandler callbacks before returning.

	return nil
}
