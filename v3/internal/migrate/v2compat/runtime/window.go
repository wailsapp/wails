package runtime

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// themeWarnOnce guards the one-time warning emitted by the theme functions.
var themeWarnOnce sync.Once

// themeWarn logs a warning the first time any of the theme functions is called.
func themeWarn() {
	themeWarnOnce.Do(func() {
		logger().Warn("v2compat: WindowSetSystemDefaultTheme/WindowSetLightTheme/WindowSetDarkTheme are no-ops in v3; configure the theme at window creation via application.WebviewWindowOptions")
	})
}

// WindowSetTitle mirrors the v2 runtime.WindowSetTitle function.
// v3 equivalent: window.SetTitle.
func WindowSetTitle(_ context.Context, title string) {
	if w := currentWindow(); w != nil {
		w.SetTitle(title)
	}
}

// WindowFullscreen mirrors the v2 runtime.WindowFullscreen function.
// v3 equivalent: window.Fullscreen.
func WindowFullscreen(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.Fullscreen()
	}
}

// WindowUnfullscreen mirrors the v2 runtime.WindowUnfullscreen function.
// v3 equivalent: window.UnFullscreen.
func WindowUnfullscreen(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.UnFullscreen()
	}
}

// WindowCenter mirrors the v2 runtime.WindowCenter function.
// v3 equivalent: window.Center.
func WindowCenter(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.Center()
	}
}

// WindowReload mirrors the v2 runtime.WindowReload function.
// v3 equivalent: window.Reload.
func WindowReload(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.Reload()
	}
}

// WindowReloadApp mirrors the v2 runtime.WindowReloadApp function.
// v3 equivalent: window.ForceReload.
func WindowReloadApp(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.ForceReload()
	}
}

// WindowSetSystemDefaultTheme mirrors the v2 runtime.WindowSetSystemDefaultTheme function.
// v3 has no runtime theme setter: the theme is configured at window creation
// via application.WebviewWindowOptions. This is a no-op that logs a warning once.
func WindowSetSystemDefaultTheme(_ context.Context) {
	themeWarn()
}

// WindowSetLightTheme mirrors the v2 runtime.WindowSetLightTheme function.
// v3 has no runtime theme setter: the theme is configured at window creation
// via application.WebviewWindowOptions. This is a no-op that logs a warning once.
func WindowSetLightTheme(_ context.Context) {
	themeWarn()
}

// WindowSetDarkTheme mirrors the v2 runtime.WindowSetDarkTheme function.
// v3 has no runtime theme setter: the theme is configured at window creation
// via application.WebviewWindowOptions. This is a no-op that logs a warning once.
func WindowSetDarkTheme(_ context.Context) {
	themeWarn()
}

// WindowShow mirrors the v2 runtime.WindowShow function.
// v3 equivalent: window.Show.
func WindowShow(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.Show()
	}
}

// WindowHide mirrors the v2 runtime.WindowHide function.
// v3 equivalent: window.Hide.
func WindowHide(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.Hide()
	}
}

// WindowSetSize mirrors the v2 runtime.WindowSetSize function.
// v3 equivalent: window.SetSize.
func WindowSetSize(_ context.Context, width int, height int) {
	if w := currentWindow(); w != nil {
		w.SetSize(width, height)
	}
}

// WindowGetSize mirrors the v2 runtime.WindowGetSize function.
// v3 equivalent: window.Size.
func WindowGetSize(_ context.Context) (int, int) {
	if w := currentWindow(); w != nil {
		return w.Size()
	}
	return 0, 0
}

// WindowSetMinSize mirrors the v2 runtime.WindowSetMinSize function.
// v3 equivalent: window.SetMinSize.
func WindowSetMinSize(_ context.Context, width int, height int) {
	if w := currentWindow(); w != nil {
		w.SetMinSize(width, height)
	}
}

// WindowSetMaxSize mirrors the v2 runtime.WindowSetMaxSize function.
// v3 equivalent: window.SetMaxSize.
func WindowSetMaxSize(_ context.Context, width int, height int) {
	if w := currentWindow(); w != nil {
		w.SetMaxSize(width, height)
	}
}

