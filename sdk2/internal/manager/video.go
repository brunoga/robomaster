package manager

import (
	"fmt"
	"sync"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

type Video struct {
	ub unitybridge.UnityBridge
	l  *logger.Logger

	gntToken token.Token
	vtsToken token.Token
	vdrToken token.Token
}

func NewVideo(ub unitybridge.UnityBridge, l *logger.Logger) (*Video, error) {
	return &Video{
		ub: ub,
		l:  l,
	}, nil
}

func (v *Video) Start() error {
	var err error
	v.gntToken, err = v.ub.AddEventTypeListener(event.TypeGetNativeTexture,
		v.onGetNativeTexture)
	if err != nil {
		return err
	}

	v.vtsToken, err = v.ub.AddEventTypeListener(event.TypeVideoTransferSpeed,
		v.onVideoTransferSpeed)
	if err != nil {
		return err
	}

	v.vdrToken, err = v.ub.AddEventTypeListener(event.TypeVideoDataRecv,
		v.onVideoDataRecv)
	if err != nil {
		return err
	}

	err = v.ub.SendEvent(event.NewFromType(event.TypeStartVideo))
	if err != nil {
		return err
	}

	// Ask for video texture information.
	err = v.ub.SendEvent(event.NewFromType(event.TypeGetNativeTexture))
	if err != nil {
		return err
	}

	return nil
}

// SetVideoResolution sets the video resolution.
//
// TODO(bga): Other then  actually limiting the available resolutions, it looks
// like changing resolutions is not working. Need to investigate further as
// there might be some setup that is needed and is not being done.
func (v *Video) SetVideoResolution(resolutionID uint64) error {
	var err error

	var wg sync.WaitGroup
	wg.Add(1)

	v.ub.SetKeyValue(key.KeyCameraVideoFormat, resolutionID, func(r *result.Result) {
		if r.ErrorCode() != 0 {
			err = fmt.Errorf("error setting video resolution: %s", r.ErrorDesc())
		}
		wg.Done()
	})

	wg.Wait()

	return err
}

func (v *Video) StartRecordingToSDCard() error {
	var err error

	var wg sync.WaitGroup
	wg.Add(1)

	v.ub.GetKeyValue(key.KeyCameraMode, func(r *result.Result) {
		if r.ErrorCode() != 0 {
			err = fmt.Errorf("error getting camera mode: %s", r.ErrorDesc())
		} else {
			if uint64(r.Value().(float64)) != 1 {
				wg.Add(1)
				v.ub.SetKeyValue(key.KeyCameraMode, uint64(1), func(r *result.Result) {
					if r.ErrorCode() != 0 {
						err = fmt.Errorf("error setting camera mode to video: %s", r.ErrorDesc())
					}
					wg.Done()
				})
			}
		}

		wg.Done()
	})

	wg.Wait()

	wg.Add(1)
	err = v.ub.PerformActionForKey(key.KeyCameraStartRecordVideo, nil, func(r *result.Result) {
		if r.ErrorCode() != 0 {
			err = fmt.Errorf("error starting video recording: %s", r.ErrorDesc())
		}
		wg.Done()
	})

	wg.Wait()

	return err
}

func (v *Video) StopRecordingToSDCard() error {
	var err error

	var wg sync.WaitGroup
	wg.Add(1)

	err = v.ub.PerformActionForKey(key.KeyCameraStopRecordVideo, nil, func(r *result.Result) {
		if r.ErrorCode() != 0 {
			err = fmt.Errorf("error stopping video recording: %s", r.ErrorDesc())
		}

		wg.Done()
	})

	return err
}

func (v *Video) Stop() error {
	var err error

	err = v.ub.SendEvent(event.NewFromType(event.TypeStopVideo))
	if err != nil {
		return err
	}

	err = v.ub.RemoveEventTypeListener(event.TypeGetNativeTexture, v.gntToken)
	if err != nil {
		return err
	}

	err = v.ub.RemoveEventTypeListener(event.TypeVideoTransferSpeed, v.vtsToken)
	if err != nil {
		return err
	}

	err = v.ub.RemoveEventTypeListener(event.TypeVideoDataRecv, v.vdrToken)
	if err != nil {
		return err
	}

	return nil
}

func (v *Video) onGetNativeTexture(data []byte, dataType event.DataType) {
	v.l.Debug("onGetNativeTexture", "data", data, "dataType", dataType)
}

func (v *Video) onVideoTransferSpeed(data []byte, dataType event.DataType) {
	v.l.Debug("onVideoTransferSpeed", "data", data, "dataType", dataType)
}

func (v *Video) onVideoDataRecv(data []byte, dataType event.DataType) {
	v.l.Debug("onVideoDataRecv", "len(data)", len(data), "dataType", dataType)
}
