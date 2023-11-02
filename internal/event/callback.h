#ifndef EVENT_CALLBACK_H_
#define EVENT_CALLBACK_H_

#include <stdint.h>

typedef void (*EventCallback)(uint64_t event_code, uintptr_t data, int length,
                              uint64_t tag);

void eventCallbackC(uint64_t event_code, uintptr_t data, int length,
                    uint64_t tag);

#endif  // EVENT_CALLBACK_H_
