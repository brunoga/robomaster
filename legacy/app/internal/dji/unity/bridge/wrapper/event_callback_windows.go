package wrapper

import "C"
import (
	"unsafe"

	"github.com/mattn/go-pointer"
)

//export eventCallbackGo
func eventCallbackGo(context unsafe.Pointer, eventCode uint64, data []byte,
	tag uint64) {
	cb := pointer.Restore(context).(EventCallback)
	cb(eventCode, data, tag)
}
