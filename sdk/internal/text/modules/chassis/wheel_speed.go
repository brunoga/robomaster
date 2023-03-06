package chassis

import (
	"fmt"
	"sync"
)

type WheelSpeed struct {
	m  sync.RWMutex
	w1 float64
	w2 float64
	w3 float64
	w4 float64
}

func NewWheelSpeed(w1, w2, w3, w4 float64) *WheelSpeed {
	return &WheelSpeed{
		sync.RWMutex{},
		w1,
		w2,
		w3,
		w4,
	}
}

func NewWheelSpeedFromData(data string) (*WheelSpeed, error) {
	w := &WheelSpeed{
		sync.RWMutex{},
		0.0,
		0.0,
		0.0,
		0.0,
	}

	err := w.UpdateFromData(data)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WheelSpeed) Update(w1, w2, w3, w4 float64) {
	w.m.Lock()

	w.w1 = w1
	w.w2 = w2
	w.w3 = w3
	w.w4 = w4

	w.m.Unlock()
}

func (w *WheelSpeed) UpdateFromData(data string) error {
	var w1, w2, w3, w4 float64

	n, err := fmt.Sscanf(data, "%f %f %f %f", &w1, &w2, &w3, &w4)
	if err != nil {
		fmt.Errorf("error parsing data: %w", err)
	}
	if n != 7 {
		fmt.Errorf(
			"unexpected number of entries in data: %w", err)
	}

	w.m.Lock()

	w.w1 = w1
	w.w2 = w2
	w.w3 = w3
	w.w4 = w4

	w.m.Unlock()

	return nil
}

func (w *WheelSpeed) W1() float64 {
	w.m.RLock()
	defer w.m.RUnlock()

	return w.w1
}

func (w *WheelSpeed) W2() float64 {
	w.m.RLock()
	defer w.m.RUnlock()

	return w.w2
}

func (w *WheelSpeed) W3() float64 {
	w.m.RLock()
	defer w.m.RUnlock()

	return w.w3
}

func (w *WheelSpeed) W4() float64 {
	w.m.RLock()
	defer w.m.RUnlock()

	return w.w4
}
