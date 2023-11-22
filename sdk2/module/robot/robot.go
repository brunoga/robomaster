package robot

import (
	"log/slog"
	"sync/atomic"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
)

// Robot is the module that controls the robot. It provides methods to
// generally set it up.
type Robot struct {
	*internal.BaseModule

	functions atomic.Pointer[map[FunctionType]bool]
}

var _ module.Module = (*Robot)(nil)

// New creates a new Robot instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Robot, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("robot_module")

	return &Robot{
		BaseModule: internal.NewBaseModule(ub, l, "Robot",
			key.KeyRobomasterSystemConnection, nil),
	}, nil
}

type functionEnableInfo struct {
	ID     FunctionType `json:"id"`
	Enable bool         `json:"enable"`
}

type functionEnableParamValue struct {
	List []functionEnableInfo `json:"list"`
}

// EnableFunction enables or disables the given function. This keeps track of
// previoudly enabled/didabled functions and always sends the full list with
// the current status of all functions to the Unity Bridge.
func (r *Robot) EnableFunction(ft FunctionType, enable bool) error {
	info := functionEnableInfo{
		ID:     ft,
		Enable: enable,
	}

	param := functionEnableParamValue{
		List: []functionEnableInfo{info},
	}

	var newFunctions map[FunctionType]bool
	for {
		oldFunctionsPtr := r.functions.Load()
		oldFunctions := *oldFunctionsPtr
		newFunctions := make(map[FunctionType]bool, len(oldFunctions)+1)
		for k, v := range oldFunctions {
			newFunctions[k] = v
		}
		newFunctions[ft] = enable

		if r.functions.CompareAndSwap(oldFunctionsPtr, &newFunctions) {
			break
		}
	}

	for ft, enabled := range newFunctions {
		param.List = append(param.List, functionEnableInfo{
			ID:     ft,
			Enable: enabled,
		})
	}

	err := r.UB().PerformActionForKeySync(key.KeyRobomasterSystemFunctionEnable,
		param)
	if err != nil {
		return err
	}

	return nil
}
