#include "event_callback_windows.h"

#include "_cgo_export.h"

#include <stdlib.h>
#include <stdint.h>

extern void eventCallbackGo(void* context, GoUint64 e, GoSlice info,
		GoUint64 tag);

void rgb_to_nrgba(char* nrgba, const char* rgb, int length) {
	for(int i = length; --i; nrgba += 4, rgb += 3) {
		*(uint32_t*)(void*)nrgba = *(const uint32_t*)(const void*)rgb;
		nrgba[3] = 255;
	}
       	
	nrgba[0] = rgb[0];
	nrgba[1] = rgb[1];
	nrgba[2] = rgb[2];
	nrgba[3] = 255;
}

void event_callback(void* context, va_alist alist) {
	va_start_void(alist);
	unsigned long long event_code = va_arg_ulonglong(alist);
        void* data = va_arg_ptr(alist, void*);
        int length = va_arg_int(alist);
        unsigned long long tag = va_arg_ulonglong(alist);

	// Create a Go slice with the data.
	GoSlice data_slice;
	data_slice.data = data;
	data_slice.len = length;
	data_slice.cap = length;

	eventCallbackGo(context, event_code, data_slice, tag);
}

