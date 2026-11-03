// The e2e probe drives a real Wails app with real OS input: key bindings
// (application-wide and per window), menu accelerators, checkbox, radio,
// disabled and re-bound menu items, Edit-menu shortcuts in the page, typed
// text, keyboard menu navigation, a context menu, window focus and close
// events, Go to JavaScript events, message dialogs and the open and save
// file dialogs. It runs on the lab VMs and in CI (e2e-v3.yml), never on a
// developer's console session:
//
//	macOS:   go run -tags purego ./tests/purego/e2e
//	Linux:   xvfb-run -a -s "-screen 0 1280x800x24" dbus-run-session -- go run -tags purego ./tests/purego/e2e
//	Windows: go run ./tests/purego/e2e (in the logged-on desktop session)
//
// Each scenario prints "E2E <name> ok" or "E2E <name> fail: <reason>", and a
// screenshot of any open menu or dialog between "EVIDENCE-BEGIN <name>.png"
// and "EVIDENCE-END" (base64); with E2E_EVIDENCE_DIR set the screenshots are
// also written there. Build it with -cover and set GOCOVERDIR to count what
// it exercises.
//
// The environment selects what runs:
//
//	E2E_SCENARIOS  comma separated scenarios to run (default: all for this OS)
//	E2E_SKIP       comma separated scenarios to leave out
//	E2E_LOCKS      lock keys to switch on first, e.g. "numlock" or "numlock,capslock"
package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const (
	title       = "E2E"
	secondTitle = "E2E second"
)

// The page: a context-menu target, two inputs and a textarea at known
// places. It reports what the user did to it (input, pointer motion) and
// what Go sent it (events) back to the asset server, which turns each report
// into a signal.
const page = `<!doctype html><html><head><script type="module" src="/wails/runtime.js"></script></head>
<body style="margin:0;font:16px system-ui,-apple-system,sans-serif;background:#1e293b;color:white">
<div id="ctx" style="position:absolute;left:20px;top:20px;width:240px;height:100px;background:#7c3aed;
display:flex;align-items:center;justify-content:center;--custom-contextmenu:e2e">Right-click me</div>
<p style="position:absolute;left:20px;top:140px">E2E probe</p>
<input id="src" value="wails-e2e-clip" style="position:absolute;left:300px;top:20px;width:200px;height:30px">
<input id="dst" style="position:absolute;left:300px;top:70px;width:200px;height:30px">
<textarea id="area" style="position:absolute;left:300px;top:120px;width:200px;height:60px"></textarea>
<script type="module">
import {Events} from '/wails/runtime.js';
const w = new URLSearchParams(location.search).get('w') || 'main';
const report = (path, params) => fetch(path + '?' + new URLSearchParams({w, ...params}));
document.addEventListener('input', (e) => report('/input', {id: e.target.id, v: e.target.value}));
document.addEventListener('mousemove', (e) => report('/mouse', {x: e.clientX, y: e.clientY}));
document.addEventListener('mousedown', () => report('/click', {}));
Events.On('e2e-ping', (e) => report('/event', {v: JSON.stringify(e.data)}));
report('/ready', {dpr: devicePixelRatio});
</script></body></html>`

// The page's controls' centres, in the page's coordinates.
const (
	srcX, srcY   = 400, 35
	dstX, dstY   = 400, 85
	areaX, areaY = 400, 150
)

// signals carries what the app saw, as text: "keybinding main",
// "input main dst=…", "ready second". It is buffered generously and never
// blocks a sender: scenarios drain it before they act.
var signals = make(chan string, 1024)

// pointer carries the page's pointer reports, apart from signals so the
// stream of mousemove reports cannot crowd them out.
var pointer = make(chan pointerReport, 64)

type pointerReport struct {
	window string
	x, y   float64
}

// dpr is each page's devicePixelRatio, from its ready report.
var dpr sync.Map

func emit(s string) {
	select {
	case signals <- s:
	default:
		fmt.Fprintf(os.Stderr, "e2e: signal buffer full, dropped %q\n", s)
	}
}

// eventPause paces the input of a multi-step scenario: the page has no signal
// for a click having moved focus or a selection having been made.
const eventPause = 300 * time.Millisecond

// quiet is how long a scenario waits to be sure something did not happen.
const quiet = 1500 * time.Millisecond

func signal(name string) func() { return func() { emit(name) } }

