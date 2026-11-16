//go:build darwin && purego && !ios && !server

package application

import (
	"testing"

	"github.com/ebitengine/purego/objc"
)

const mediaCapturePermissionSelector = "webView:requestMediaCapturePermissionForOrigin:initiatedByFrame:type:decisionHandler:"

// The values are WKMediaCaptureType as WebKit defines them. A wrong number here
// would ask for the wrong device, so they are pinned explicitly.
func TestWKMediaCaptureTypeValues(t *testing.T) {
	if wkMediaCaptureTypeCamera != 0 {
		t.Error("wkMediaCaptureTypeCamera should be 0 (WKMediaCaptureTypeCamera)")
	}
	if wkMediaCaptureTypeMicrophone != 1 {
		t.Error("wkMediaCaptureTypeMicrophone should be 1 (WKMediaCaptureTypeMicrophone)")
	}
	if wkMediaCaptureTypeCameraAndMicrophone != 2 {
		t.Error("wkMediaCaptureTypeCameraAndMicrophone should be 2 (WKMediaCaptureTypeCameraAndMicrophone)")
	}
}

func TestMediaCaptureNeeds(t *testing.T) {
	tests := []struct {
		name      string
		mediaType int
		wantAudio bool
		wantVideo bool
	}{
		{"camera asks for video only", int(wkMediaCaptureTypeCamera), false, true},
		{"microphone asks for audio only", int(wkMediaCaptureTypeMicrophone), true, false},
		{"camera and microphone ask for both", int(wkMediaCaptureTypeCameraAndMicrophone), true, true},
		// Unknown values must ask for nothing, so the decision stays WebKit's prompt.
		{"unknown type asks for neither", 3, false, false},
		{"negative type asks for neither", -1, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAudio, gotVideo := mediaCaptureNeeds(tt.mediaType)
			if gotAudio != tt.wantAudio || gotVideo != tt.wantVideo {
				t.Errorf("mediaCaptureNeeds(%d) = (audio=%t, video=%t), want (audio=%t, video=%t)",
					tt.mediaType, gotAudio, gotVideo, tt.wantAudio, tt.wantVideo)
			}
		})
	}
}

// The full path from a window's Permissions to the WKPermissionDecision integer
// that the handler receives: 0 prompt, 1 grant, 2 deny.
func TestMediaCaptureDecisionForType(t *testing.T) {
	const windowID = 1

	tests := []struct {
		name        string
		permissions map[PermissionType]Permission
		mediaType   int
		want        int
	}{
		{"camera allowed", map[PermissionType]Permission{PermissionCamera: PermissionAllow}, int(wkMediaCaptureTypeCamera), 1},
		{"camera denied", map[PermissionType]Permission{PermissionCamera: PermissionDeny}, int(wkMediaCaptureTypeCamera), 2},
		{"camera unset prompts", nil, int(wkMediaCaptureTypeCamera), 0},
		{"microphone allowed", map[PermissionType]Permission{PermissionMicrophone: PermissionAllow}, int(wkMediaCaptureTypeMicrophone), 1},
		{"microphone denied", map[PermissionType]Permission{PermissionMicrophone: PermissionDeny}, int(wkMediaCaptureTypeMicrophone), 2},
		{"microphone unset prompts", nil, int(wkMediaCaptureTypeMicrophone), 0},
		// A denied camera is not applied to a microphone-only request.
		{
			"denied camera ignored for microphone request",
			map[PermissionType]Permission{PermissionMicrophone: PermissionAllow, PermissionCamera: PermissionDeny},
			int(wkMediaCaptureTypeMicrophone), 1,
		},
		{
			"both allowed for camera and microphone request",
			map[PermissionType]Permission{PermissionMicrophone: PermissionAllow, PermissionCamera: PermissionAllow},
			int(wkMediaCaptureTypeCameraAndMicrophone), 1,
		},
		{
			"denied camera denies camera and microphone request",
			map[PermissionType]Permission{PermissionMicrophone: PermissionAllow, PermissionCamera: PermissionDeny},
			int(wkMediaCaptureTypeCameraAndMicrophone), 2,
		},
		{
			"unset camera prompts a camera and microphone request",
			map[PermissionType]Permission{PermissionMicrophone: PermissionAllow},
			int(wkMediaCaptureTypeCameraAndMicrophone), 0,
		},
		{
			"unknown type gets the prompt even when both are allowed",
			map[PermissionType]Permission{PermissionMicrophone: PermissionAllow, PermissionCamera: PermissionAllow},
			3, 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withWindowPermissions(t, windowID, tt.permissions)

			if got := int(mediaCaptureDecisionForType(windowID, tt.mediaType)); got != tt.want {
				t.Errorf("mediaCaptureDecisionForType(%d) = %d, want %d", tt.mediaType, got, tt.want)
			}
		})
	}
}

