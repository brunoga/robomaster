package h264

/*
#cgo LDFLAGS: -lavcodec -lavutil -lswscale -lz
#cgo CFLAGS: -O3

#include "decoder.h"
#include "frame_callback.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/mattn/go-pointer"
)

type FrameCallback func(frameData []byte)

type Decoder struct {
	m             sync.Mutex
	cDecoder      *C.decoder
	frameCallback FrameCallback
	userData      unsafe.Pointer
}

func NewDecoder(frameCallback FrameCallback) (*Decoder, error) {
	if frameCallback == nil {
		return nil, fmt.Errorf("frame callback must noit be nil")
	}

	d := &Decoder{
		sync.Mutex{},
		nil,
		frameCallback,
		nil,
	}

	d.userData = pointer.Save(d)

	return d, nil
}

func (d *Decoder) Open() error {
	d.m.Lock()
	defer d.m.Unlock()

	if d.cDecoder != nil {
		return fmt.Errorf("decoder already opened")
	}

	d.cDecoder = C.decoder_new(C.frame_callback(C.c_frame_callback),
		d.userData)
	if d.cDecoder == nil {
		pointer.Unref(d.userData)
		return fmt.Errorf("unable to allocate C decoder")
	}

	return nil
}

func (d *Decoder) Close() error {
	d.m.Lock()
	defer d.m.Unlock()

	if d.cDecoder == nil {
		return fmt.Errorf("decoder already closed")
	}

	C.decoder_free(d.cDecoder)
	d.cDecoder = nil

	pointer.Unref(d.userData)
	d.userData = nil

	return nil
}

func (d *Decoder) SendData(data []byte) {
	C.decoder_send_data(d.cDecoder, (*C.char)(unsafe.Pointer(&data[0])),
		C.int(len(data)))
}

//export goFrameCallback
func goFrameCallback(frameData []byte, userData unsafe.Pointer) {
	decoder := pointer.Restore(userData).(*Decoder)
	decoder.frameCallback(frameData)
}
