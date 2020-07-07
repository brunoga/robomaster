#include "decoder.h"

#include <stddef.h>

#include <libavutil/imgutils.h>

decoder* decoder_new(frame_callback frame_callback, void* user_data) {
    decoder* d = malloc(sizeof(decoder));
    if (!d) return NULL;

    d->codec = avcodec_find_decoder(AV_CODEC_ID_H264);
    if (!d->codec) return NULL;

    d->codec_context = avcodec_alloc_context3(d->codec);
    if (!d->codec_context) return NULL;

    if (d->codec->capabilities & AV_CODEC_CAP_TRUNCATED) {
        d->codec_context->flags |= AV_CODEC_FLAG_TRUNCATED;
    }

    int err = avcodec_open2(d->codec_context, d->codec, NULL);
    if (err < 0) {
        decoder_free(d);
        return NULL;
    }

    d->codec_parser_context = av_parser_init(AV_CODEC_ID_H264);
    if (!d->codec_parser_context) {
        decoder_free(d);
        return NULL;
    }

    d->sws_context = NULL;

    d->frame = av_frame_alloc();
    if (!d->frame) {
        decoder_free(d);
        return NULL;
    }

    d->frame_rgb = av_frame_alloc();
    if (!d->frame_rgb) {
        decoder_free(d);
        return NULL;
    }

    d->frame_callback = frame_callback;
    d->user_data = user_data;

    d->output_buffer = NULL;
    d->output_buffer_size = 0;

    return d;
}

void decoder_free(decoder* d) {
    if (!d) return;

    if (d->codec_parser_context) {
        av_parser_close(d->codec_parser_context);
        d->codec_parser_context = NULL;
    }

    if (d->sws_context) {
        sws_freeContext(d->sws_context);
        d->sws_context = NULL;
    }

    if (d->frame)
        av_frame_free(&(d->frame));

    if (d->frame_rgb)
        av_frame_free(&(d->frame_rgb));

    if (d->codec_context)
        avcodec_free_context(&(d->codec_context));

    d->frame_callback = NULL;
    d->user_data = NULL;

    if (d->output_buffer != NULL) {
        free(d->output_buffer);

        d->output_buffer = NULL;
        d->output_buffer_size = 0;
    }

    free(d);
}

int decoder_parse_data(decoder* d, const char* data, int size,
    char* parsed_data, int* parsed_data_size) {
    if (!d) return -1;
    if (!data) return -1;
    if (size < 0) return -1;

    return av_parser_parse2(d->codec_parser_context,
        d->codec_context, (uint8_t **)&parsed_data, parsed_data_size,
        (uint8_t*)data, size, 0, 0, AV_NOPTS_VALUE);
}

void decoder_send_data(decoder* d, char* data, int size) {
    if (!d) return;
    if (!data) return;
    if (size < 0) return;

    uint8_t* parsed_data = NULL;
    int parsed_data_size = 0;

    uint8_t* data_ptr = data;
    while (size > 0) {
        int len = av_parser_parse2(d->codec_parser_context,
            d->codec_context, &parsed_data, &parsed_data_size,
            (uint8_t*)data_ptr, size, 0, 0, AV_NOPTS_VALUE);

        data_ptr += len;
        size -= len;

        if (parsed_data_size > 0) {
            AVPacket packet;
            av_init_packet(&packet);

            packet.data = parsed_data;
            packet.size = parsed_data_size;

            if (avcodec_send_packet(d->codec_context, &packet) != 0)
                continue;

            if (avcodec_receive_frame(d->codec_context, d->frame) != 0)
                continue;

            const AVFrame* frame = d->frame;

            d->sws_context = sws_getCachedContext(d->sws_context,
                frame->width, frame->height, frame->format, frame->width,
                frame->height, AV_PIX_FMT_RGBA, SWS_BILINEAR, NULL,
                NULL, NULL);
            if (!d->sws_context) continue;

            if (!d->output_buffer) {
                d->output_buffer_size = av_image_get_buffer_size(
                    AV_PIX_FMT_RGB32, frame->width, frame->height, 1);
                d->output_buffer = malloc(d->output_buffer_size);
            }

            av_image_fill_arrays(d->frame_rgb->data, d->frame_rgb->linesize,
                d->output_buffer, AV_PIX_FMT_RGB32, frame->width, frame->height, 1);

            sws_scale(d->sws_context, (const uint8_t * const*)frame->data,
                frame->linesize, 0, frame->height, d->frame_rgb->data,
                d->frame_rgb->linesize);

            d->frame_rgb->width = frame->width;
            d->frame_rgb->height = frame->height;

            d->frame_callback(d->output_buffer, d->output_buffer_size,
                d->user_data);
        }
    }
}