// AppKit: the window delegate class, which is also the window's WKUIDelegate,
// must answer the selector. Without it WebKit shows its own prompt and the
// Permissions are ignored.
func TestMediaCapturePermissionSelectorRegistered(t *testing.T) {
	cls := registerWindowDelegateClass()
	if !get[bool](cls, "instancesRespondToSelector:", sel_(mediaCapturePermissionSelector)) {
		t.Fatalf("%s does not implement %s", "WailsWebviewWindowDelegate", mediaCapturePermissionSelector)
	}
}

// AppKit: call the registered IMP through the Objective-C runtime with a real
// block as the decision handler, and read back the integer it receives. This
// covers the handler invocation and the delegate-to-window lookup together.
func TestMediaCapturePermissionDecisionHandler(t *testing.T) {
	const windowID = 1
	cls := registerWindowDelegateClass()

	tests := []struct {
		name        string
		permissions map[PermissionType]Permission
		mediaType   int
		registered  bool // false: the delegate has no entry in delegateToWindowID
		want        int
	}{
		{"camera allowed", map[PermissionType]Permission{PermissionCamera: PermissionAllow}, int(wkMediaCaptureTypeCamera), true, 1},
		{"camera denied", map[PermissionType]Permission{PermissionCamera: PermissionDeny}, int(wkMediaCaptureTypeCamera), true, 2},
		{"camera unset prompts", nil, int(wkMediaCaptureTypeCamera), true, 0},
		{"microphone denied", map[PermissionType]Permission{PermissionMicrophone: PermissionDeny}, int(wkMediaCaptureTypeMicrophone), true, 2},
		{
			"camera and microphone, microphone allowed, camera unset",
			map[PermissionType]Permission{PermissionMicrophone: PermissionAllow},
			int(wkMediaCaptureTypeCameraAndMicrophone), true, 0,
		},
		// A delegate that is not mapped to a window keeps WebKit's prompt even
		// though the Permissions would have granted the request.
		{
			"unmapped delegate prompts",
			map[PermissionType]Permission{PermissionMicrophone: PermissionAllow, PermissionCamera: PermissionAllow},
			int(wkMediaCaptureTypeCameraAndMicrophone), false, 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withWindowPermissions(t, windowID, tt.permissions)

			del := cls.send("new")
			t.Cleanup(func() { del.send("release") })
			if tt.registered {
				delegateToWindowID.Store(del.ptr(), uint(windowID))
				t.Cleanup(func() { delegateToWindowID.Delete(del.ptr()) })
			}

			got := -1
			block := objc.NewBlock(func(b objc.Block, decision int) {
				got = decision
			})
			defer block.Release()

			del.send(mediaCapturePermissionSelector, objc.ID(0), objc.ID(0), objc.ID(0), tt.mediaType, block)

			if got != tt.want {
				t.Errorf("decision handler received %d, want %d", got, tt.want)
			}
		})
	}
}
