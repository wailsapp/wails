//go:build darwin && !ios && !server

package application

/*
#include <stdbool.h>
*/
import "C"

// captureDecision mirrors WKPermissionDecision, whose cases the delegate in
// webview_window_darwin.m hands straight back to WebKit. Written out rather
// than cast from Permission: the two enums happen to agree today, and a
// coincidence between unrelated ABIs is not something to build on.
type captureDecision int

const (
	captureDecisionPrompt captureDecision = iota
	captureDecisionGrant
	captureDecisionDeny
)

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

func mediaCaptureDecision(windowID uint, needAudio bool, needVideo bool) captureDecision {
	if !needAudio && !needVideo {
		return captureDecisionPrompt
	}

	decision := captureDecisionGrant
	if needAudio {
		decision = strictestCaptureDecision(decision, captureDecisionFor(windowID, PermissionMicrophone))
	}
	if needVideo {
		decision = strictestCaptureDecision(decision, captureDecisionFor(windowID, PermissionCamera))
	}

	return decision
}

func captureDecisionFor(windowID uint, kind PermissionType) captureDecision {
	switch resolvePermission(windowID, kind) {
	case PermissionAllow:
		return captureDecisionGrant
	case PermissionDeny:
		return captureDecisionDeny
	default:
		return captureDecisionPrompt
	}
}

// A request for the camera and the microphone together gets a single answer,
// and it can be no more permissive than either half on its own: deny beats
// prompt beats grant.
func strictestCaptureDecision(a, b captureDecision) captureDecision {
	if a == captureDecisionDeny || b == captureDecisionDeny {
		return captureDecisionDeny
	}
	if a == captureDecisionPrompt || b == captureDecisionPrompt {
		return captureDecisionPrompt
	}
	return captureDecisionGrant
}

// resolvePermission returns the configured Permission for the given type on the
// window identified by windowID, defaulting to PermissionDefault when the
// window or an entry is not found.
func resolvePermission(windowID uint, kind PermissionType) Permission {
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		return PermissionDefault
	}
	webviewWindow, ok := window.(*WebviewWindow)
	if !ok || webviewWindow.options.Permissions == nil {
		return PermissionDefault
	}
	return webviewWindow.options.Permissions[kind]
}
