package mobile

import (
	"sync"
)

type WaitGroup struct {
	wg *sync.WaitGroup
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		wg: &sync.WaitGroup{},
	}
}

func (wg *WaitGroup) Add(delta int) {
	wg.wg.Add(delta)
}

func (wg *WaitGroup) Done() {
	wg.wg.Done()
}

func (wg *WaitGroup) Wait() {
	wg.wg.Wait()
}


