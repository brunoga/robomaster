package robot

import (
	"log/slog"
	"sync/atomic"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

// Robot is the module that controls the robot. It provides methods to
// generally set it up.
type Robot struct {
	*internal.BaseModule

	wdToken token.Token

	functions atomic.Pointer[map[FunctionType]bool]
	wds       atomic.Pointer[[]DeviceType]
}

var _ module.Module = (*Robot)(nil)

// New creates a new Robot instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Robot, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("robot_module")

	r := &Robot{
		BaseModule: internal.NewBaseModule(ub, l, "Robot",
			key.KeyRobomasterSystemConnection, nil),
	}

	functions := make(map[FunctionType]bool)
	r.functions.Store(&functions)

	wds := make([]DeviceType, 0)
	r.wds.Store(&wds)

	return r, nil
}

type functionEnableInfo struct {
	ID     FunctionType `json:"id"`
	Enable bool         `json:"enable"`
}

type functionEnableParamValue struct {
	List []functionEnableInfo `json:"list"`
}

func (r *Robot) Start() error {
	err := r.BaseModule.Start()
	if err != nil {
		return err
	}

	r.wdToken, err = r.UB().AddKeyListener(key.KeyRobomasterSystemWorkingDevices,
		r.onWorkingDevices, false)
	if err != nil {
		return err
	}

	return nil
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

// HasFunction returns true if the given device is connected to the robot and
// is reported as working.
func (r *Robot) HasDevice(device DeviceType) bool {
	wds := r.wds.Load()
	for _, wd := range *wds {
		if wd == device {
			return true
		}
	}

	return false
}

// Devices returns the list of working devices (i.e. devices connected to the
// robot and that are reported as working).
func (r *Robot) Devices() []DeviceType {
	return *r.wds.Load()
}

func (r *Robot) onWorkingDevices(res *result.Result) {
	if res == nil || !res.Succeeded() {
		return
	}

	// 2 or more updates at the same time are *VERY* unlikely. If they happen,
	// we just accept whatever ordering the Store bellow gives us.
	newWds := wdsListToWds(res.Value().([]interface{}))
	r.wds.Store(&newWds)
	r.Logger().Debug("working devices updated", "working_devices", newWds)
}

func wdsListToWds(wdsList []interface{}) []DeviceType {
	wds := make([]DeviceType, 0, len(wdsList)+1)
	for _, wd := range wdsList {
		wdStr, ok := wd.(float64)
		if !ok {
			continue
		}
		deviceType := DeviceType(wdStr)

		wds = append(wds, deviceType)

		if deviceType == DeviceTypeWaterGun {
			// Apparently the water gun implies infrared gun. This is based on
			// the fact that a robot with a working infrared gun is not
			// reporting it back so we kinda cheat here.
			//
			// TODO(bga): Double check this.
			wds = append(wds, DeviceTypeInfraredGun)
		}
	}

	return wds
}
