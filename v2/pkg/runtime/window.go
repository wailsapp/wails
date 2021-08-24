package runtime

import (
	"context"
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

// WindowUnFullscreen makes the window UnFullscreen
func WindowUnFullscreen(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowUnFullscreen()
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

// WindowSetPosition sets the position of the window
func WindowSetPosition(ctx context.Context, x int, y int) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowSetPos(x, y)
}

// WindowMaximise the window
func WindowMaximise(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.WindowMaximise()
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
