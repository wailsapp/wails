//go:build linux && cgo && !android && !server

package application

// resolvePermission returns the configured Permission for the given type on the
// window identified by windowID, defaulting to PermissionDefault when the
// window or an entry is not found.
func resolvePermission(windowID uint, kind PermissionType) Permission {
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		return PermissionDefault
	}
	lw := getLinuxWebviewWindow(window)
	if lw == nil || lw.parent.options.Permissions == nil {
		return PermissionDefault
	}
	return lw.parent.options.Permissions[kind]
}

// allowMediaCapture decides whether a getUserMedia request for the requested
// device types is permitted, applying the window's Permissions. WebKitGTK has
// no native permission prompt, so PermissionDefault allows media capture
// (restoring getUserMedia for app content, see #5552); an explicit
// PermissionDeny turns it off.
func allowMediaCapture(windowID uint, needAudio, needVideo bool) bool {
	allows := func(kind PermissionType) bool {
		switch resolvePermission(windowID, kind) {
		case PermissionDeny:
			return false
		default: // PermissionDefault and PermissionAllow
			return true
		}
	}
	if needAudio && !allows(PermissionMicrophone) {
		return false
	}
	if needVideo && !allows(PermissionCamera) {
		return false
	}
	return true
}
