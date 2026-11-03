package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// The API scenarios drive what the input scenarios cannot reach: the role
// menus' own items, menu items changed at run time, and the window API. They
// still check results the way an app would see them (getters, events, the
// page reporting back), and they run on every OS.

func init() {
	scenarios = append(scenarios,
		scenario{name: "menu-api", run: menuAPI},
		scenario{name: "window-api", run: windowAPI},
		scenario{name: "runtime-api", run: runtimeAPI},
		scenario{name: "second-instance", run: secondInstance},
		scenario{name: "tray", os: "linux", run: trayHost},
		// Last: it reloads the page (see menuRoles).
		scenario{name: "menu-roles", run: menuRoles},
	)
}

// eventually polls cond for up to three seconds.
func eventually(cond func() bool) bool { return within(3*time.Second, cond) }

// within polls cond for up to d.
func within(d time.Duration, cond func() bool) bool {
	for deadline := time.Now().Add(d); time.Now().Before(deadline); time.Sleep(25 * time.Millisecond) {
		if cond() {
			return true
		}
	}
	return cond()
}

// menuRoles presses the View menu's role shortcuts, as a user would.
func menuRoles(p *probe) error {
	zoom := func() float64 { return p.main.GetZoom() }
	if err := press(p.d, title, "cmdorctrl+0"); err != nil {
		return err
	}
	if !eventually(func() bool { return zoom() == 1 }) {
		return fmt.Errorf("Actual Size (Ctrl+0): zoom %v, want 1", zoom())
	}
	// "Zoom In" is CmdOrCtrl+plus: the "+" character, Shift+= on a US keyboard.
	if err := press(p.d, title, "cmdorctrl+shift+="); err != nil {
		return err
	}
	if !eventually(func() bool { return zoom() > 1 }) {
		return fmt.Errorf("Zoom In (Ctrl++): zoom %v, want more than 1", zoom())
	}
	// Zoom Out from the zoomed-in level: Windows and macOS stop at 100%,
	// GTK goes below it.
	zoomed := zoom()
	if err := press(p.d, title, "cmdorctrl+-"); err != nil {
		return err
	}
	if !eventually(func() bool { return zoom() < zoomed }) {
		return fmt.Errorf("Zoom Out (Ctrl+-): zoom %v, want less than %v", zoom(), zoomed)
	}
	if err := press(p.d, title, "cmdorctrl+0"); err != nil {
		return err
	}
	if !eventually(func() bool { return zoom() == 1 }) {
		return fmt.Errorf("Actual Size after zooming: zoom %v, want 1", zoom())
	}
	// Reload last: the page loads again and says so. (On macOS, menu items
	// on Cmd+Shift+letter stop firing after a reload while the page has
	// focus, but only for the driver's posted key events: with a hardware
	// keyboard (a VM's VNC input) they fire after one or two reloads, in the
	// page and in a text field. So nothing runs after this scenario.)
	drain()
	if err := press(p.d, title, "cmdorctrl+r"); err != nil {
		return err
	}
	if err := await("ready main", 10*time.Second); err != nil {
		return fmt.Errorf("Reload (Ctrl+R): %w", err)
	}
	return nil
}

// menuAPI changes a menu item at run time: a hidden item's shortcut does
// nothing, shown again it works, and label, tooltip and bitmap changes leave
// it working.
func menuAPI(p *probe) error {
	item := p.hidden
	item.SetLabel("Shown").SetTooltip("A tooltip").SetBitmap(bitmap())
	application.InvokeSync(p.menu.Update)
	if err := press(p.d, title, "cmdorctrl+shift+l"); err != nil {
		return err
	}
	if err := await("hidden-action", 5*time.Second); err != nil {
		shot("menu-api")
		return fmt.Errorf("after label, tooltip and bitmap changes: %w (window visible %v, focused %v)",
			err, p.main.IsVisible(), p.main.IsFocused())
	}
	item.SetHidden(true)
	application.InvokeSync(p.menu.Update)
	if !item.Hidden() {
		return fmt.Errorf("Hidden() is false after SetHidden(true)")
	}
	if err := press(p.d, title, "cmdorctrl+shift+l"); err != nil {
		return err
	}
	if err := awaitNone("hidden-action"); err != nil {
		return fmt.Errorf("a hidden item's shortcut: %w", err)
	}
	item.SetHidden(false)
	application.InvokeSync(p.menu.Update)
	if err := press(p.d, title, "cmdorctrl+shift+l"); err != nil {
		return err
	}
	if err := await("hidden-action", 5*time.Second); err != nil {
		return fmt.Errorf("shown again: %w", err)
	}
	return nil
}