// await waits for a named signal, discarding others.
func await(name string, timeout time.Duration) error {
	deadline := time.After(timeout)
	for {
		select {
		case got := <-signals:
			if got == name {
				return nil
			}
		case <-deadline:
			return fmt.Errorf("no %q within %s", name, timeout)
		}
	}
}

// awaitLatest waits for the signal want, remembering the last signal that
// starts with prefix, so a failure can say what arrived instead.
func awaitLatest(want, prefix string, timeout time.Duration) error {
	deadline := time.After(timeout)
	last := "nothing"
	for {
		select {
		case got := <-signals:
			if got == want {
				return nil
			}
			if strings.HasPrefix(got, prefix) {
				last = fmt.Sprintf("%q", got)
			}
		case <-deadline:
			return fmt.Errorf("no %q within %s; last %s", want, timeout, last)
		}
	}
}

// awaitNone fails if the named signal arrives within quiet.
func awaitNone(name string) error {
	if await(name, quiet) == nil {
		return fmt.Errorf("unexpected %q", name)
	}
	return nil
}

// drain discards pending signals, so a scenario sees only its own.
func drain() {
	for {
		select {
		case <-signals:
		default:
			return
		}
	}
}

// steps runs input steps in order, pausing between them.
func steps(fns ...func() error) error {
	for _, step := range fns {
		if err := step(); err != nil {
			return err
		}
		time.Sleep(eventPause)
	}
	return nil
}

// press focuses the window titled win and presses chord.
func press(d driver, win, chord string) error {
	if err := d.focus(win); err != nil {
		return err
	}
	return d.key(chord)
}

// probe holds what the scenarios share: the app, its windows and menu items.
type probe struct {
	app      *application.App
	d        driver
	main     *application.WebviewWindow
	menu     *application.Menu
	toggle   *application.MenuItem
	alpha    *application.MenuItem
	beta     *application.MenuItem
	disabled *application.MenuItem
	rebound  *application.MenuItem
	hidden   *application.MenuItem
}

type scenario struct {
	name string
	os   string // comma separated GOOS values; "" = every OS
	run  func(p *probe) error
}

func (s scenario) runsOn(goos string) bool {
	return s.os == "" || slices.Contains(strings.Split(s.os, ","), goos)
}

