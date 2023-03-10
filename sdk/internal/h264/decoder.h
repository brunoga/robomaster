#ifndef H264_DECODER_H
#define H264_DECODER_H

#include "libavcodec/avcodec.h"
#include "libswscale/swscale.h"

typedef void(*frame_callback)(char* frame_data, int frame_data_size,
    int frame_width, int frame_height, void* user_data);

typedef struct decoder {
    const AVCodec* codec;
    AVCodecContext* codec_context;
    AVCodecParserContext* codec_parser_context;

    AVFrame* frame;
    AVFrame* frame_rgb;

    struct SwsContext* sws_context;

    frame_callback frame_callback;
    void* user_data;

    char* output_buffer;
    int output_buffer_size;
} decoder;

decoder* decoder_new(frame_callback frame_callback, void* user_data);
void decoder_free(decoder* d);
void decoder_send_data(decoder* d, char* data, int size);

#endif  // H264_DECODER_H