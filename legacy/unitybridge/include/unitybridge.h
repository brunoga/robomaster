// This file includes relevant prototypes for functions that exist in DJI's
// Unity bridge library for the RoboMaster robots (S1 and EP only as far as
// I know at the time of writting this).
//
// These functions are the interface used by their Unity app to talk to the
// robot in lieu of using the existing SDK for the EP or nothing for the S1.
// Although the interface is small, it allows full control of the robot and
// basically doing anything that the DJI app itself can do. Because of that,
// it is obvious that the complexity is somewhere else and that is in their
// app that knows what events to send and when and what to do with the
// replies. Without that knowledge, all of this here is not very useful but
// I do want to point out the the DJI app for Windows is just a C# app
// (*wink* *wink*).
//
// Comments below were derived from the function names and the behavior
// checked empirically so they might not be as complete or as correct as
// one would like. Be warned.
//
// Also, the library where these functions are is closed source so some measure
// of reverse engineering was needed to get to this point. The most important
// thing to keep in mind is that Linux is not among the supported platforms
// (unless you consider Android as Linux, that is), although there is at least
// one workaround implemented to get the Windows library (DLL) to work on
// Linux by using Wine.

#include <stdbool.h>
#include <stdint.h>

// EventCallback is the type for functions that are passed to the
// UnitySetEventCallback call.
typedef void (*EventCallbackFunc)(uint64_t event_code, uintptr_t data,
                                  int length, uint64_t tag);

// CreateUnityBridge creates a new bridge that will be dynamically loaded from
// the library at the given path and using the given name and debuggable status.
// Note that the only name that seems valid is "Robomaster" so use that unless
// you want to try to figure out if there are other valid names. Log will be
// written to the file represented by log path unless it is the empty string in
// which case log will be sent to stdout. This should always be called before
// anything else. Multiple calls will be ignored.
void CreateUnityBridge(const char* path, const char* name, bool debuggable,
                       const char* log_path);

// DestroyUnityBridge destroys the bridge created with the CreateUnityBridge
// call above. Should always be called when cleaning up (and, of course, only
// after CreateUnityBridge is called).
void DestroyUnityBridge();

// UnityBridgeInitialize effectivelly starts the created bridge and sets it up
// so it can be used to talk to the RoboMaster robot. Returns true on success
// and false on failure.
bool UnityBridgeInitialize();

// UnityBridgeUninitialze (yep, the typo is there in the symbol name so we also
// have it here) tears down the setup done during initialization. Should also be
// called on cleanup (again, only after UnityBridgeInitialize() was called).
void UnityBridgeUninitialize();

// UnitySendEvent sends an event to be handled by the RoboMaster robot. The
// event is defined by its code and includes the given data payload. The tag
// parameter is used to track event completion (when there are multiple events
// in flight with the same code.
void UnitySendEvent(uint64_t event_code, uintptr_t data, int length,
                    uint64_t tag);

// UnitySendEventWithString is just like UnitySendEvent but the event payload is
// a string instead of arbitrary data.
void UnitySendEventWithString(uint64_t event_code, const char* data,
                              uint64_t tag);

// UnitySendEventWithNumber is just like UnitySendEvent but the event payload is
// a number instead of arbitrary data.
void UnitySendEventWithNumber(uint64_t event_code, uint64_t data, uint64_t tag);

// UnitySetEventCallback sets a callback to be called whenever an event with the
// given event code is completed and reported back by the RoboMaster robot.
void UnitySetEventCallback(uint64_t event_code,
                           EventCallbackFunc event_callback);

// UnityGetSecurityKeyByKeyChainIndex returns the security key associated with
// the given keychain index.
uintptr_t UnityGetSecurityKeyByKeyChainIndex(int index);