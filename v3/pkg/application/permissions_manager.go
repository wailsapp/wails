package application

import (
	"errors"
)

// PermissionKind identifies a system-level privacy permission (TCC on macOS).
// It is distinct from PermissionType, which governs what web content inside a
// window may request; see WebviewWindowOptions.Permissions for that.
type PermissionKind string

const (
	// PermissionKindCamera is access to video capture devices.
	PermissionKindCamera PermissionKind = "camera"
	// PermissionKindMicrophone is access to audio capture devices.
	PermissionKindMicrophone PermissionKind = "microphone"
	// PermissionKindScreenRecording is permission to capture screen contents.
	PermissionKindScreenRecording PermissionKind = "screen-recording"
	// PermissionKindAccessibility is the accessibility (assistive access)
	// permission that allows controlling other applications and observing
	// keyboard input system-wide.
	PermissionKindAccessibility PermissionKind = "accessibility"
	// PermissionKindLocation is access to location services.
	PermissionKindLocation PermissionKind = "location"
	// PermissionKindNotifications is permission to post user notifications.
	PermissionKindNotifications PermissionKind = "notifications"
	// PermissionKindInputMonitoring is permission to observe HID input events
	// from other applications.
	PermissionKindInputMonitoring PermissionKind = "input-monitoring"
	// PermissionKindFullDiskAccess is permission to read files protected by
	// the system (Mail, Messages, Safari data, Time Machine backups, ...).
	PermissionKindFullDiskAccess PermissionKind = "full-disk-access"
)

// String returns the wire name of the permission kind.
func (k PermissionKind) String() string { return string(k) }

// AllPermissionKinds returns every PermissionKind in declaration order.
func AllPermissionKinds() []PermissionKind {
	return []PermissionKind{
		PermissionKindCamera,
		PermissionKindMicrophone,
		PermissionKindScreenRecording,
		PermissionKindAccessibility,
		PermissionKindLocation,
		PermissionKindNotifications,
		PermissionKindInputMonitoring,
		PermissionKindFullDiskAccess,
	}
}

// PermissionStatus is the current authorisation state of a PermissionKind.
type PermissionStatus string

const (
	// PermissionStatusNotDetermined means the user has not been asked yet.
	PermissionStatusNotDetermined PermissionStatus = "not-determined"
	// PermissionStatusDenied means the user (or a probe) refused access.
	PermissionStatusDenied PermissionStatus = "denied"
	// PermissionStatusAuthorized means access is granted.
	PermissionStatusAuthorized PermissionStatus = "authorized"
	// PermissionStatusRestricted means access is blocked by policy (MDM,
	// parental controls) or the service is disabled system-wide; the user
	// cannot change it from within the app.
	PermissionStatusRestricted PermissionStatus = "restricted"
	// PermissionStatusUnsupported means the platform, OS version or process
	// (for example an unbundled binary asking about notifications) cannot
	// answer the question.
	PermissionStatusUnsupported PermissionStatus = "unsupported"
)

// String returns the wire name of the status.
func (s PermissionStatus) String() string { return string(s) }

// ErrPermissionsUnsupported is returned by Request and OpenSystemSettings on
// platforms that have no permission API for the requested kind.
var ErrPermissionsUnsupported = errors.New("permissions: not supported on this platform")

// ErrPermissionNotRequestable is returned by Request for kinds that cannot be
// requested programmatically (full disk access): direct the user to
// OpenSystemSettings instead.
var ErrPermissionNotRequestable = errors.New("permissions: this permission cannot be requested programmatically; use OpenSystemSettings")

// ErrPermissionRequestOnMainThread is returned when Request is called from the
// main thread. Requests wait for a system callback that is delivered on the
// main thread, so the caller must be a different goroutine.
var ErrPermissionRequestOnMainThread = errors.New("permissions: Request must not be called from the main thread")

// ErrPermissionRequestTimeout is returned when the system never answered a
// request. On macOS this usually means the Info.plist usage-description key
// for the kind is missing, in which case the OS silently drops the prompt.
var ErrPermissionRequestTimeout = errors.New("permissions: request timed out waiting for the system")

// ErrPermissionUnknownKind is returned for a PermissionKind this package does
// not know about.
var ErrPermissionUnknownKind = errors.New("permissions: unknown permission kind")

// permissionsImpl is the platform side of PermissionsManager.
type permissionsImpl interface {
	status(kind PermissionKind) PermissionStatus
	request(kind PermissionKind) (PermissionStatus, error)
	openSystemSettings(kind PermissionKind) error
}

