package connection

import (
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/finder"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
)

const (
	subTypeConnectionOpen = iota
	subTypeConnectionClose
	subTypeConnectionSetIP
	subTypeConnectionSetPort
)

// Connection provides support for managing the connection to the robot.
type Connection struct {
	ub    unitybridge.UnityBridge
	l     *logger.Logger
	appID uint64

	f                     *finder.Finder
	connectionStatusToken token.Token

	connRL *support.ResultListener
}

var _ module.Module = (*Connection)(nil)

// New creates a new Connection instance with the given UnityBridge instance and
// logger.
func New(ub unitybridge.UnityBridge,
	l *logger.Logger, appID uint64) (*Connection, error) {

	cm := &Connection{
		ub:     ub,
		l:      l,
		f:      finder.New(appID, l),
		appID:  appID,
		connRL: support.NewResultListener(ub, l, key.KeyAirLinkConnection),
	}

	return cm, nil
}

// Start starts the connection module. It will try to find a robot broadcasting
// in the network and connect to it.
func (cm *Connection) Start() error {
	err := cm.connRL.Start(nil)
	if err != nil {
		return err
	}

	b, err := cm.f.Find(30 * time.Second)
	if err != nil {
		return err
	}

	cm.f.SendACK(b.SourceIp(), b.AppId())

	e := event.NewFromType(event.TypeConnection)

	e.ResetSubType(subTypeConnectionClose)
	err = cm.ub.SendEvent(e)
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionSetIP)
	err = cm.ub.SendEventWithString(e, b.SourceIp().String())
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionSetPort)
	err = cm.ub.SendEventWithUint64(e, 10607)
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionOpen)
	err = cm.ub.SendEvent(e)
	if err != nil {
		return err
	}

	return nil
}

func (cm *Connection) WaitForConnection() bool {
	connected, ok := cm.connRL.Result().Value().(bool)
	if ok && connected {
		return true
	}

	return cm.connRL.WaitForNewResult(5 * time.Second).Value().(bool)
}

// Stop stops the connection module.
func (cm *Connection) Stop() error {
	e := event.NewFromType(event.TypeConnection)

	e.ResetSubType(subTypeConnectionClose)
	err := cm.ub.SendEvent(e)
	if err != nil {
		return err
	}

	err = cm.ub.RemoveKeyListener(key.KeyAirLinkConnection,
		cm.connectionStatusToken)
	if err != nil {
		return err
	}

	return cm.connRL.Stop()
}

// Connected returns true if the connection to the robot is established.
func (cm *Connection) Connected() bool {
	connected, ok := cm.connRL.Result().Value().(bool)
	if !ok {
		return false
	}

	return connected
}

// String returns a string representation of the Connection module.
func (cm *Connection) String() string {
	return "Connection"
}
