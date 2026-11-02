// The e2e probe drives a real Wails app with real OS input: key bindings,
// menu accelerators, Edit-menu shortcuts in the page, keyboard menu
// navigation, a context menu, message
// dialogs and the open-file dialog. It runs on the lab VMs (never on a
// developer's console session):
//
//	macOS:  go run -tags purego ./tests/purego/e2e
//	Linux:  xvfb-run -a -s "-screen 0 1280x800x24" dbus-run-session -- go run ./tests/purego/e2e
//
// Each scenario prints "E2E <name> ok" or "E2E <name> fail: <reason>", and a
// screenshot of any open menu or dialog between "EVIDENCE-BEGIN <name>.png"
// and "EVIDENCE-END" (base64). Build it with -cover and set GOCOVERDIR to
// count what it exercises. E2E_SCENARIOS (comma separated) limits the run.
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const title = "E2E"

// The page: a target for the context menu at a known place.
const page = `<!doctype html><html><head><script type="module" src="/wails/runtime.js"></script></head>
<body style="margin:0;font:16px system-ui,-apple-system,sans-serif;background:#1e293b;color:white">
<div id="ctx" style="position:absolute;left:20px;top:20px;width:240px;height:100px;background:#7c3aed;
display:flex;align-items:center;justify-content:center;--custom-contextmenu:e2e">Right-click me</div>
<p style="position:absolute;left:20px;top:140px">E2E probe</p>
<input id="src" value="wails-e2e-clip" style="position:absolute;left:300px;top:20px;width:200px;height:30px">
<input id="dst" style="position:absolute;left:300px;top:70px;width:200px;height:30px">
<script>document.addEventListener('input', (e) =>
 fetch('/input?id=' + e.target.id + '&v=' + encodeURIComponent(e.target.value)));</script></body></html>`

// The inputs' centres, in the page's coordinates.
const srcX, srcY, dstX, dstY = 400, 35, 400, 85

var signals = make(chan string, 32)

// eventPause paces the input of a multi-step scenario: the page has no signal
// for a click having moved focus or a selection having been made.
const eventPause = 300 * time.Millisecond

func signal(name string) func() { return func() { signals <- name } }

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

type scenario struct {
	name string
	os   string // "" = every OS
	run  func(app *application.App, d driver) error
}

