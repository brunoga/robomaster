//go:build (darwin && amd64) || (android && arm) || (android && arm64) || (ios && arm64)

#include <dlfcn.h>
#include <stdio.h>

#include "unitybridge_loader.h"

UnityBridgeFunctions unity_bridge_functions;

void* unity_bridge_handle = NULL;

void* set_function(const char* name) {
  void* func = dlsym(unity_bridge_handle, name);
  if (!func) {
    char* error = dlerror();
    if (error) {
      fprintf(stderr, "Error loading symbol \"%s\": %s\n", name, error);
    } else {
      fprintf(stderr, "Error loading symbol \"%s\": unknown error\n", name);
    }
  }

  return func;
}

void CreateUnityBridge(const char* path, const char* name, bool debuggable,
                       const char* log_path) {
  if (unity_bridge_handle != NULL) {
    return;
  }

  unity_bridge_handle = dlopen(path, RTLD_LAZY);
  if (!unity_bridge_handle) {
    char* error = dlerror();
    if (error) {
      fprintf(stderr, "Error loading library \"%s\": %s\n", path, error);
    } else {
      fprintf(stderr, "Error loading library \"%s\": unknown error\n", path);
    }
    return;
  }

  unity_bridge_functions.create_unity_bridge =
      set_function("CreateUnityBridge");
  unity_bridge_functions.destroy_unity_bridge =
      set_function("DestroyUnityBridge");
  unity_bridge_functions.unity_bridge_initialize =
      set_function("UnityBridgeInitialize");
  unity_bridge_functions.unity_bridge_uninitialize =
      set_function("UnityBridgeUninitialze");
  unity_bridge_functions.unity_send_event =
      set_function("UnitySendEvent");
  unity_bridge_functions.unity_send_event_with_string =
      set_function("UnitySendEventWithString");
  unity_bridge_functions.unity_send_event_with_number =
      set_function("UnitySendEventWithNumber");
  unity_bridge_functions.unity_set_event_callback =
      set_function("UnitySetEventCallback");
  unity_bridge_functions.unity_get_security_key_by_keychain_index =
      set_function("UnityGetSecurityKeyByKeyChainIndex");

  unity_bridge_functions.create_unity_bridge(name, debuggable, log_path);
}

void DestroyUnityBridge() {
  if (unity_bridge_handle == NULL) {
    return;
  }

  unity_bridge_functions.destroy_unity_bridge();

  unity_bridge_functions.create_unity_bridge = NULL;
  unity_bridge_functions.destroy_unity_bridge = NULL;
  unity_bridge_functions.unity_bridge_initialize = NULL;
  unity_bridge_functions.unity_bridge_uninitialize = NULL;
  unity_bridge_functions.unity_send_event = NULL;
  unity_bridge_functions.unity_send_event_with_string = NULL;
  unity_bridge_functions.unity_send_event_with_number = NULL;
  unity_bridge_functions.unity_set_event_callback = NULL;
  unity_bridge_functions.unity_get_security_key_by_keychain_index = NULL;

  dlclose(unity_bridge_handle);

  unity_bridge_handle = NULL;
}

bool UnityBridgeInitialize() {
  return unity_bridge_functions.unity_bridge_initialize();
}

void UnityBridgeUninitialize() {
  unity_bridge_functions.unity_bridge_uninitialize();
}

void UnitySendEvent(uint64_t event_id, uintptr_t data, int data_size,
                    uint64_t callback_id) {
  unity_bridge_functions.unity_send_event(event_id, data, data_size,
                                          callback_id);
}

void UnitySendEventWithString(uint64_t event_id, const char* data,
                              uint64_t callback_id) {
  unity_bridge_functions.unity_send_event_with_string(event_id, data,
                                                      callback_id);
}

void UnitySendEventWithNumber(uint64_t event_id, uint64_t data,
                              uint64_t callback_id) {
  unity_bridge_functions.unity_send_event_with_number(event_id, data,
                                                      callback_id);
}

void UnitySetEventCallback(uint64_t event_id, EventCallbackFunc callback) {
  unity_bridge_functions.unity_set_event_callback(event_id, callback);
}

uintptr_t UnityGetSecurityKeyByKeyChainIndex(int index) {
  return unity_bridge_functions.unity_get_security_key_by_keychain_index(index);
}