// bitmap is a small PNG for SetBitmap.
func bitmap() []byte {
	img := image.NewNRGBA(image.Rect(0, 0, 16, 16))
	for i := 0; i < 16; i++ {
		img.Set(i, i, color.NRGBA{R: 124, G: 58, B: 237, A: 255})
	}
	var b bytes.Buffer
	_ = png.Encode(&b, img)
	return b.Bytes()
}

// windowStatesChange reports whether minimise, maximise and fullscreen take
// effect: a bare X server has no window manager to grant them.
func windowStatesChange() bool {
	return runtime.GOOS != "linux" || os.Getenv("E2E_WM") != ""
}

// windowAPI drives a window through the public Window API and checks each
// change is visible through the getters, events or the page.
func windowAPI(p *probe) error {
	const apiTitle = "E2E api"
	w := p.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name: "api", Title: apiTitle, Width: 500, Height: 350, X: 200, Y: 200, URL: "/?w=api",
	})
	closed := make(chan struct{}, 1)
	w.OnWindowEvent(events.Common.WindowClosing, func(*application.WindowEvent) {
		select {
		case closed <- struct{}{}:
		default:
		}
	})
	defer func() {
		w.Close()
		select {
		case <-closed:
		case <-time.After(5 * time.Second):
		}
	}()
	if err := await("ready api", 15*time.Second); err != nil {
		return err
	}
	check := func(what string, cond func() bool, got func() any) error {
		if !eventually(cond) {
			return fmt.Errorf("%s: got %v", what, got())
		}
		return nil
	}

	w.SetTitle("E2E api renamed")
	w.SetTitle(apiTitle)

	w.SetSize(520, 360)
	if err := check("SetSize(520, 360)", func() bool { a, b := w.Size(); return a == 520 && b == 360 },
		func() any { a, b := w.Size(); return fmt.Sprint(a, "x", b) }); err != nil {
		return err
	}
	if w.Width() != 520 || w.Height() != 360 {
		return fmt.Errorf("Width/Height %dx%d disagree with Size", w.Width(), w.Height())
	}
	w.SetPosition(220, 230)
	if err := check("SetPosition(220, 230)", func() bool { x, y := w.Position(); return x == 220 && y == 230 },
		func() any { x, y := w.Position(); return fmt.Sprint(x, ",", y) }); err != nil {
		return err
	}
	want := application.Rect{X: 240, Y: 250, Width: 540, Height: 380}
	w.SetBounds(want)
	if err := check(fmt.Sprintf("SetBounds(%+v)", want), func() bool { return w.Bounds() == want },
		func() any { return w.Bounds() }); err != nil {
		return err
	}
	// Size constraints limit the user's resizing everywhere, but only GTK
	// also clamps SetSize, so they are exercised, not checked.
	w.SetMaxSize(600, 400)
	w.SetSize(800, 600)
	w.SetMaxSize(0, 0)
	w.SetMinSize(300, 200)
	w.SetSize(100, 100)
	w.SetMinSize(0, 0)
	w.SetSize(520, 360)

	w.SetResizable(false)
	if w.Resizable() {
		return fmt.Errorf("Resizable() after SetResizable(false)")
	}
	w.SetResizable(true)

	w.SetZoom(1.5)
	if z := w.GetZoom(); z != 1.5 {
		return fmt.Errorf("GetZoom() = %v after SetZoom(1.5)", z)
	}
	w.ZoomIn()
	if !eventually(func() bool { return w.GetZoom() > 1.5 }) {
		return fmt.Errorf("ZoomIn from 1.5: zoom %v", w.GetZoom())
	}
	w.ZoomOut()
	w.ZoomReset()
	if !eventually(func() bool { return w.GetZoom() == 1 }) {
		return fmt.Errorf("ZoomReset: zoom %v", w.GetZoom())
	}

	w.SetIgnoreMouseEvents(true)
	if !w.IsIgnoreMouseEvents() {
		return fmt.Errorf("IsIgnoreMouseEvents() after SetIgnoreMouseEvents(true)")
	}
	// A click on a window that ignores the mouse goes to whatever is below
	// it. Only where the driver's clicks are the OS's own pointer input: the
	// macOS driver posts its events to the window itself.
	physicalPointer := runtime.GOOS != "darwin"
	if physicalPointer {
		drain()
		if err := p.d.click(apiTitle, 100, 100); err != nil {
			return err
		}
		if err := awaitNone("click api"); err != nil {
			return fmt.Errorf("SetIgnoreMouseEvents(true): the window still took a click (%w)", err)
		}
	}
	w.SetIgnoreMouseEvents(false)
	if physicalPointer {
		if err := p.d.click(apiTitle, 100, 100); err != nil {
			return err
		}
		if err := await("click api", 5*time.Second); err != nil {
			return fmt.Errorf("SetIgnoreMouseEvents(false): %w", err)
		}
	}
	w.SetAlwaysOnTop(true)
	w.SetAlwaysOnTop(false)
	w.SetBackgroundColour(application.RGBA{Red: 30, Green: 41, Blue: 59, Alpha: 255})
	for _, state := range []application.ButtonState{application.ButtonDisabled, application.ButtonHidden, application.ButtonEnabled} {
		w.SetMinimiseButtonState(state)
		w.SetMaximiseButtonState(state)
		w.SetCloseButtonState(state)
		w.SetFullscreenButtonState(state)
	}
	w.SetEnabled(false)
	w.SetEnabled(true)
	w.Center()
	if s, err := w.GetScreen(); err != nil || s == nil {
		return fmt.Errorf("GetScreen: %v, %v", s, err)
	}
	_ = w.GetBorderSizes()
	x, y := w.RelativePosition()
	w.SetRelativePosition(x, y)

	w.Hide()
	if err := check("Hide", func() bool { return !w.IsVisible() }, func() any { return w.IsVisible() }); err != nil {
		return err
	}
	w.Show()
	if err := check("Show", w.IsVisible, func() any { return w.IsVisible() }); err != nil {
		return err
	}

	// The page runs JavaScript sent from Go, reloads, and follows SetURL.
	drain()
	w.ExecJS(`fetch('/input?w=api&id=js&v=ran')`)
	if err := await("input api js=ran", 5*time.Second); err != nil {
		return fmt.Errorf("ExecJS: %w", err)
	}
	w.Reload()
	if err := await("ready api", 10*time.Second); err != nil {
		return fmt.Errorf("Reload: %w", err)
	}
	w.SetURL("/?w=api2")
	if err := await("ready api2", 10*time.Second); err != nil {
		return fmt.Errorf("SetURL: %w", err)
	}

	w.Focus()
	if err := check("Focus", w.IsFocused, func() any { return w.IsFocused() }); err != nil {
		return err
	}

	// Window states: always exercised, checked where something grants them.
	states := windowStatesChange()
	expect := func(what string, cond func() bool) error {
		// macOS animates full screen in and out over about a second, and
		// ignores a change requested while one is running.
		if states && !within(10*time.Second, cond) {
			return fmt.Errorf("%s did not take effect", what)
		}
		return nil
	}
	for _, step := range []struct {
		what string
		do   func()
		cond func() bool
	}{
		{"Maximise", func() { w.Maximise() }, w.IsMaximised},
		{"UnMaximise", w.UnMaximise, func() bool { return !w.IsMaximised() }},
		{"ToggleMaximise", w.ToggleMaximise, w.IsMaximised},
		{"ToggleMaximise back", w.ToggleMaximise, func() bool { return !w.IsMaximised() }},
		{"Fullscreen", func() { w.Fullscreen() }, w.IsFullscreen},
		{"UnFullscreen", w.UnFullscreen, func() bool { return !w.IsFullscreen() }},
		{"Minimise", func() { w.Minimise() }, w.IsMinimised},
		{"UnMinimise", w.UnMinimise, func() bool { return !w.IsMinimised() }},
		{"Minimise again", func() { w.Minimise() }, w.IsMinimised},
		{"Restore", w.Restore, func() bool { return !w.IsMinimised() && !w.IsMaximised() && !w.IsFullscreen() }},
	} {
		step.do()
		if err := expect(step.what, step.cond); err != nil {
			return err
		}
		if runtime.GOOS == "darwin" {
			// Kept: the state reads as changed before macOS's animation
			// ends, and a change requested during it is ignored. Nothing
			// observable marks the animation's end.
			time.Sleep(1500 * time.Millisecond)
		}
	}
	return nil
}
