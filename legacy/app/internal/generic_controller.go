package internal

import (
	"fmt"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge"
)

type GenericController struct {
	eventHandler bridge.EventHandler

	eventHandlerIndexMap map[unity.EventType]uint64
}

func NewGenericController(eventHandler bridge.EventHandler) *GenericController {
	return &GenericController{
		eventHandler,
		make(map[unity.EventType]uint64),
	}
}

func (g *GenericController) StartControllingEvent(eventType unity.EventType) error {
	b := bridge.Instance()

	_, ok := g.eventHandlerIndexMap[eventType]
	if ok {
		return fmt.Errorf("event type %s is already being controlled",
			unity.EventTypeName(eventType))
	}

	i, err := b.AddEventHandler(eventType, g.eventHandler)
	if err != nil {
		return err
	}

	g.eventHandlerIndexMap[eventType] = i

	return nil
}

func (g *GenericController) StopControllingEvent(eventType unity.EventType) error {
	b := bridge.Instance()

	i, ok := g.eventHandlerIndexMap[eventType]
	if !ok {
		return fmt.Errorf("event type %s is not being controlled",
			unity.EventTypeName(eventType))
	}

	err := b.RemoveEventHandler(eventType, i)
	if err != nil {
		return err
	}

	delete(g.eventHandlerIndexMap, eventType)

	return nil
}