var scenarios = []scenario{
	// calibrate maps page coordinates to the driver's pointer coordinates
	// from where the page reports the pointer. It runs first: the other
	// pointer scenarios click through the mapping it finds.
	{name: "calibrate", os: "linux,windows", run: func(p *probe) error {
		return calibrate(p.d, title, "main")
	}},
	{name: "keybinding", run: func(p *probe) error {
		if err := press(p.d, title, "cmdorctrl+shift+k"); err != nil {
			return err
		}
		return await("keybinding main", 5*time.Second)
	}},
	{name: "menu-accelerator", run: func(p *probe) error {
		if err := press(p.d, title, "cmdorctrl+m"); err != nil {
			return err
		}
		return await("menu-action", 5*time.Second)
	}},
	{name: "menu-checkbox", run: func(p *probe) error {
		// The accelerator toggles the item: the click handler sees the new
		// state, and so does Checked.
		for _, want := range []bool{!p.toggle.Checked(), p.toggle.Checked()} {
			if err := press(p.d, title, "cmdorctrl+shift+b"); err != nil {
				return err
			}
			if err := await(fmt.Sprintf("toggle %v", want), 5*time.Second); err != nil {
				return err
			}
			if got := p.toggle.Checked(); got != want {
				return fmt.Errorf("Checked() = %v after the click, want %v", got, want)
			}
		}
		return nil
	}},
	{name: "menu-radio", run: func(p *probe) error {
		// Choosing a radio item by its accelerator checks it and unchecks
		// the rest of its group. (Not Ctrl+Shift+U: GTK's input method takes
		// it to start Unicode entry.)
		for _, c := range []struct {
			chord, name string
			on, off     *application.MenuItem
		}{{"cmdorctrl+shift+g", "Beta", p.beta, p.alpha}, {"cmdorctrl+shift+y", "Alpha", p.alpha, p.beta}} {
			if err := press(p.d, title, c.chord); err != nil {
				return err
			}
			if err := await("radio "+c.name, 5*time.Second); err != nil {
				return err
			}
			if !c.on.Checked() || c.off.Checked() {
				return fmt.Errorf("after choosing %s: %s checked=%v, %s checked=%v",
					c.name, c.on.Label(), c.on.Checked(), c.off.Label(), c.off.Checked())
			}
		}
		return nil
	}},
	{name: "menu-disabled", run: func(p *probe) error {
		// A disabled item's accelerator does nothing; enabled again at run
		// time, it works.
		if err := press(p.d, title, "cmdorctrl+shift+d"); err != nil {
			return err
		}
		if err := awaitNone("disabled-action"); err != nil {
			return fmt.Errorf("disabled item: %w", err)
		}
		p.disabled.SetEnabled(true)
		application.InvokeSync(p.menu.Update)
		defer func() {
			p.disabled.SetEnabled(false)
			application.InvokeSync(p.menu.Update)
		}()
		if err := press(p.d, title, "cmdorctrl+shift+d"); err != nil {
			return err
		}
		return await("disabled-action", 5*time.Second)
	}},
	{name: "menu-rebind", run: func(p *probe) error {
		// Changing an accelerator at run time moves the binding: the old
		// chord stops working and the new one starts.
		if err := press(p.d, title, "cmdorctrl+shift+e"); err != nil {
			return err
		}
		if err := await("rebound-action", 5*time.Second); err != nil {
			return fmt.Errorf("before rebinding: %w", err)
		}
		p.rebound.SetAccelerator("CmdOrCtrl+Shift+O")
		application.InvokeSync(p.menu.Update)
		if err := press(p.d, title, "cmdorctrl+shift+o"); err != nil {
			return err
		}
		if err := await("rebound-action", 5*time.Second); err != nil {
			return fmt.Errorf("new accelerator: %w", err)
		}
		if err := press(p.d, title, "cmdorctrl+shift+e"); err != nil {
			return err
		}
		if err := awaitNone("rebound-action"); err != nil {
			return fmt.Errorf("old accelerator: %w", err)
		}
		if runtime.GOOS == "darwin" {
			return nil // the menu is the application's, not the window's
		}
		// Setting the menu the window already shows (the old way to make it
		// pick up changes) keeps it working.
		p.main.SetMenu(p.menu)
		if err := press(p.d, title, "cmdorctrl+shift+o"); err != nil {
			return err
		}
		if err := await("rebound-action", 5*time.Second); err != nil {
			return fmt.Errorf("after SetMenu with the same menu: %w", err)
		}
		return nil
	}},
	{name: "menu-navigate", os: "linux,windows", run: func(p *probe) error {
		// F10 opens the menu bar; the first item of the first menu is
		// "Action". This is the user's path through the native menu.
		if err := press(p.d, title, "f10"); err != nil {
			return err
		}
		// Kept: lets the menu bar paint before the screenshot. The driver has no
		// way to observe the open menu, and the keys below need it open.
		time.Sleep(time.Second)
		shot("menu-navigate")
		for _, k := range menuNavigateKeys(p.app) {
			if err := p.d.key(k); err != nil {
				return err
			}
		}
		return await("menu-action", 5*time.Second)
	}},
	{name: "context-menu", run: func(p *probe) error {
		if err := p.d.rightClick(title, 100, 60); err != nil {
			return err
		}
		// Kept: lets the popover paint before the screenshot. The driver cannot
		// observe it, and the keys below need it open.
		time.Sleep(time.Second)
		shot("context-menu")
		if contextMenuByMouse {
			// The popover opens at the pointer; its first item is just below
			// and to the right of it.
			if err := p.d.click(title, 120, 83); err != nil {
				return err
			}
			return await("context-action", 5*time.Second)
		}
		for _, k := range []string{"down", "return"} {
			if err := p.d.key(k); err != nil {
				return err
			}
		}
		return await("context-action", 5*time.Second)
	}},
	{name: "type-text", run: func(p *probe) error {
		// Typed text arrives exactly: mixed case, digits and the punctuation
		// a keyboard layout moves around (Shift on one layout, not another).
		const text = "Wails E2E: 0-9, x_y/z.txt?"
		// Focus first: on macOS a posted click does not make the window key,
		// and just after a context menu closes it may briefly not be.
		if err := p.d.focus(title); err != nil {
			return err
		}
		if err := steps(
			func() error { return p.d.click(title, dstX, dstY) },
			func() error { return p.d.key("cmdorctrl+a") },
			func() error { return p.d.key("delete") },
			func() error { return p.d.typeText(text) },
		); err != nil {
			return err
		}
		if err := awaitLatest("input main dst="+text, "input main dst=", 5*time.Second); err != nil {
			return err
		}
		// Return in a textarea is a new line: no menu, binding or dialog
		// default takes it.
		if err := steps(
			func() error { return p.d.click(title, areaX, areaY) },
			func() error { return p.d.typeText("one") },
			func() error { return p.d.key("return") },
			func() error { return p.d.typeText("two") },
		); err != nil {
			return err
		}
		return awaitLatest("input main area=one\ntwo", "input main area=", 5*time.Second)
	}},
	{name: "edit-copy-paste", run: func(p *probe) error {
		// The Edit menu binds Ctrl+A/C/V/X/Z (Cmd on macOS). Those shortcuts must
		// still edit the page: GTK3 used to swallow them in empty role bodies.
		if err := p.d.focus(title); err != nil {
			return err
		}
		if err := steps(
			func() error { return p.d.click(title, dstX, dstY) },
			func() error { return p.d.key("cmdorctrl+a") },
			func() error { return p.d.key("delete") },
			func() error { return p.d.click(title, srcX, srcY) },
			func() error { return p.d.key("cmdorctrl+a") },
			func() error { return p.d.key("cmdorctrl+c") },
			func() error { return p.d.click(title, dstX, dstY) },
			func() error { return p.d.key("cmdorctrl+v") },
		); err != nil {
			return err
		}
		if err := await("input main dst=wails-e2e-clip", 5*time.Second); err != nil {
			return fmt.Errorf("paste: %w", err)
		}
		if err := steps(
			func() error { return p.d.click(title, srcX, srcY) },
			func() error { return p.d.key("cmdorctrl+a") },
			func() error { return p.d.key("cmdorctrl+x") },
		); err != nil {
			return err
		}
		if err := await("input main src=", 5*time.Second); err != nil {
			return fmt.Errorf("cut: %w", err)
		}
		if text, ok := p.app.Clipboard.Text(); !ok || text != "wails-e2e-clip" {
			return fmt.Errorf("clipboard after cut = %q, %v", text, ok)
		}
		if err := p.d.key("cmdorctrl+z"); err != nil {
			return err
		}
		if err := await("input main src=wails-e2e-clip", 5*time.Second); err != nil {
			return fmt.Errorf("undo: %w", err)
		}
		return nil
	}},
	{name: "event-to-page", run: func(p *probe) error {
		p.app.Event.Emit("e2e-ping", "pong")
		return await(`event main "pong"`, 5*time.Second)
	}},
	{name: "second-window", run: func(p *probe) error {
		// A window opened at run time loads, takes focus when the user gives
		// it focus (and the first window loses it), has its own key bindings
		// as well as the application's, and closes without ending the app.
		second := p.app.Window.NewWithOptions(application.WebviewWindowOptions{
			Name: "second", Title: secondTitle, Width: 600, Height: 400, X: 300, Y: 250,
			URL: "/?w=second",
			KeyBindings: map[string]func(application.Window){
				"CmdOrCtrl+Shift+J": func(w application.Window) { emit("second-binding " + w.Name()) },
			},
		})
		for _, w := range []*application.WebviewWindow{p.main, second} {
			name := w.Name()
			w.OnWindowEvent(events.Common.WindowFocus, func(*application.WindowEvent) { emit("focus " + name) })
		}
		second.OnWindowEvent(events.Common.WindowClosing, func(*application.WindowEvent) { emit("closing second") })
		if err := await("ready second", 15*time.Second); err != nil {
			return err
		}
		if _, ok := p.d.(calibrator); ok {
			if err := calibrate(p.d, secondTitle, "second"); err != nil {
				return err
			}
		}
		// A new window may or may not take focus when it maps: give focus to
		// the first window, so the second's focus event below comes from the
		// user moving focus. (Kept: the pause lets that focus change land; it
		// may produce no event if the first window still had focus.)
		if err := p.d.focus(title); err != nil {
			return err
		}
		time.Sleep(eventPause)
		drain()
		if err := p.d.focus(secondTitle); err != nil {
			return err
		}
		if err := await("focus second", 5*time.Second); err != nil {
			return err
		}
		if err := p.d.key("cmdorctrl+shift+j"); err != nil {
			return err
		}
		if err := await("second-binding second", 5*time.Second); err != nil {
			return fmt.Errorf("the window's own binding: %w", err)
		}
		if err := p.d.key("cmdorctrl+shift+k"); err != nil {
			return err
		}
		if err := await("keybinding second", 5*time.Second); err != nil {
			return fmt.Errorf("the application's binding in the second window: %w", err)
		}
		if err := p.d.focus(title); err != nil {
			return err
		}
		if err := await("focus main", 5*time.Second); err != nil {
			return err
		}
		if err := p.d.key("cmdorctrl+shift+j"); err != nil {
			return err
		}
		if err := awaitNone("second-binding second"); err != nil {
			return fmt.Errorf("the second window's binding fired in the first window: %w", err)
		}
		second.Close()
		if err := await("closing second", 5*time.Second); err != nil {
			return err
		}
		for deadline := time.Now().Add(5 * time.Second); ; time.Sleep(50 * time.Millisecond) {
			if _, open := p.app.Window.GetByName("second"); !open {
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("the closed window is still registered")
			}
		}
		// The app is still alive: the first window still takes input.
		if err := press(p.d, title, "cmdorctrl+shift+k"); err != nil {
			return err
		}
		return await("keybinding main", 5*time.Second)
	}},
	{name: "dialog-default-button", run: func(p *probe) error {
		dialog := p.app.Dialog.Info().SetTitle("E2E dialog").SetMessage("Return presses the default button")
		dialog.AddButton("Later").OnClick(signal("dialog-later"))
		dialog.AddButton("OK").SetAsDefault().OnClick(signal("dialog-ok"))
		application.InvokeAsync(dialog.Show)
		// Kept: the driver cannot detect this dialog (GTK4 message dialogs are
		// untitled windows, and darwin focus reports success regardless), so
		// the probe waits for it to paint before the key is sent.
		time.Sleep(1500 * time.Millisecond)
		shot("dialog-default-button")
		if err := p.d.key("return"); err != nil {
			return err
		}
		return await("dialog-ok", 5*time.Second)
	}},
	{name: "dialog-cancel-button", run: func(p *probe) error {
		dialog := p.app.Dialog.Question().SetTitle("E2E question").SetMessage("Escape presses the cancel button")
		dialog.AddButton("Yes").SetAsDefault().OnClick(signal("question-yes"))
		dialog.AddButton("No").SetAsCancel().OnClick(signal("question-no"))
		application.InvokeAsync(dialog.Show)
		// Kept: as in dialog-default-button, the driver cannot detect this
		// dialog, so the probe waits for it to paint before the key is sent.
		time.Sleep(1500 * time.Millisecond)
		shot("dialog-cancel-button")
		if err := p.d.key("escape"); err != nil {
			return err
		}
		return await("question-no", 5*time.Second)
	}},
	{name: "open-file", run: func(p *probe) error {
		dir, err := os.MkdirTemp("", "wails-e2e")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		want := filepath.Join(dir, "chosen.txt")
		if err := os.WriteFile(want, []byte("e2e"), 0o644); err != nil {
			return err
		}
		done := prompt(func() (string, error) {
			return p.app.Dialog.OpenFile().SetTitle("E2E open").SetDirectory(dir).CanChooseFiles(true).PromptForSingleSelection()
		})
		// Kept: lets the open panel paint before the screenshot. The Linux path
		// below polls for the dialog; the macOS path has no readable signal.
		time.Sleep(2 * time.Second)
		shot("open-file")
		if c, ok := p.d.(dialogCanceller); ok {
			// macOS: the panel's contents run in an Apple service process
			// that posted events cannot reach, so the probe checks the panel
			// opens (screenshot) and that cancelling it returns cleanly.
			if err := c.cancelDialog(); err != nil {
				return err
			}
			return expectPath(done, "")
		}
		if c, ok := p.d.(fileChooser); ok {
			err = c.chooseFile("E2E open", want)
		} else {
			err = gtkChooseFile(p.d, "E2E open", want)
		}
		if err != nil {
			return err
		}
		return expectPath(done, want)
	}},
	{name: "save-file", run: func(p *probe) error {
		// The save dialog offers the name the app suggested in the folder it
		// chose; Return accepts it.
		dir, err := os.MkdirTemp("", "wails-e2e")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		done := prompt(func() (string, error) {
			dialog := p.app.Dialog.SaveFile()
			dialog.SetOptions(&application.SaveFileDialogOptions{
				Title: "E2E save", Directory: dir, Filename: "saved.txt", CanCreateDirectories: true,
			})
			return dialog.PromptForSingleSelection()
		})
		// Kept: as in open-file, lets the panel paint before the screenshot.
		time.Sleep(2 * time.Second)
		shot("save-file")
		if c, ok := p.d.(dialogCanceller); ok {
			if err := c.cancelDialog(); err != nil {
				return err
			}
			return expectPath(done, "")
		}
		if err := waitFocus(p.d, "E2E save"); err != nil {
			return err
		}
		if err := p.d.key("return"); err != nil {
			return err
		}
		return expectPath(done, filepath.Join(dir, "saved.txt"))
	}},
}

