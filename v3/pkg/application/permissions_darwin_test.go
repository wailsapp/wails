//go:build darwin && !ios && !server

package application

import "testing"

// webview_window_darwin.m casts what resolveMediaCapturePermission returns
// straight to WKPermissionDecision, so captureDecision has to keep agreeing
// with it case for case. That agreement is why the mapping in
// captureDecisionFor is written out rather than cast from Permission, and it
// is what these three numbers pin.
func TestCaptureDecision_Constants(t *testing.T) {
	if captureDecisionPrompt != 0 {
		t.Error("captureDecisionPrompt should be 0 (WKPermissionDecisionPrompt)")
	}
	if captureDecisionGrant != 1 {
		t.Error("captureDecisionGrant should be 1 (WKPermissionDecisionGrant)")
	}
	if captureDecisionDeny != 2 {
		t.Error("captureDecisionDeny should be 2 (WKPermissionDecisionDeny)")
	}
}

func TestStrictestCaptureDecision(t *testing.T) {
	tests := []struct {
		name string
		a    captureDecision
		b    captureDecision
		want captureDecision
	}{
		{"grant with grant", captureDecisionGrant, captureDecisionGrant, captureDecisionGrant},
		{"grant with prompt", captureDecisionGrant, captureDecisionPrompt, captureDecisionPrompt},
		{"grant with deny", captureDecisionGrant, captureDecisionDeny, captureDecisionDeny},
		{"prompt with prompt", captureDecisionPrompt, captureDecisionPrompt, captureDecisionPrompt},
		{"prompt with deny", captureDecisionPrompt, captureDecisionDeny, captureDecisionDeny},
		{"deny with deny", captureDecisionDeny, captureDecisionDeny, captureDecisionDeny},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strictestCaptureDecision(tt.a, tt.b); got != tt.want {
				t.Errorf("strictestCaptureDecision(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
			// The two halves of a camera-and-microphone request arrive in
			// whichever order the caller asks about them.
			if got := strictestCaptureDecision(tt.b, tt.a); got != tt.want {
				t.Errorf("strictestCaptureDecision(%d, %d) = %d, want %d", tt.b, tt.a, got, tt.want)
			}
		})
	}
}

func TestCaptureDecisionFor(t *testing.T) {
	const windowID = 1

	tests := []struct {
		name       string
		permission Permission
		want       captureDecision
	}{
		{"allow grants", PermissionAllow, captureDecisionGrant},
		{"deny denies", PermissionDeny, captureDecisionDeny},
		{"default prompts", PermissionDefault, captureDecisionPrompt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withWindowPermissions(t, windowID, map[PermissionType]Permission{
				PermissionMicrophone: tt.permission,
			})

			if got := captureDecisionFor(windowID, PermissionMicrophone); got != tt.want {
				t.Errorf("captureDecisionFor(PermissionMicrophone) = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMediaCaptureDecision(t *testing.T) {
	const windowID = 1

	allowMicrophone := map[PermissionType]Permission{PermissionMicrophone: PermissionAllow}

	tests := []struct {
		name        string
		permissions map[PermissionType]Permission
		needAudio   bool
		needVideo   bool
		want        captureDecision
	}{
		{
			name:      "an unconfigured window keeps WebKit's prompt",
			needAudio: true,
			needVideo: true,
			want:      captureDecisionPrompt,
		},
		{
			name:        "microphone allowed, microphone asked",
			permissions: allowMicrophone,
			needAudio:   true,
			want:        captureDecisionGrant,
		},
		{
			name:        "microphone denied, microphone asked",
			permissions: map[PermissionType]Permission{PermissionMicrophone: PermissionDeny},
			needAudio:   true,
			want:        captureDecisionDeny,
		},
		{
			name:        "camera allowed, camera asked",
			permissions: map[PermissionType]Permission{PermissionCamera: PermissionAllow},
			needVideo:   true,
			want:        captureDecisionGrant,
		},
		{
			name: "both allowed, both asked",
			permissions: map[PermissionType]Permission{
				PermissionMicrophone: PermissionAllow,
				PermissionCamera:     PermissionAllow,
			},
			needAudio: true,
			needVideo: true,
			want:      captureDecisionGrant,
		},
		{
			name:        "an allowed microphone does not carry an unset camera",
			permissions: allowMicrophone,
			needAudio:   true,
			needVideo:   true,
			want:        captureDecisionPrompt,
		},
		{
			name: "a denied camera outweighs an allowed microphone",
			permissions: map[PermissionType]Permission{
				PermissionMicrophone: PermissionAllow,
				PermissionCamera:     PermissionDeny,
			},
			needAudio: true,
			needVideo: true,
			want:      captureDecisionDeny,
		},
		{
			name: "a denial the request does not ask about is not applied",
			permissions: map[PermissionType]Permission{
				PermissionMicrophone: PermissionAllow,
				PermissionCamera:     PermissionDeny,
			},
			needAudio: true,
			want:      captureDecisionGrant,
		},
		{
			name:        "a request for neither device is nothing to grant",
			permissions: allowMicrophone,
			want:        captureDecisionPrompt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withWindowPermissions(t, windowID, tt.permissions)

			got := mediaCaptureDecision(windowID, tt.needAudio, tt.needVideo)
			if got != tt.want {
				t.Errorf("mediaCaptureDecision(audio=%t, video=%t) = %d, want %d", tt.needAudio, tt.needVideo, got, tt.want)
			}
		})
	}
}

// A request can outlive the window it came from, and windowID is then a window
// the manager no longer knows about. WebKit's own prompt is the answer that
// changes nothing.
func TestMediaCaptureDecisionForUnknownWindow(t *testing.T) {
	withWindowPermissions(t, 1, map[PermissionType]Permission{
		PermissionMicrophone: PermissionAllow,
		PermissionCamera:     PermissionAllow,
	})

	if got := resolvePermission(2, PermissionMicrophone); got != PermissionDefault {
		t.Errorf("resolvePermission for an unknown window = %d, want PermissionDefault", got)
	}
	if got := mediaCaptureDecision(2, true, true); got != captureDecisionPrompt {
		t.Errorf("mediaCaptureDecision for an unknown window = %d, want %d", got, captureDecisionPrompt)
	}
}

// withWindowPermissions installs an application holding one window with the
// given Permissions, and restores the previous one when the test ends.
func withWindowPermissions(t *testing.T, windowID uint, permissions map[PermissionType]Permission) {
	t.Helper()

	previous := globalApplication

	app := &App{
		windows: map[uint]Window{
			windowID: &WebviewWindow{
				id:      windowID,
				options: WebviewWindowOptions{Permissions: permissions},
			},
		},
	}
	app.Window = newWindowManager(app)

	globalApplication = app
	t.Cleanup(func() { globalApplication = previous })
}
