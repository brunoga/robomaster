package robot

import (
	"fmt"
	"log/slog"
	"sort"
	"sync/atomic"
	"time"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/connection"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/listener"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

// Robot is the module that controls the robot. It provides methods to
// generally set it up.
type Robot struct {
	*internal.BaseModule

	functions           atomic.Pointer[map[FunctionType]bool]
	workingDevices      atomic.Pointer[map[DeviceType]struct{}]
	batteryPowerPercent atomic.Pointer[uint8]

	workingDevicesRL      *listener.Listener
	batteryPowerPercentRL *listener.Listener
	actionStatusRL        *listener.Listener
}

var _ module.Module = (*Robot)(nil)

// New creates a new Robot instance.
func New(ub unitybridge.UnityBridge, l *logger.Logger,
	cm *connection.Connection) (*Robot, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("robot_module")

	rb := &Robot{}

	rb.BaseModule = internal.NewBaseModule(ub, l, "Robot",
		key.KeyRobomasterSystemConnection, func(r *result.Result) {
			if !r.Succeeded() {
				rb.Logger().Error(
					"Connection: Unexpected result.", "result", r)
				return
			}

			value, ok := r.Value().(*value.Bool)
			if !ok {
				rb.Logger().Error("Connection: Unexpected value.", "value",
					r.Value())
				return
			}

			if value.Value {
				// Connection is up. Start listeners.
				rb.Logger().Debug(
					"Connection: Connected. Starting listeners.")
				err := rb.workingDevicesRL.Start()
				if err != nil {
					rb.Logger().Error("Connection: Failed to start working "+
						"devices result listener.", "error", err)
				}

				err = rb.batteryPowerPercentRL.Start()
				if err != nil {
					rb.Logger().Error("Connection: Failed to start battery "+
						"power percent result listener", "error", err)
				}

				err = rb.actionStatusRL.Start()
				if err != nil {
					rb.Logger().Error("Connection: Failed to start action "+
						"status result listener", "error", err)
				}
			} else {
				// Connection is down. Stop listeners.
				rb.Logger().Debug(
					"Connection: Disconnected. Stopping listeners.")
				err := rb.workingDevicesRL.Stop()
				if err != nil {
					rb.Logger().Error("Connection: Failed to stop working "+
						"devices result listener.", "error", err)
				}

				err = rb.batteryPowerPercentRL.Stop()
				if err != nil {
					rb.Logger().Error("Connection: Failed to stop battery "+
						"power percent result listener", "error", err)
				}

				err = rb.actionStatusRL.Stop()
				if err != nil {
					rb.Logger().Error("Connection: Failed to stop action "+
						"status result listener", "error", err)
				}
			}
		}, cm)

	functions := make(map[FunctionType]bool)
	rb.functions.Store(&functions)

	wds := make(map[DeviceType]struct{})
	rb.workingDevices.Store(&wds)

	rb.workingDevicesRL = listener.New(ub, l,
		key.KeyRobomasterSystemWorkingDevices, func(res *result.Result) {
			rb.onWorkingDevices(res)
		})

	rb.batteryPowerPercentRL = listener.New(ub, l,
		key.KeyRobomasterBatteryPowerPercent, func(res *result.Result) {
			rb.onBatteryPowerPercent(res)
		})

	rb.actionStatusRL = listener.New(ub, l,
		key.KeyRobomasterSystemTaskStatus, func(res *result.Result) {
			rb.onActionStatus(res)
		})

	return rb, nil
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
	return r.workingDevicesRL.WaitForAnyResult(timeout) != nil
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

	return uint8(res.Value().(*value.Uint64).Value), nil
}

// SetSpeakerVolume sets the speaker volume.
func (r *Robot) SetSpeakerVolume(volume uint8) error {
	return r.UB().SetKeyValueSync(key.KeyRobomasterSystemSpeakerVolumn,
		&value.Uint64{Value: uint64(volume)})
}

// BatteryPowerPercent returns the current battery power percent.
func (r *Robot) BatteryPowerPercent() uint8 {
	return *r.batteryPowerPercent.Load()
}

// ChassisSpeedLevel returns the current chassis speed level.
func (r *Robot) ChassisSpeedLevel() (ChassisSpeedLevelType, error) {
	res, err := r.UB().GetKeyValueSync(key.KeyRobomasterSystemChassisSpeedLevel, true)
	if err != nil {
		return 0, err
	}

	return ChassisSpeedLevelType(res.Value().(*value.Uint64).Value - 1), nil
}

// SetChassisSpeedLevel sets the chassis speed level.
func (r *Robot) SetChassisSpeedLevel(speedLevel ChassisSpeedLevelType) error {
	if speedLevel >= ChassisSpeedLevelTypeCount {
		return fmt.Errorf("invalid chassis speed level: %d", speedLevel)
	}

	return r.UB().SetKeyValueSync(key.KeyRobomasterSystemChassisSpeedLevel,
		&value.Uint64{Value: uint64(speedLevel + 1)})
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
		r.Logger().Error("Unexpected battery power percent result.", "result",
			res)
		return
	}

	value, ok := res.Value().(*value.Uint64)
	if !ok {
		r.Logger().Error("Unexpected battery power percent value.", "value",
			res.Value())
		return
	}

	batterhyPowerPercent := uint8(value.Value)

	r.batteryPowerPercent.Store(&batterhyPowerPercent)
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