type promptResult struct {
	path string
	err  error
}

// prompt runs a blocking file dialog on its own goroutine.
func prompt(show func() (string, error)) chan promptResult {
	done := make(chan promptResult, 1)
	go func() {
		path, err := show()
		done <- promptResult{path, err}
	}()
	return done
}

// expectPath waits for a file dialog to return want ("" for cancelled).
func expectPath(done chan promptResult, want string) error {
	select {
	case r := <-done:
		if r.err != nil {
			return r.err
		}
		if want == "" {
			if r.path != "" {
				return fmt.Errorf("cancelled dialog returned %q; want \"\"", r.path)
			}
			return nil
		}
		if mustEval(filepath.Dir(r.path)) != mustEval(filepath.Dir(want)) || filepath.Base(r.path) != filepath.Base(want) {
			return fmt.Errorf("dialog returned %q, want %q", r.path, want)
		}
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("the file dialog did not return")
	}
}

// waitFocus focuses the window titled win once it exists. A dialog is its own
// top-level window; with no window manager it does not get the keyboard by
// itself.
func waitFocus(d driver, win string) error {
	err := d.focus(win)
	for deadline := time.Now().Add(5 * time.Second); err != nil && time.Now().Before(deadline); err = d.focus(win) {
		time.Sleep(100 * time.Millisecond)
	}
	return err
}

