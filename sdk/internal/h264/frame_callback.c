#include "_cgo_export.h"

extern void goFrameCallback(GoSlice frame_data, GoInt width, GoInt height,
	void* user_data);

void c_frame_callback(char* frame_data, int frame_data_size, int frame_width,
	int frame_height, void* user_data) {
	GoSlice frameDataSlice;

	frameDataSlice.data = frame_data;
	frameDataSlice.len = frame_data_size;
	frameDataSlice.cap = frame_data_size;

	goFrameCallback(frameDataSlice, frame_width, frame_height, user_data);
}
