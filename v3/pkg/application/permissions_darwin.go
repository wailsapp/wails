//go:build darwin && !ios && !server && !purego

package application

/*
#include <stdbool.h>
*/
import "C"

// resolveMediaCapturePermission answers a getUserMedia request for the given
// window, applying its Permissions.
//
// With no delegate method, WebKit takes the default action, which on Cocoa is
// its own prompt — so capture works, but the window's Permissions are never
// consulted and macOS behaves as if the option were permanently
// PermissionDefault. This is what gives PermissionAllow and PermissionDeny
// their meaning there; PermissionDefault keeps that same prompt.
//
//export resolveMediaCapturePermission
func resolveMediaCapturePermission(windowID C.uint, needAudio C.bool, needVideo C.bool) C.int {
	return C.int(mediaCaptureDecision(uint(windowID), bool(needAudio), bool(needVideo)))
}
