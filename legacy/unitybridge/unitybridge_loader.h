//go:build (darwin && amd64) || (android && arm) || (android && arm64) || (ios && arm64) || (windows && amd64)

#include <stdbool.h>
#include <stdint.h>
#include <unitybridge.h>

typedef struct {
  void (*create_unity_bridge)(const char*, bool, const char*);
  void (*destroy_unity_bridge)();
  bool (*unity_bridge_initialize)();
  void (*unity_bridge_uninitialize)();
  void (*unity_send_event)(uint64_t, uintptr_t, int, uint64_t);
  void (*unity_send_event_with_string)(uint64_t, const char*, uint64_t);
  void (*unity_send_event_with_number)(uint64_t, uint64_t, uint64_t);
  void (*unity_set_event_callback)(uint64_t, EventCallbackFunc);
  uintptr_t (*unity_get_security_key_by_keychain_index)(int);
} UnityBridgeFunctions;