// WindowSetAlwaysOnTop mirrors the v2 runtime.WindowSetAlwaysOnTop function.
// v3 equivalent: window.SetAlwaysOnTop.
func WindowSetAlwaysOnTop(_ context.Context, b bool) {
	if w := currentWindow(); w != nil {
		w.SetAlwaysOnTop(b)
	}
}

// WindowSetPosition mirrors the v2 runtime.WindowSetPosition function.
// v3 equivalent: window.SetRelativePosition (v2 positions were relative to
// the window's current screen).
func WindowSetPosition(_ context.Context, x int, y int) {
	if w := currentWindow(); w != nil {
		w.SetRelativePosition(x, y)
	}
}

// WindowGetPosition mirrors the v2 runtime.WindowGetPosition function.
// v3 equivalent: window.RelativePosition.
func WindowGetPosition(_ context.Context) (int, int) {
	if w := currentWindow(); w != nil {
		return w.RelativePosition()
	}
	return 0, 0
}

// WindowMaximise mirrors the v2 runtime.WindowMaximise function.
// v3 equivalent: window.Maximise.
func WindowMaximise(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.Maximise()
	}
}

// WindowToggleMaximise mirrors the v2 runtime.WindowToggleMaximise function.
// v3 equivalent: window.ToggleMaximise.
func WindowToggleMaximise(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.ToggleMaximise()
	}
}

// WindowUnmaximise mirrors the v2 runtime.WindowUnmaximise function.
// v3 equivalent: window.UnMaximise.
func WindowUnmaximise(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.UnMaximise()
	}
}

// WindowMinimise mirrors the v2 runtime.WindowMinimise function.
// v3 equivalent: window.Minimise.
func WindowMinimise(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.Minimise()
	}
}

// WindowUnminimise mirrors the v2 runtime.WindowUnminimise function.
// v3 equivalent: window.UnMinimise.
func WindowUnminimise(_ context.Context) {
	if w := currentWindow(); w != nil {
		w.UnMinimise()
	}
}

// WindowIsFullscreen mirrors the v2 runtime.WindowIsFullscreen function.
// v3 equivalent: window.IsFullscreen.
func WindowIsFullscreen(_ context.Context) bool {
	if w := currentWindow(); w != nil {
		return w.IsFullscreen()
	}
	return false
}

// WindowIsMaximised mirrors the v2 runtime.WindowIsMaximised function.
// v3 equivalent: window.IsMaximised.
func WindowIsMaximised(_ context.Context) bool {
	if w := currentWindow(); w != nil {
		return w.IsMaximised()
	}
	return false
}

// WindowIsMinimised mirrors the v2 runtime.WindowIsMinimised function.
// v3 equivalent: window.IsMinimised.
func WindowIsMinimised(_ context.Context) bool {
	if w := currentWindow(); w != nil {
		return w.IsMinimised()
	}
	return false
}

// WindowIsNormal mirrors the v2 runtime.WindowIsNormal function.
// v3 equivalent: !window.IsFullscreen() && !window.IsMaximised() && !window.IsMinimised().
func WindowIsNormal(_ context.Context) bool {
	w := currentWindow()
	if w == nil {
		return false
	}
	return !w.IsFullscreen() && !w.IsMaximised() && !w.IsMinimised()
}

// WindowExecJS mirrors the v2 runtime.WindowExecJS function.
// v3 equivalent: window.ExecJS.
func WindowExecJS(_ context.Context, js string) {
	if w := currentWindow(); w != nil {
		w.ExecJS(js)
	}
}

// WindowSetBackgroundColour mirrors the v2 runtime.WindowSetBackgroundColour function.
// v3 equivalent: window.SetBackgroundColour(application.RGBA{...}).
func WindowSetBackgroundColour(_ context.Context, R, G, B, A uint8) {
	if w := currentWindow(); w != nil {
		w.SetBackgroundColour(application.RGBA{Red: R, Green: G, Blue: B, Alpha: A})
	}
}

// WindowPrint mirrors the v2 runtime.WindowPrint function.
// v3 equivalent: window.Print.
func WindowPrint(_ context.Context) {
	if w := currentWindow(); w != nil {
		if err := w.Print(); err != nil {
			logger().Warn("v2compat: WindowPrint failed", "error", err)
		}
	}
}
