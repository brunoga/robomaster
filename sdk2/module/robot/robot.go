package robot

import (
	"encoding/json"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
)

type Robot struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger
}

var _ module.Module = (*Robot)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Robot, error) {
	return &Robot{
		ub: ub,
		l:  l,
	}, nil
}

func (r *Robot) Start() error {
	return nil
}

type functionEnableInfo struct {
	ID     FunctionType `json:"id"`
	Enable bool         `json:"enable"`
}

type functionEnableParamValue struct {
	List []functionEnableInfo `json:"list"`
}

func (r *Robot) EnableFunction(function FunctionType, enable bool) error {
	info := functionEnableInfo{
		ID:     function,
		Enable: enable,
	}

	param := functionEnableParamValue{
		List: []functionEnableInfo{info},
	}

	data, err := json.Marshal(param)
	if err != nil {
		return err
	}

	err = r.ub.PerformActionForKeySync(key.KeyRobomasterSystemFunctionEnable, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Robot) Stop() error {
	return nil
}

func (r *Robot) String() string {
	return "Robot"
}
