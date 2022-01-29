package video

import (
	"fmt"
	"git.bug-br.org.br/bga/robomasters1/app/internal"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge"
)

type DataHandlerFunc func([]byte, *sync.WaitGroup)

type Video struct {
	*internal.GenericController

	videoDataHandlers map[int]DataHandlerFunc
}

func New() (*Video, error) {
	vc := &Video{
		nil,
		make(map[int]DataHandlerFunc),
	}

	vc.GenericController = internal.NewGenericController(vc.HandleEvent)

	var err error
	err = vc.StartControllingEvent(unity.EventTypeGetNativeTexture)
	if err != nil {
		return nil, err
	}
	err = vc.StartControllingEvent(unity.EventTypeVideoTransferSpeed)
	if err != nil {
		return nil, err
	}
	err = vc.StartControllingEvent(unity.EventTypeVideoDataRecv)
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func (v *Video) StartVideo() {
	ub := bridge.Instance()

	ub.SendEvent(unity.NewEvent(unity.EventTypeStartVideo))
}

func (v *Video) StopVideo() {
	ub := bridge.Instance()

	ub.SendEvent(unity.NewEvent(unity.EventTypeStopVideo))
}

func (v *Video) AddDataHandler(dataHandlerFunc DataHandlerFunc) (int, error) {
	if dataHandlerFunc == nil {
		return -1, fmt.Errorf("dataHandlerFunc must not be nil")
	}

	var i int
	for i = 0; ; i++ {
		_, ok := v.videoDataHandlers[i]
		if !ok {
			v.videoDataHandlers[i] = dataHandlerFunc
			break
		}
	}

	return i, nil
}

func (v *Video) RemoveDataHandler(index int) error {
	_, ok := v.videoDataHandlers[index]
	if !ok {
		return fmt.Errorf("no dataHandlerFunc at given index")
	}

	delete(v.videoDataHandlers, index)

	return nil
}

func (v *Video) HandleEvent(event *unity.Event, info []byte,
	tag uint64, wg *sync.WaitGroup) {
	switch event.Type() {
	case unity.EventTypeGetNativeTexture:
	case unity.EventTypeVideoTransferSpeed:
	case unity.EventTypeVideoDataRecv:
		for _, handler := range v.videoDataHandlers {
			wg.Add(1)
			go handler(info, wg)
		}
	default:
	}

	wg.Done()
}