// gtkChooseFile types path into the GTK file chooser titled win.
func gtkChooseFile(d driver, win, path string) error {
	if err := waitFocus(d, win); err != nil {
		return err
	}
	// Type the path into the dialog's location field (Ctrl+L in GTK). The
	// pauses in these steps are kept: they pace the input, and GTK has no
	// signal the driver can read for when its pre-fill or completion is done.
	return steps(
		func() error { return d.key("ctrl+l") },
		// GTK4 pre-fills the field with the current folder: replace it.
		func() error { time.Sleep(500 * time.Millisecond); return d.key("ctrl+a") },
		func() error { return d.typeText(path) },
		// GTK completes the path inline while it is typed, inserting the
		// rest of the name after the cursor. Select from the cursor to the
		// end and delete it, so Return takes exactly what was typed.
		func() error { time.Sleep(300 * time.Millisecond); return d.key("shift+end") },
		func() error { return d.key("delete") },
		func() error { return d.key("return") },
	)
}

// calibrate finds how page coordinates map to the driver's pointer
// coordinates for the window titled win, whose page reports as page: it
// moves the pointer to two points and reads where the page saw it. That
// replaces guessing the menu bar's height (it depends on the theme and font)
// and the scale (GDK_SCALE, Windows display scaling), and checks the pointer
// reaches the page at all.
func calibrate(d driver, win, page string) error {
	c, ok := d.(calibrator)
	if !ok {
		return nil
	}
	if err := d.focus(win); err != nil {
		return err
	}
	at := func(x, y int) (pointerReport, error) {
		for {
			select {
			case <-pointer:
				continue
			default:
			}
			break
		}
		if err := c.moveRaw(win, x, y); err != nil {
			return pointerReport{}, err
		}
		// The last report once they stop: a report of earlier motion can
		// still be on its way when the pointer moves.
		var last *pointerReport
		deadline := time.After(3 * time.Second)
		for {
			settled := time.After(250 * time.Millisecond)
			select {
			case r := <-pointer:
				if r.window == page {
					last = &r
				}
			case <-settled:
				if last != nil {
					return *last, nil
				}
			case <-deadline:
				if last != nil {
					return *last, nil
				}
				return pointerReport{}, fmt.Errorf("the page saw no pointer motion to (%d, %d) of %q", x, y, win)
			}
		}
	}
	// The points are far enough apart for a precise scale, and inside the
	// page of an 800x500 window at up to twice the scale.
	const x1, y1, x2, y2 = 120, 140, 520, 340
	a, err := at(x1, y1)
	if err != nil {
		return err
	}
	b, err := at(x2, y2)
	if err != nil {
		return err
	}
	m, err := solveMapping(x1, y1, x2, y2, a, b)
	if err != nil {
		return err
	}
	if v, ok := dpr.Load(page); ok {
		if r := v.(float64); r > 0 && (m.scale < r*0.9 || m.scale > r*1.1) {
			return fmt.Errorf("pointer scale %.2f but the page's devicePixelRatio is %.2f", m.scale, r)
		}
	}
	c.setMapping(win, m)
	fmt.Printf("e2e: %q page origin (%.0f, %.0f), scale %.2f\n", win, m.ox, m.oy, m.scale)
	return nil
}

