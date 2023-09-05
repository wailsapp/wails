package runtime

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/options"
)

// WindowSetTitle sets the title of the window
func WindowSetTitle(ctx context.Context, title string) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetTitle(title)
}

// WindowFullscreen makes the window fullscreen
func WindowFullscreen(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowFullscreen()
}

// WindowUnfullscreen makes the window UnFullscreen
func WindowUnfullscreen(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowUnfullscreen()
}

// WindowCenter the window on the current screen
func WindowCenter(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowCenter()
}

// WindowReload will reload the window contents
func WindowReload(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowReload()
}

// WindowReloadApp will reload the application
func WindowReloadApp(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowReloadApp()
}

func WindowSetSystemDefaultTheme(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetSystemDefaultTheme()
}

func WindowSetLightTheme(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetLightTheme()
}

func WindowSetDarkTheme(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetDarkTheme()
}

// WindowShow shows the window if hidden
func WindowShow(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowShow()
}

// WindowHide the window
func WindowHide(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowHide()
}

// WindowSetSize sets the size of the window
func WindowSetSize(ctx context.Context, width int, height int) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetSize(width, height)
}

func WindowGetSize(ctx context.Context) (int, int) {
	appFrontend := getFrontend(ctx)
	return appFrontend.WindowGetSize()
}

// WindowSetMinSize sets the minimum size of the window
func WindowSetMinSize(ctx context.Context, width int, height int) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetMinSize(width, height)
}

// WindowSetMaxSize sets the maximum size of the window
func WindowSetMaxSize(ctx context.Context, width int, height int) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetMaxSize(width, height)
}

// WindowSetAlwaysOnTop sets the window AlwaysOnTop or not on top
func WindowSetAlwaysOnTop(ctx context.Context, b bool) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetAlwaysOnTop(b)
}

// WindowSetPosition sets the position of the window
func WindowSetPosition(ctx context.Context, x int, y int) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetPosition(x, y)
}

func WindowGetPosition(ctx context.Context) (int, int) {
	appFrontend := getFrontend(ctx)
	return appFrontend.WindowGetPosition()
}

// WindowMaximise the window
func WindowMaximise(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowMaximise()
}

// WindowToggleMaximise the window
func WindowToggleMaximise(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowToggleMaximise()
}

// WindowUnmaximise the window
func WindowUnmaximise(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowUnmaximise()
}

// WindowMinimise the window
func WindowMinimise(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowMinimise()
}

// WindowUnminimise the window
func WindowUnminimise(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowUnminimise()
}

// WindowIsFullscreen get the window state is window Fullscreen
func WindowIsFullscreen(ctx context.Context) bool {
	appFrontend := getFrontend(ctx)
	return appFrontend.WindowIsFullscreen()
}

// WindowIsMaximised get the window state is window Maximised
func WindowIsMaximised(ctx context.Context) bool {
	appFrontend := getFrontend(ctx)
	return appFrontend.WindowIsMaximised()
}

// WindowIsMinimised get the window state is window Minimised
func WindowIsMinimised(ctx context.Context) bool {
	appFrontend := getFrontend(ctx)
	return appFrontend.WindowIsMinimised()
}

// WindowIsNormal get the window state is window Normal
func WindowIsNormal(ctx context.Context) bool {
	appFrontend := getFrontend(ctx)
	return appFrontend.WindowIsNormal()
}

// WindowExecJS executes the given Js in the window
func WindowExecJS(ctx context.Context, js string) {
	appFrontend := getFrontend(ctx)
	appFrontend.ExecJS(js)
}

func WindowSetBackgroundColour(ctx context.Context, R, G, B, A uint8) {
	appFrontend := getFrontend(ctx)
	col := &options.RGBA{
		R: R,
		G: G,
		B: B,
		A: A,
	}
	appFrontend.WindowSetBackgroundColour(col)
}

func WindowPrint(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowPrint()
}
