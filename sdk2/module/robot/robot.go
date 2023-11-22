package robot

import (
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
)

type Robot struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger

	connRL *support.ResultListener
}

var _ module.Module = (*Robot)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Robot, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("robot_module")

	r := &Robot{
		ub: ub,
		l:  l,
		connRL: support.NewResultListener(ub, l,
			key.KeyRobomasterSystemConnection, nil),
	}

	return r, nil
}

func (r *Robot) Start() error {
	return r.connRL.Start()
}

type functionEnableInfo struct {
	ID     FunctionType `json:"id"`
	Enable bool         `json:"enable"`
}

type functionEnableParamValue struct {
	List []functionEnableInfo `json:"list"`
}

func (r *Robot) Connected() bool {
	connected, ok := r.connRL.Result().Value().(bool)
	if !ok {
		return false
	}

	return connected
}

func (r *Robot) WaitForConnection(timeout time.Duration) bool {
	connected, ok := r.connRL.Result().Value().(bool)
	if ok && connected {
		return true
	}

	return r.connRL.WaitForNewResult(5 * time.Second).Value().(bool)
}

func (r *Robot) EnableFunction(function FunctionType, enable bool) error {
	info := functionEnableInfo{
		ID:     function,
		Enable: enable,
	}

	param := functionEnableParamValue{
		List: []functionEnableInfo{info},
	}

	err := r.ub.PerformActionForKeySync(key.KeyRobomasterSystemFunctionEnable,
		param)
	if err != nil {
		return err
	}

	return nil
}

func (r *Robot) Stop() error {
	return r.connRL.Stop()
}

func (r *Robot) String() string {
	return "Robot"
}