// solveMapping finds raw = origin + scale*page from two raw points and the
// page points they landed on.
func solveMapping(x1, y1, x2, y2 int, a, b pointerReport) (mapping, error) {
	dx, dy := b.x-a.x, b.y-a.y
	if dx <= 0 || dy <= 0 {
		return mapping{}, fmt.Errorf("the page saw the pointer at (%v, %v) then (%v, %v): not a mapping", a.x, a.y, b.x, b.y)
	}
	sx, sy := float64(x2-x1)/dx, float64(y2-y1)/dy
	if sx/sy < 0.95 || sx/sy > 1.05 {
		return mapping{}, fmt.Errorf("pointer scale differs by axis: %.3f by %.3f", sx, sy)
	}
	s := (sx + sy) / 2
	return mapping{ox: float64(x1) - s*a.x, oy: float64(y1) - s*a.y, scale: s}, nil
}

func mustEval(path string) string {
	if p, err := filepath.EvalSymlinks(path); err == nil {
		return p
	}
	return path
}

// shot prints a full-screen screenshot as evidence, and writes it to
// E2E_EVIDENCE_DIR when that is set.
func shot(name string) {
	file := fmt.Sprintf("%s-%s-e2e-%s.png", runtime.GOOS, runtime.GOARCH, name)
	path := filepath.Join(os.TempDir(), file)
	if err := screenshot(path); err != nil {
		fmt.Fprintf(os.Stderr, "screenshot %s: %v\n", name, err)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	if dir := os.Getenv("E2E_EVIDENCE_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, file), data, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "evidence %s: %v\n", name, err)
		}
	}
	fmt.Printf("EVIDENCE-BEGIN %s\n", file)
	enc := base64.StdEncoding.EncodeToString(data)
	for len(enc) > 0 {
		n := min(len(enc), 76)
		fmt.Println(enc[:n])
		enc = enc[n:]
	}
	fmt.Println("EVIDENCE-END")
}

