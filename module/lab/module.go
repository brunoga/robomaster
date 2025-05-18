package lab

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/internal"
	"github.com/brunoga/robomaster/module/lab/types"
	"github.com/brunoga/robomaster/support/dsp"
	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge"
	"github.com/brunoga/robomaster/unitybridge/unity/key"
	"github.com/brunoga/robomaster/unitybridge/unity/result"
	"github.com/brunoga/robomaster/unitybridge/unity/result/value"
)

type Module struct {
	*internal.BaseModule

	uploadWg sync.WaitGroup

	m            sync.Mutex
	currentState types.ScratchState
}

var _ module.Module = (*Module)(nil)

func New(ub unitybridge.UnityBridge, l *logger.Logger) (*Module, error) {
	if l == nil {
		l = logger.New(slog.LevelError)
	}

	l = l.WithGroup("lab_module")

	m := &Module{}

	bm, err := internal.NewBaseModule(ub, l, "lab", key.KeyRobomasterSystemConnection, func(r *result.Result) {
		if !r.Succeeded() {
			l.Error("Lab Connection: Unsuccessfull result.", "result", r)
			return
		}

		connectedValue, ok := r.Value().(*value.Bool)

		if !ok {
			l.Error("Lab Connection: Unexpected value.", "value", r.Value())
			return
		}

		if connectedValue.Value {
			l.Debug("Lab Connected.")
		} else {
			l.Debug("Lab Disconnected.")
		}
	})
	if err != nil {
		return nil, err
	}

	m.BaseModule = bm

	return m, nil
}

func (m *Module) Start() error {
	var err error

	err = m.AddKeyListener(key.KeyRobomasterSystemScratchCallback, func(r *result.Result) {
		m.Logger().Error("Got scratch callback result.", "result", r)
		m.uploadWg.Done()
	}, false)
	if err != nil {
		return err
	}

	err = m.AddKeyListener(key.KeyRobomasterSystemScratchExecuteState, func(r *result.Result) {
		m.Logger().Error("Got scratch execute state result.", "result", r)
	}, false)
	if err != nil {
		return err
	}

	err = m.AddKeyListener(key.KeyRobomasterSystemScratchErrorInfo, func(r *result.Result) {
		m.Logger().Trace("Got scratch error info result.", "result", r)

		result := r.Value().(*value.ScratchErrorInfo)
		fmt.Println(result.Key, result.ErrorLine, result.ErrorMsg)
	}, false)
	if err != nil {
		return err
	}

	err = m.AddKeyListener(key.KeyRobomasterSystemScratchState, func(r *result.Result) {
		m.Logger().Trace("Got scratch state result.", "result", r)

		result := r.Value().(*value.ScratchState)
		fmt.Println("Key:", result.Key, "BlockID:", result.BlockID, "BlockType:", result.BlockType, "BlockExecuteState:", result.BlockExecuteState, "BlockExecutePercentage:", result.BlockExecutePercentage, "State:", result.State)
	}, false)
	if err != nil {
		return err
	}

	err = m.AddKeyListener(key.KeyRobomasterSystemScratchOutputInfo, func(r *result.Result) {
		m.Logger().Error("Got scratch output info result.", "result", r)

		result := r.Value().(*value.ScratchOutputInfo)

		fmt.Println(result.InfoLevel, result.InfoType, result.InfoContent)
	}, false)
	if err != nil {
		return err
	}

	return m.BaseModule.Start()
}

func (m *Module) RunDSP(f *dsp.File) (string, string, error) {
	guid, sign, err := m.uploadDSP(f)
	if err != nil {
		return "", "", err
	}

	param := &value.ControlScratch{
		Type: types.StartControlScratch,
		GUID: guid,
		Sign: sign,
	}

	err = m.UB().PerformActionForKeySync(key.KeyRobomasterSystemControlScratch, param)
	if err != nil {
		return "", "", err
	}

	return guid, sign, nil
}

func (m *Module) StopDSP(guid, sign string) error {
	param := &value.ControlScratch{
		Type: types.ExitControlScratch, // We really want to quit, so use exit instead of stop.
		GUID: guid,
		Sign: sign,
	}

	err := m.UB().PerformActionForKeySync(key.KeyRobomasterSystemControlScratch, param)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) uploadDSP(f *dsp.File) (string, string, error) {
	// Save DSP to temp dir.
	tmpDir, err := os.MkdirTemp("", "lab-dsp-*")
	if err != nil {
		return "", "", err
	}

	// Make sure we clean up after ourselves.
	defer os.RemoveAll(tmpDir)

	// Save DSP to temporary directory.
	fileName, err := f.Save(tmpDir, false)
	if err != nil {
		return "", "", err
	}

	guid := f.GUID()
	sign := f.Sign()

	param2 := &value.UploadScratchByFTP{
		SkillIndex: -1,
		GUID:       guid,
		Sign:       sign,
		FilePath:   fileName,
	}

	m.uploadWg.Add(1)
	defer m.uploadWg.Wait()

	err = m.UB().SetKeyValue(key.KeyRobomasterSystemUploadScratchByFTP, param2, nil)
	if err != nil {
		return "", "", err
	}

	return guid, sign, nil
}