// PermissionsManager requests and reports system permissions such as camera, microphone, screen recording and location.
//
// On macOS the answers come from TCC through the framework that owns each
// permission (AVFoundation, CoreGraphics, ApplicationServices, CoreLocation,
// UserNotifications, IOKit). Prompts only appear for kinds that are
// PermissionStatusNotDetermined and require the matching usage-description
// key in Info.plist (NSCameraUsageDescription, NSMicrophoneUsageDescription,
// NSLocationUsageDescription); without it macOS logs a message and never
// shows the dialog.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type PermissionsManager struct {
	app  *App
	impl permissionsImpl
}

func newPermissionsManager(app *App) *PermissionsManager {
	return &PermissionsManager{app: app, impl: newPermissionsImpl(app)}
}

// Status returns the current authorisation state for kind without prompting.
//
// macOS: answered natively for every kind. Full disk access has no query
// API, so it is probed by opening a TCC-protected file: a readable file is
// PermissionStatusAuthorized, a permission error is PermissionStatusDenied
// (macOS does not distinguish "never asked" from "refused" here).
// Screen recording and accessibility likewise only report Authorized or
// Denied. Notifications report PermissionStatusUnsupported in an unbundled
// binary because UNUserNotificationCenter needs a bundle identifier.
//
// Other platforms: PermissionStatusUnsupported.
func (pm *PermissionsManager) Status(kind PermissionKind) PermissionStatus {
	if pm.impl == nil {
		return PermissionStatusUnsupported
	}
	return pm.impl.status(kind)
}

// Request asks the user for kind and blocks until the system answers,
// returning the resulting status. Kinds that are already determined return
// immediately without a prompt. Request must be called from a goroutine
// other than the main thread; see ErrPermissionRequestOnMainThread.
//
// macOS: camera and microphone use AVCaptureDevice; screen recording uses
// CGRequestScreenCaptureAccess (the answer only takes effect after the app is
// relaunched); accessibility shows the "allow in System Settings" dialog and
// returns the current value; input monitoring uses IOHIDRequestAccess;
// location uses CLLocationManager and keeps the manager alive for the
// process; notifications use UNUserNotificationCenter and require a bundled
// app. Full disk access cannot be requested (ErrPermissionNotRequestable).
//
// Other platforms: ErrPermissionsUnsupported.
func (pm *PermissionsManager) Request(kind PermissionKind) (PermissionStatus, error) {
	if pm.impl == nil {
		return PermissionStatusUnsupported, ErrPermissionsUnsupported
	}
	return pm.impl.request(kind)
}

// OpenSystemSettings opens the system preference pane where the user can
// change kind for this application.
//
// macOS: opens the matching Privacy & Security pane through an
// x-apple.systempreferences: URL.
//
// Other platforms: ErrPermissionsUnsupported.
func (pm *PermissionsManager) OpenSystemSettings(kind PermissionKind) error {
	if pm.impl == nil {
		return ErrPermissionsUnsupported
	}
	return pm.impl.openSystemSettings(kind)
}

// mediaCaptureDecision is the answer given to a webview media-capture
// request. The values match WebKit's WKPermissionDecision so the macOS
// delegate can pass them straight through.
type mediaCaptureDecision int

const (
	mediaCapturePrompt mediaCaptureDecision = 0
	mediaCaptureGrant  mediaCaptureDecision = 1
	mediaCaptureDeny   mediaCaptureDecision = 2
)

// mediaCapturePermissionDecision applies a window's Permissions map to a
// getUserMedia request. Any requested kind set to PermissionDeny denies the
// whole request; the request is granted silently only when every requested
// kind is PermissionAllow; otherwise the platform prompt is shown
// (PermissionDefault). PermissionAllow only skips the webview's own prompt:
// the OS-level TCC prompt for camera and microphone still appears the first
// time, see PermissionsManager.Request.
func mediaCapturePermissionDecision(perms map[PermissionType]Permission, needAudio, needVideo bool) mediaCaptureDecision {
	requested := make([]PermissionType, 0, 2)
	if needAudio {
		requested = append(requested, PermissionMicrophone)
	}
	if needVideo {
		requested = append(requested, PermissionCamera)
	}
	if len(requested) == 0 {
		return mediaCapturePrompt
	}
	allAllowed := true
	for _, kind := range requested {
		switch perms[kind] {
		case PermissionDeny:
			return mediaCaptureDeny
		case PermissionAllow:
		default:
			allAllowed = false
		}
	}
	if allAllowed {
		return mediaCaptureGrant
	}
	return mediaCapturePrompt
}

// windowPermissions returns the Permissions map configured on the window
// with the given ID, or nil when the window is unknown.
func windowPermissions(windowID uint) map[PermissionType]Permission {
	if globalApplication == nil {
		return nil
	}
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		return nil
	}
	webviewWindow, ok := window.(*WebviewWindow)
	if !ok {
		return nil
	}
	return webviewWindow.options.Permissions
}
