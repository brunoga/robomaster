#include "_cgo_export.h"

extern void goFrameCallback(GoSlice frame_data, void* user_data);

void c_frame_callback(char* frame_data, int frame_data_size, void* user_data) {
	GoSlice frameDataSlice;

	frameDataSlice.data = frame_data;
	frameDataSlice.len = frame_data_size;
	frameDataSlice.cap = frame_data_size;

	goFrameCallback(frameDataSlice, user_data);
}
