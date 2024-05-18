#include "callback.h"

#include <stdint.h>

#include "_cgo_export.h"

void eventCallbackC(uint64_t event_code, uintptr_t data, int length,
                    uint64_t tag) {
  GoSlice data_slice;
  data_slice.data = (void *)data;
  data_slice.len = length;
  data_slice.cap = length;

  eventCallbackGo(event_code, data_slice, tag);
}