package robot

import (
	"log/slog"
	"sort"
	"sync/atomic"
	"time"

	"github.com/brunoga/robomaster/sdk2/module"
	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/robomaster/sdk2/module/internal"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
	"github.com/brunoga/unitybridge/unity/result/value"
)

// Robot is the module that controls the robot. It provides methods to
// generally set it up.
type Robot struct {
	*internal.BaseModule

	functions           atomic.Pointer[map[FunctionType]bool]
	workingDevices      atomic.Pointer[map[DeviceType]struct{}]
	batteryPowerPercent atomic.Pointer[uint8]

	workingDevicesRL      *support.ResultListener
	batteryPowerPercentRL *support.ResultListener
	actionStatusRL        *support.ResultListener
}

var _ module.Module = (*Robot)(nil)

// New creates a new Robot instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection) (*Robot, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("robot_module")

	r := &Robot{
		BaseModule: internal.NewBaseModule(ub, l, "Robot",
			key.KeyRobomasterSystemConnection, nil, cm),
	}

	functions := make(map[FunctionType]bool)
	r.functions.Store(&functions)

	wds := make(map[DeviceType]struct{})
	r.workingDevices.Store(&wds)

	r.workingDevicesRL = support.NewResultListener(ub, l,
		key.KeyRobomasterSystemWorkingDevices, func(res *result.Result) {
			r.onWorkingDevices(res)
		})

	r.batteryPowerPercentRL = support.NewResultListener(ub, l,
		key.KeyRobomasterBatteryPowerPercent, func(res *result.Result) {
			r.onBatteryPowerPercent(res)
		})

	r.actionStatusRL = support.NewResultListener(ub, l,
		key.KeyRobomasterSystemTaskStatus, func(res *result.Result) {
			r.onActionStatus(res)
		})

	return r, nil
}

// Start starts the Robot module.
func (r *Robot) Start() error {
	err := r.BaseModule.Start()
	if err != nil {
		return err
	}

	err = r.workingDevicesRL.Start()
	if err != nil {
		return err
	}

	err = r.batteryPowerPercentRL.Start()
	if err != nil {
		return err
	}

	return r.actionStatusRL.Start()
}

// EnableFunction enables or disables the given function. This keeps track of
// previoudly enabled/didabled functions and always sends the full list with
// the current status of all functions to the Unity Bridge.
func (r *Robot) EnableFunction(ft FunctionType, enable bool) error {
	var newFunctions map[FunctionType]bool
	for {
		oldFunctionsPtr := r.functions.Load()
		oldFunctions := *oldFunctionsPtr

		newFunctions = make(map[FunctionType]bool, len(oldFunctions)+1)
		for k, v := range oldFunctions {
			newFunctions[k] = v
		}
		newFunctions[ft] = enable

		if r.functions.CompareAndSwap(oldFunctionsPtr, &newFunctions) {
			break
		}
	}

	v := &value.FunctionEnable{
		List: []value.FunctionEnableInfo{},
	}

	for ft, enabled := range newFunctions {
		v.List = append(v.List, value.FunctionEnableInfo{
			ID:     uint8(ft),
			Enable: enabled,
		})
	}

	err := r.UB().PerformActionForKeySync(key.KeyRobomasterSystemFunctionEnable,
		v)
	if err != nil {
		return err
	}

	return nil
}

// WaitForDevices waits for the robot to report devices as working. Returns
// true if an actual result was obtained and false otherwise (i.e. timeout).
func (r *Robot) WaitForDevices(timeout time.Duration) bool {
	// Just wait for a result to be available.
	if r.workingDevicesRL.Result() != nil {
		return true
	}

	return r.workingDevicesRL.WaitForNewResult(timeout) != nil
}

// HasFunction returns true if the given device is connected to the robot and
// is reported as working.
func (r *Robot) HasDevice(device DeviceType) bool {
	wds := r.workingDevices.Load()
	_, ok := (*wds)[device]

	return ok
}

// Devices returns the list of working devices (i.e. devices connected to the
// robot and that are reported as working).
func (r *Robot) Devices() []DeviceType {
	wds := *r.workingDevices.Load()

	ks := make([]DeviceType, 0, len(wds))
	for k := range wds {
		ks = append(ks, k)
	}

	sort.SliceStable(ks, func(i, j int) bool {
		return ks[i] < ks[j]
	})

	return ks
}

// SpeakerVolume returns the current speaker volume.
func (r *Robot) SpeakerVolume() (uint8, error) {
	res, err := r.UB().GetKeyValueSync(key.KeyRobomasterSystemSpeakerVolumn, true)
	if err != nil {
		return 0, err
	}

	return uint8(res.Value().(float64)), nil
}

// SetSpeakerVolume sets the speaker volume.
func (r *Robot) SetSpeakerVolume(volume uint8) error {
	return r.UB().SetKeyValueSync(key.KeyRobomasterSystemSpeakerVolumn,
		volume)
}

// BatteryPowerPercent returns the current battery power percent.
func (r *Robot) BatteryPowerPercent() uint8 {
	return *r.batteryPowerPercent.Load()
}

// Stop stops the Robot module.
func (r *Robot) Stop() error {
	err := r.actionStatusRL.Stop()
	if err != nil {
		return err
	}

	err = r.batteryPowerPercentRL.Stop()
	if err != nil {
		return err
	}

	err = r.workingDevicesRL.Stop()
	if err != nil {
		return err
	}

	return r.BaseModule.Stop()
}

func (r *Robot) onWorkingDevices(res *result.Result) {
	if res == nil || !res.Succeeded() {
		return
	}

	// 2 or more updates at the same time are *VERY* unlikely. If they happen,
	// we just accept whatever ordering the Store bellow gives us.
	oldWds := *r.workingDevices.Load()
	newWds := wdsListToWds(res.Value().(*value.List[uint16]).List)
	r.workingDevices.Store(&newWds)

	removed, added := r.checkDiff(oldWds, newWds)

	for device := range removed {
		r.Logger().Warn("Device removed", "device", device)
	}

	for device := range added {
		r.Logger().Warn("Device added", "device", device)
	}

	// TODO(bga): If we ever associate device types with modules, we might
	// want to automatically enable/disable them here.
}

func (r *Robot) onBatteryPowerPercent(res *result.Result) {
	if res == nil || !res.Succeeded() {
		return
	}

	v := uint8(res.Value().(*value.Uint64).Value)

	r.batteryPowerPercent.Store(&v)
}

func (r *Robot) onActionStatus(res *result.Result) {
	// Just log for now.
	//
	// TODO(bga): Implement me.
	r.Logger().Debug("Action status", "result", res)
}

func (r *Robot) checkDiff(oldWds, newWds map[DeviceType]struct{}) (
	map[DeviceType]struct{}, map[DeviceType]struct{}) {
	removed := make(map[DeviceType]struct{})
	for k := range oldWds {
		if _, ok := newWds[k]; !ok {
			removed[k] = struct{}{}
		}
	}

	added := make(map[DeviceType]struct{})
	for k := range newWds {
		if _, ok := oldWds[k]; !ok {
			added[k] = struct{}{}
		}
	}

	return removed, added
}

func wdsListToWds(wdsList []uint16) map[DeviceType]struct{} {
	wds := make(map[DeviceType]struct{}, len(wdsList))
	for _, wd := range wdsList {
		wds[DeviceType(wd)] = struct{}{}
	}

	return wds
}
