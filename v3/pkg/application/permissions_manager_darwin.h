//go:build darwin && !ios && !server

#ifndef permissions_manager_darwin_h
#define permissions_manager_darwin_h

#include <stdbool.h>

// Permission kinds, kept in step with permissionKindCode in
// permissions_manager_darwin.go.
enum {
    WailsPermissionCamera = 0,
    WailsPermissionMicrophone = 1,
    WailsPermissionScreenRecording = 2,
    WailsPermissionAccessibility = 3,
    WailsPermissionLocation = 4,
    WailsPermissionNotifications = 5,
    WailsPermissionInputMonitoring = 6,
    WailsPermissionFullDiskAccess = 7,
};

// Permission statuses, kept in step with permissionStatusFromCode in
// permissions_manager_darwin.go.
enum {
    WailsPermissionStatusNotDetermined = 0,
    WailsPermissionStatusDenied = 1,
    WailsPermissionStatusAuthorized = 2,
    WailsPermissionStatusRestricted = 3,
    WailsPermissionStatusUnsupported = 4,
};

// Returned by wailsPermissionRequest when the answer arrives later through
// permissionsRequestResult.
#define WailsPermissionRequestPending (-1)

// wailsPermissionStatus returns the current status of a kind without
// prompting. Safe to call from any thread.
int wailsPermissionStatus(int kind);

// wailsPermissionRequest starts a request for a kind. It returns the final
// status when the framework answers synchronously, or
// WailsPermissionRequestPending when the answer will be delivered through
// permissionsRequestResult(requestID, status). Must be called on the main
// thread (CoreLocation needs the main run loop).
int wailsPermissionRequest(int kind, unsigned int requestID);

// wailsPermissionOpenSettings opens the System Settings pane for a kind.
// Returns false when the URL could not be opened.
bool wailsPermissionOpenSettings(int kind);

// wailsMediaCapturePermissionDecision applies the window's Permissions map to
// a WKUIDelegate media-capture request. captureType is a WKMediaCaptureType
// (0 camera, 1 microphone, 2 both); the result is a WKPermissionDecision
// (0 prompt, 1 grant, 2 deny). Called by the WebviewWindowDelegate hook in
// webview_window_darwin.m.
int wailsMediaCapturePermissionDecision(unsigned int windowId, int captureType);

#endif /* permissions_manager_darwin_h */
