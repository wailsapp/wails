package application

import (
	"errors"
	"testing"
)

func TestPermissionKindStrings(t *testing.T) {
	want := map[PermissionKind]string{
		PermissionKindCamera:          "camera",
		PermissionKindMicrophone:      "microphone",
		PermissionKindScreenRecording: "screen-recording",
		PermissionKindAccessibility:   "accessibility",
		PermissionKindLocation:        "location",
		PermissionKindNotifications:   "notifications",
		PermissionKindInputMonitoring: "input-monitoring",
		PermissionKindFullDiskAccess:  "full-disk-access",
	}
	all := AllPermissionKinds()
	if len(all) != len(want) {
		t.Fatalf("AllPermissionKinds returned %d kinds, want %d", len(all), len(want))
	}
	seen := map[PermissionKind]bool{}
	for _, kind := range all {
		if seen[kind] {
			t.Errorf("AllPermissionKinds lists %q twice", kind)
		}
		seen[kind] = true
		if got := kind.String(); got != want[kind] {
			t.Errorf("%q.String() = %q, want %q", kind, got, want[kind])
		}
	}
}

func TestPermissionStatusStrings(t *testing.T) {
	want := map[PermissionStatus]string{
		PermissionStatusNotDetermined: "not-determined",
		PermissionStatusDenied:        "denied",
		PermissionStatusAuthorized:    "authorized",
		PermissionStatusRestricted:    "restricted",
		PermissionStatusUnsupported:   "unsupported",
	}
	for status, text := range want {
		if got := status.String(); got != text {
			t.Errorf("%v.String() = %q, want %q", status, got, text)
		}
	}
}

func TestPermissionsManagerWithoutImpl(t *testing.T) {
	pm := &PermissionsManager{}
	if got := pm.Status(PermissionKindCamera); got != PermissionStatusUnsupported {
		t.Errorf("Status without impl = %q, want unsupported", got)
	}
	status, err := pm.Request(PermissionKindCamera)
	if !errors.Is(err, ErrPermissionsUnsupported) {
		t.Errorf("Request without impl err = %v, want ErrPermissionsUnsupported", err)
	}
	if status != PermissionStatusUnsupported {
		t.Errorf("Request without impl status = %q, want unsupported", status)
	}
	if err := pm.OpenSystemSettings(PermissionKindCamera); !errors.Is(err, ErrPermissionsUnsupported) {
		t.Errorf("OpenSystemSettings without impl err = %v, want ErrPermissionsUnsupported", err)
	}
}

type fakePermissions struct {
	statuses map[PermissionKind]PermissionStatus
	requests []PermissionKind
}

func (f *fakePermissions) status(kind PermissionKind) PermissionStatus {
	if status, ok := f.statuses[kind]; ok {
		return status
	}
	return PermissionStatusUnsupported
}

func (f *fakePermissions) request(kind PermissionKind) (PermissionStatus, error) {
	f.requests = append(f.requests, kind)
	f.statuses[kind] = PermissionStatusAuthorized
	return PermissionStatusAuthorized, nil
}

func (f *fakePermissions) openSystemSettings(kind PermissionKind) error { return nil }

func TestPermissionsManagerDelegatesToImpl(t *testing.T) {
	fake := &fakePermissions{statuses: map[PermissionKind]PermissionStatus{
		PermissionKindMicrophone: PermissionStatusNotDetermined,
	}}
	pm := &PermissionsManager{impl: fake}
	if got := pm.Status(PermissionKindMicrophone); got != PermissionStatusNotDetermined {
		t.Fatalf("Status = %q, want not-determined", got)
	}
	status, err := pm.Request(PermissionKindMicrophone)
	if err != nil || status != PermissionStatusAuthorized {
		t.Fatalf("Request = %q, %v; want authorized, nil", status, err)
	}
	if len(fake.requests) != 1 || fake.requests[0] != PermissionKindMicrophone {
		t.Fatalf("impl saw requests %v, want [microphone]", fake.requests)
	}
	if got := pm.Status(PermissionKindMicrophone); got != PermissionStatusAuthorized {
		t.Fatalf("Status after request = %q, want authorized", got)
	}
}

func TestMediaCapturePermissionDecision(t *testing.T) {
	cases := []struct {
		name       string
		perms      map[PermissionType]Permission
		audio      bool
		video      bool
		wantResult mediaCaptureDecision
	}{
		{"nil map prompts", nil, true, true, mediaCapturePrompt},
		{"default prompts", map[PermissionType]Permission{}, false, true, mediaCapturePrompt},
		{"camera allowed grants video", map[PermissionType]Permission{PermissionCamera: PermissionAllow}, false, true, mediaCaptureGrant},
		{"camera allowed still prompts for mic", map[PermissionType]Permission{PermissionCamera: PermissionAllow}, true, true, mediaCapturePrompt},
		{"both allowed grants both", map[PermissionType]Permission{PermissionCamera: PermissionAllow, PermissionMicrophone: PermissionAllow}, true, true, mediaCaptureGrant},
		{"mic denied denies both", map[PermissionType]Permission{PermissionCamera: PermissionAllow, PermissionMicrophone: PermissionDeny}, true, true, mediaCaptureDeny},
		{"mic denied does not affect video", map[PermissionType]Permission{PermissionMicrophone: PermissionDeny}, false, true, mediaCapturePrompt},
		{"nothing requested prompts", map[PermissionType]Permission{PermissionCamera: PermissionDeny}, false, false, mediaCapturePrompt},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := mediaCapturePermissionDecision(tc.perms, tc.audio, tc.video); got != tc.wantResult {
				t.Errorf("got %d, want %d", got, tc.wantResult)
			}
		})
	}
}

func TestMediaCaptureDecisionValuesMatchWebKit(t *testing.T) {
	// WKPermissionDecision: Prompt = 0, Grant = 1, Deny = 2.
	if mediaCapturePrompt != 0 || mediaCaptureGrant != 1 || mediaCaptureDeny != 2 {
		t.Fatalf("decision values drifted from WKPermissionDecision: %d %d %d", mediaCapturePrompt, mediaCaptureGrant, mediaCaptureDeny)
	}
}