func list(env string) map[string]bool {
	out := map[string]bool{}
	for _, s := range strings.Split(os.Getenv(env), ",") {
		if s = strings.TrimSpace(s); s != "" {
			out[s] = true
		}
	}
	return out
}

func selected(goos string) []scenario {
	only, skip := list("E2E_SCENARIOS"), list("E2E_SKIP")
	var out []scenario
	for _, s := range scenarios {
		if s.runsOn(goos) && (len(only) == 0 || only[s.name]) && !skip[s.name] {
			out = append(out, s)
		}
	}
	return out
}

// assets serves the page, and turns its reports into signals.
func assets(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	win := q.Get("w")
	switch r.URL.Path {
	case "/input":
		emit("input " + win + " " + q.Get("id") + "=" + q.Get("v"))
	case "/click":
		emit("click " + win)
	case "/event":
		emit("event " + win + " " + q.Get("v"))
	case "/rt":
		emit("rt " + q.Get("step") + " " + q.Get("r") + " " + q.Get("v"))
	case "/runtime":
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, runtimePage)
	case "/ready":
		if v, err := strconv.ParseFloat(q.Get("dpr"), 64); err == nil {
			dpr.Store(win, v)
		}
		emit("ready " + win)
	case "/mouse":
		x, errX := strconv.ParseFloat(q.Get("x"), 64)
		y, errY := strconv.ParseFloat(q.Get("y"), 64)
		if errX == nil && errY == nil {
			select {
			case pointer <- pointerReport{win, x, y}:
			default:
			}
		}
	default:
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, page)
	}
}