var scenarios = []scenario{
	{name: "keybinding", run: func(app *application.App, d driver) error {
		if err := d.focus(title); err != nil {
			return err
		}
		if err := d.key("cmdorctrl+shift+k"); err != nil {
			return err
		}
		return await("keybinding", 5*time.Second)
	}},
	{name: "menu-accelerator", run: func(app *application.App, d driver) error {
		if err := d.focus(title); err != nil {
			return err
		}
		if err := d.key("cmdorctrl+m"); err != nil {
			return err
		}
		return await("menu-action", 5*time.Second)
	}},
	{name: "menu-navigate", os: "linux", run: func(app *application.App, d driver) error {
		// F10 opens the menu bar; the first item of the first menu is
		// "Action". This is the user's path through the native menu.
		if err := d.focus(title); err != nil {
			return err
		}
		if err := d.key("f10"); err != nil {
			return err
		}
		// Kept: lets the menu bar paint before the screenshot. The driver has no
		// way to observe the open menu, and the keys below need it open.
		time.Sleep(time.Second)
		shot("menu-navigate")
		keys := []string{"down", "return"}
		if menuOpensOnFirstItem {
			keys = []string{"return"}
		}
		for _, k := range keys {
			if err := d.key(k); err != nil {
				return err
			}
		}
		return await("menu-action", 5*time.Second)
	}},
	{name: "context-menu", run: func(app *application.App, d driver) error {
		if err := d.rightClick(title, 100, 60); err != nil {
			return err
		}
		// Kept: lets the popover paint before the screenshot. The driver cannot
		// observe it, and the keys below need it open.
		time.Sleep(time.Second)
		shot("context-menu")
		if runtime.GOOS == "linux" && contextMenuByMouse {
			// The popover opens at the pointer; its first item is just below
			// and to the right of it.
			if err := d.click(title, 120, 83); err != nil {
				return err
			}
			return await("context-action", 5*time.Second)
		}
		for _, k := range []string{"down", "return"} {
			if err := d.key(k); err != nil {
				return err
			}
		}
		return await("context-action", 5*time.Second)
	}},
	{name: "dialog-default-button", run: func(app *application.App, d driver) error {
		dialog := app.Dialog.Info().SetTitle("E2E dialog").SetMessage("Return presses the default button")
		dialog.AddButton("Later").OnClick(signal("dialog-later"))
		dialog.AddButton("OK").SetAsDefault().OnClick(signal("dialog-ok"))
		application.InvokeAsync(dialog.Show)
		// Kept: the driver cannot detect this dialog (GTK4 message dialogs are
		// untitled windows, and darwin focus reports success regardless), so
		// the probe waits for it to paint before the key is sent.
		time.Sleep(1500 * time.Millisecond)
		shot("dialog-default-button")
		if err := d.key("return"); err != nil {
			return err
		}
		return await("dialog-ok", 5*time.Second)
	}},
	{name: "dialog-cancel-button", run: func(app *application.App, d driver) error {
		dialog := app.Dialog.Question().SetTitle("E2E question").SetMessage("Escape presses the cancel button")
		dialog.AddButton("Yes").SetAsDefault().OnClick(signal("question-yes"))
		dialog.AddButton("No").SetAsCancel().OnClick(signal("question-no"))
		application.InvokeAsync(dialog.Show)
		// Kept: as in dialog-default-button, the driver cannot detect this
		// dialog, so the probe waits for it to paint before the key is sent.
		time.Sleep(1500 * time.Millisecond)
		shot("dialog-cancel-button")
		if err := d.key("escape"); err != nil {
			return err
		}
		return await("question-no", 5*time.Second)
	}},
	{name: "edit-copy-paste", run: func(app *application.App, d driver) error {
		// The Edit menu binds Ctrl+A/C/V/X/Z (Cmd on macOS). Those shortcuts must
		// still edit the page: GTK3 used to swallow them in empty role bodies.
		if err := d.focus(title); err != nil {
			return err
		}
		for _, step := range []func() error{
			func() error { return d.click(title, srcX, srcY) },
			func() error { return d.key("cmdorctrl+a") },
			func() error { return d.key("cmdorctrl+c") },
			func() error { return d.click(title, dstX, dstY) },
			func() error { return d.key("cmdorctrl+v") },
		} {
			if err := step(); err != nil {
				return err
			}
			time.Sleep(eventPause)
		}
		if err := await("input dst=wails-e2e-clip", 5*time.Second); err != nil {
			return fmt.Errorf("paste: %w", err)
		}
		for _, step := range []func() error{
			func() error { return d.click(title, srcX, srcY) },
			func() error { return d.key("cmdorctrl+a") },
			func() error { return d.key("cmdorctrl+x") },
		} {
			if err := step(); err != nil {
				return err
			}
			time.Sleep(eventPause)
		}
		if err := await("input src=", 5*time.Second); err != nil {
			return fmt.Errorf("cut: %w", err)
		}
		if text, ok := app.Clipboard.Text(); !ok || text != "wails-e2e-clip" {
			return fmt.Errorf("clipboard after cut = %q, %v", text, ok)
		}
		if err := d.key("cmdorctrl+z"); err != nil {
			return err
		}
		if err := await("input src=wails-e2e-clip", 5*time.Second); err != nil {
			return fmt.Errorf("undo: %w", err)
		}
		return nil
	}},
	{name: "open-file", run: func(app *application.App, d driver) error {
		dir, err := os.MkdirTemp("", "wails-e2e")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		want := filepath.Join(dir, "chosen.txt")
		if err := os.WriteFile(want, []byte("e2e"), 0o644); err != nil {
			return err
		}
		type result struct {
			path string
			err  error
		}
		done := make(chan result, 1)
		go func() {
			path, err := app.Dialog.OpenFile().SetTitle("E2E open").SetDirectory(dir).CanChooseFiles(true).PromptForSingleSelection()
			done <- result{path, err}
		}()
		// Kept: lets the open panel paint before the screenshot. The Linux path
		// below polls for the dialog; the macOS path has no readable signal.
		time.Sleep(2 * time.Second)
		shot("open-file")
		if c, ok := d.(interface{ cancelDialog() error }); ok {
			// macOS: the panel's contents run in an Apple service process
			// that posted events cannot reach, so the probe checks the panel
			// opens (screenshot) and that cancelling it returns cleanly.
			if err := c.cancelDialog(); err != nil {
				return err
			}
			select {
			case r := <-done:
				if r.err != nil || r.path != "" {
					return fmt.Errorf("cancelled dialog returned %q, %v; want \"\", nil", r.path, r.err)
				}
				return nil
			case <-time.After(10 * time.Second):
				return fmt.Errorf("the cancelled open-file dialog did not return")
			}
		}
		// The dialog is its own top-level window; with no window manager it
		// does not get the keyboard by itself. Wait until it exists.
		for deadline := time.Now().Add(5 * time.Second); d.focus("E2E open") != nil && time.Now().Before(deadline); {
			time.Sleep(100 * time.Millisecond)
		}
		// Type the path into the dialog's location field (Ctrl+L in GTK). The
		// pauses in these steps are kept: they pace the input, and GTK has no
		// signal the driver can read for when its pre-fill or completion is done.
		steps := []func() error{
			func() error { return d.key("ctrl+l") },
			// GTK4 pre-fills the field with the current folder: replace it.
			func() error { time.Sleep(500 * time.Millisecond); return d.key("ctrl+a") },
			func() error { return d.typeText(want) },
			// GTK completes the path inline while it is typed, inserting the
			// rest of the name after the cursor. Select from the cursor to the
			// end and delete it, so Return takes exactly what was typed.
			func() error { time.Sleep(300 * time.Millisecond); return d.key("shift+end") },
			func() error { return d.key("delete") },
			func() error { return d.key("return") },
		}
		for _, step := range steps {
			if err := step(); err != nil {
				return err
			}
			// Kept: pacing between input events (see the steps above).
			time.Sleep(300 * time.Millisecond)
		}
		select {
		case r := <-done:
			if r.err != nil {
				return r.err
			}
			if got, _ := filepath.EvalSymlinks(r.path); got != mustEval(want) {
				return fmt.Errorf("dialog returned %q, want %q", r.path, want)
			}
			return nil
		case <-time.After(10 * time.Second):
			return fmt.Errorf("the open-file dialog did not return")
		}
	}},
}

