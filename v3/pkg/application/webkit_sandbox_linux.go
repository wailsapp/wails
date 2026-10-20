//go:build linux && cgo && !android && !server

package application

import "github.com/wailsapp/wails/v3/internal/usernssandbox"

// warnIfWebKitSandboxBlocked logs actionable guidance when WebKitGTK's bwrap
// sandbox cannot start (typically AppArmor restricting unprivileged user
// namespaces on Ubuntu 24.04), which otherwise ends in an opaque SIGTRAP when
// the first web view is created. It only reports; the sandbox is never disabled.
// See https://github.com/wailsapp/wails/issues/6172.
func warnIfWebKitSandboxBlocked(app *App) {
	res := usernssandbox.Check(usernssandbox.Env{}, false)
	if res.Status == usernssandbox.OK {
		return
	}
	app.warning("WebKitGTK sandbox is likely to fail: %s\n%s", res.Message(), usernssandbox.Advice())
}