func main() {
	if os.Getenv(secondInstanceEnv) != "" {
		// The copy secondInstance starts: New hands over and exits.
		application.New(application.Options{Name: "E2E probe", SingleInstance: singleInstance()})
		fmt.Println("E2E second instance: New returned; the first copy was not found")
		exit(4)
	}
	d, err := newDriver()
	if err != nil {
		fmt.Fprintln(os.Stderr, "E2E driver:", err)
		exit(2)
	}

	p := &probe{d: d}
	p.menu = application.NewMenu()
	actions := p.menu.AddSubmenu("E2E menu")
	actions.Add("Action").SetAccelerator("CmdOrCtrl+M").OnClick(func(*application.Context) { emit("menu-action") })
	actions.Add("Other")
	actions.AddSeparator()
	p.toggle = actions.AddCheckbox("Toggle", false).SetAccelerator("CmdOrCtrl+Shift+B").
		OnClick(func(ctx *application.Context) { emit(fmt.Sprintf("toggle %v", ctx.ClickedMenuItem().Checked())) })
	p.alpha = actions.AddRadio("Alpha", true).SetAccelerator("CmdOrCtrl+Shift+Y").
		OnClick(func(*application.Context) { emit("radio Alpha") })
	p.beta = actions.AddRadio("Beta", false).SetAccelerator("CmdOrCtrl+Shift+G").
		OnClick(func(*application.Context) { emit("radio Beta") })
	actions.AddSeparator()
	p.disabled = actions.Add("Disabled").SetAccelerator("CmdOrCtrl+Shift+D").SetEnabled(false).
		OnClick(func(*application.Context) { emit("disabled-action") })
	p.rebound = actions.Add("Rebound").SetAccelerator("CmdOrCtrl+Shift+E").
		OnClick(func(*application.Context) { emit("rebound-action") })
	p.hidden = actions.Add("Hidden").SetAccelerator("CmdOrCtrl+Shift+L").
		OnClick(func(*application.Context) { emit("hidden-action") })
	p.menu.AddRole(application.EditMenu)
	// The View role: Reload, Zoom In/Out, Actual Size, full screen. (Not the
	// Window role: its Minimise takes CmdOrCtrl+M, the Action item's.)
	p.menu.AddRole(application.ViewMenu)

	p.app = application.New(application.Options{
		Name:   "E2E probe",
		Assets: application.AssetOptions{Handler: http.HandlerFunc(assets)},
		KeyBindings: map[string]func(window application.Window){
			"CmdOrCtrl+Shift+K": func(w application.Window) { emit("keybinding " + w.Name()) },
		},
		Mac:            application.MacOptions{ActivationPolicy: application.ActivationPolicyRegular},
		Services:       services,
		SingleInstance: singleInstance(),
	})
	for _, setup := range setups {
		setup(p)
	}
	if runtime.GOOS == "darwin" {
		p.app.Menu.Set(p.menu)
	}
	ctx := p.app.ContextMenu.New()
	ctx.Add("Context action").OnClick(func(*application.Context) { emit("context-action") })
	ctx.Add("Other")
	p.app.ContextMenu.Add("e2e", ctx)

	p.main = p.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name: "main", Title: title, Width: 800, Height: 500, X: 100, Y: 100,
		Linux:   application.LinuxWindow{Menu: p.menu},
		Windows: application.WindowsWindow{Menu: p.menu},
	})

	failed := 0
	p.app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		go func() {
			if err := await("ready main", 30*time.Second); err != nil {
				fmt.Println("E2E startup fail:", err)
				exit(1)
			}
			for _, lock := range strings.Split(os.Getenv("E2E_LOCKS"), ",") {
				if lock = strings.TrimSpace(lock); lock != "" {
					if err := press(d, title, lock); err != nil {
						fmt.Println("E2E locks fail:", err)
						exit(1)
					}
				}
			}
			todo := selected(runtime.GOOS)
			if len(todo) == 0 {
				// A misspelt E2E_SCENARIOS must not pass having run nothing.
				fmt.Println("E2E startup fail: no scenario selected for", runtime.GOOS)
				exit(1)
			}
			for _, s := range todo {
				drain()
				if err := s.run(p); err != nil {
					failed++
					fmt.Printf("E2E %s fail: %v\n", s.name, err)
				} else {
					fmt.Printf("E2E %s ok\n", s.name)
				}
				// Kept: lets the previous scenario's windows and menus go before
				// the next one starts. Nothing observable marks that point.
				time.Sleep(500 * time.Millisecond)
			}
			if failed > 0 {
				exit(1)
			}
			exit(0)
		}()
	})
	if err := p.app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exit(1)
	}
}