func mustEval(path string) string {
	if p, err := filepath.EvalSymlinks(path); err == nil {
		return p
	}
	return path
}

// shot prints a full-screen screenshot as evidence.
func shot(name string) {
	file := filepath.Join(os.TempDir(), "e2e-"+name+".png")
	cmd := exec.Command("import", "-window", "root", file)
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("screencapture", "-x", file)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "screenshot %s: %v: %s\n", name, err, bytes.TrimSpace(out))
		return
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return
	}
	fmt.Printf("EVIDENCE-BEGIN %s-%s-e2e-%s.png\n", runtime.GOOS, runtime.GOARCH, name)
	enc := base64.StdEncoding.EncodeToString(data)
	for len(enc) > 0 {
		n := min(len(enc), 76)
		fmt.Println(enc[:n])
		enc = enc[n:]
	}
	fmt.Println("EVIDENCE-END")
}

func selected() []scenario {
	only := map[string]bool{}
	for _, s := range strings.Split(os.Getenv("E2E_SCENARIOS"), ",") {
		if s = strings.TrimSpace(s); s != "" {
			only[s] = true
		}
	}
	var out []scenario
	for _, s := range scenarios {
		if (s.os == "" || s.os == runtime.GOOS) && (len(only) == 0 || only[s.name]) {
			out = append(out, s)
		}
	}
	return out
}

func main() {
	d, err := newDriver()
	if err != nil {
		fmt.Fprintln(os.Stderr, "E2E driver:", err)
		exit(2)
	}

	menu := application.NewMenu()
	actions := menu.AddSubmenu("E2E menu")
	actions.Add("Action").SetAccelerator("CmdOrCtrl+M").OnClick(func(*application.Context) { signals <- "menu-action" })
	actions.Add("Other")
	menu.AddRole(application.EditMenu)

	app := application.New(application.Options{
		Name: "E2E probe",
		Assets: application.AssetOptions{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/input" {
				signals <- "input " + r.URL.Query().Get("id") + "=" + r.URL.Query().Get("v")
				return
			}
			w.Header().Set("Content-Type", "text/html")
			io.WriteString(w, page)
		})},
		KeyBindings: map[string]func(window application.Window){
			"CmdOrCtrl+Shift+K": func(application.Window) { signals <- "keybinding" },
		},
		Mac: application.MacOptions{ActivationPolicy: application.ActivationPolicyRegular},
	})
	if runtime.GOOS == "darwin" {
		app.Menu.Set(menu)
	}
	ctx := app.ContextMenu.New()
	ctx.Add("Context action").OnClick(func(*application.Context) { signals <- "context-action" })
	ctx.Add("Other")
	app.ContextMenu.Add("e2e", ctx)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: title, Width: 800, Height: 500, X: 100, Y: 100,
		Linux: application.LinuxWindow{Menu: menu},
	})

	failed := 0
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		go func() {
			// Kept: the probe has no JS evaluation to poll the page or runtime
			// with, so it gives the page and its runtime time to load.
			time.Sleep(2 * time.Second) // the page and its runtime load
			for _, s := range selected() {
				if err := s.run(app, d); err != nil {
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
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exit(1)
	}
}
