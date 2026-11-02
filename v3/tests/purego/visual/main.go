// The visual evidence probe opens each window kind the backend supports, one
// at a time, and screenshots the whole screen while it is up. It is run on the
// lab VMs (never on a developer's console session):
//
//	macOS:  go run -tags purego ./tests/purego/visual
//	Linux:  xvfb-run -a -s "-screen 0 1280x800x24" dbus-run-session -- go run ./tests/purego/visual
//	Windows: go run ./tests/purego/visual  (in an interactive desktop session)
//
// Each screenshot is printed to stdout as base64 between
// "EVIDENCE-BEGIN <name>.png" and "EVIDENCE-END" lines, because the lab only
// returns the job's output. The harness decodes them into PNG files.
// VISUAL_SCENARIOS (comma separated) limits the run to some scenarios; the
// scenarios marked alone run only when it names them.
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// settle is how long a scenario stays on screen before the screenshot, so the
// webview has painted.
const settle = 3 * time.Second

func page(title, colour string) string {
	return `<!doctype html><html><body style="margin:0;height:100vh;display:flex;align-items:center;justify-content:center;` +
		`background:` + colour + `;color:white;font:600 28px -apple-system,system-ui,sans-serif">` + title + `</body></html>`
}

type scenario struct {
	name string
	os   string // "" = every OS, else a comma-separated list
	// alone scenarios run only when VISUAL_SCENARIOS names them.
	alone bool
	// prepare, if set, runs before the application starts.
	prepare func(app *application.App)
	run     func(app *application.App) (cleanup func())
}

// preStarted is the window menu-nil-before-start makes before the app runs.
var preStarted *application.WebviewWindow

var scenarios = []scenario{
	{name: "window", run: func(app *application.App) func() {
		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "Standard window (" + backend() + ")", Width: 640, Height: 400, X: 80, Y: 80,
			HTML: page("Standard window - "+backend(), "#2563eb"),
		})
		return w.Close
	}},
	{name: "window-frameless", run: func(app *application.App) func() {
		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "Frameless", Width: 480, Height: 300, X: 120, Y: 120, Frameless: true,
			HTML: page("Frameless window", "#059669"),
		})
		return w.Close
	}},
	{name: "square-corners", os: "darwin", run: func(app *application.App) func() {
		// cgo honours CornerType only for frameless windows (windowStyleMask),
		// so the window is frameless: a borderless window with square corners.
		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "Square corners", Width: 480, Height: 300, X: 120, Y: 120, Frameless: true,
			HTML: page("Square corners (borderless)", "#0f766e"),
			Mac:  application.MacWindow{CornerType: application.MacWindowCornerTypeSquare},
		})
		return w.Close
	}},
	{name: "corner-radius", os: "darwin", run: func(app *application.App) func() {
		// A frameless window with a custom radius: borderless, with the content
		// view's layer rounded to 24 points.
		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "Corner radius", Width: 480, Height: 300, X: 120, Y: 120, Frameless: true,
			HTML: page("Corner radius 24", "#be185d"),
			Mac:  application.MacWindow{CornerRadius: 24},
		})
		return w.Close
	}},
	{name: "menu-label", os: "linux", run: func(app *application.App) func() {
		// SetLabel after the menu is attached must show in the menu bar
		// (pg-12m: GMenu copies items on append, so the label used to stay
		// "Before"). The radio read-back is logged for the record.
		menu := app.NewMenu()
		menu.AddSubmenu("Before").Add("Item")
		radios := menu.AddSubmenu("Radios")
		first := radios.AddRadio("First", true)
		second := radios.AddRadio("Second", false)
		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "Menu label", Width: 640, Height: 400,
			HTML: page("Menu bar should read: After | Radios", "#b45309"),
			// The menu bar is built when the window starts.
			Linux: application.LinuxWindow{Menu: menu},
		})
		go func() {
			waitFor(2*time.Second, w.IsVisible) // the window is up, so its menu bar exists
			application.InvokeSync(func() {
				menu.FindByLabel("Before").SetLabel("After")
				second.SetChecked(true)
			})
			fmt.Printf("RADIO first=%v second=%v (want false true)\n", first.Checked(), second.Checked())
		}()
		return w.Close
	}},
	{name: "menu-setmenu-running", os: "linux", run: func(app *application.App) func() {
		// A window with no menu gets one after it is running: SetMenu on a
		// realised window must show the menu bar at once (pg-88b).
		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "SetMenu at runtime", Width: 640, Height: 400,
			HTML: page("Menu bar should read: Set at runtime", "#0f766e"),
		})
		go func() {
			waitFor(2*time.Second, w.IsVisible) // SetMenu must run on a realised window
			application.InvokeSync(func() {
				menu := app.NewMenu()
				menu.AddSubmenu("Set at runtime").Add("Item")
				w.SetMenu(menu)
			})
		}()
		return w.Close
	}},
	{name: "menu-setmenu-nil", os: "linux,windows", run: func(app *application.App) func() {
		// A running window's menu bar goes when SetMenu(nil) is called (pg-ifz).
		menu := app.NewMenu()
		menu.AddSubmenu("Should vanish").Add("Item")
		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "SetMenu(nil) at runtime", Width: 640, Height: 400, X: 80, Y: 80,
			HTML:    page("No menu bar above this", "#7c2d12"),
			Linux:   application.LinuxWindow{Menu: menu},
			Windows: application.WindowsWindow{Menu: menu},
		})
		go func() {
			waitFor(2*time.Second, w.IsVisible) // SetMenu(nil) must run on a realised window
			application.InvokeSync(func() { w.SetMenu(nil) })
		}()
		return w.Close
	}},
	{name: "menu-nil-before-start", os: "linux,windows", alone: true,
		prepare: func(app *application.App) {
			// The application menu is set and the window opts into it, but
			// SetMenu(nil) is called before the window starts: it must be
			// created with no menu bar at all (pg-ifz). The app menu would
			// show in every other scenario's window, so this one runs alone.
			menu := app.NewMenu()
			menu.AddSubmenu("Application menu").Add("Item")
			app.Menu.Set(menu)
			preStarted = app.Window.NewWithOptions(application.WebviewWindowOptions{
				Title: "SetMenu(nil) before start", Width: 640, Height: 400, X: 80, Y: 80,
				HTML:               page("No menu bar above this either", "#1e3a8a"),
				UseApplicationMenu: true,
				Hidden:             true,
			})
			preStarted.SetMenu(nil)
		},
		run: func(app *application.App) func() {
			preStarted.Show()
			return preStarted.Close
		}},
	{name: "panel", os: "darwin", run: func(app *application.App) func() {
		normal := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "Behind the panel", Width: 640, Height: 400, X: 60, Y: 60,
			HTML: page("Ordinary window", "#475569"),
		})
		panel := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title: "Floating panel", Width: 360, Height: 220, X: 400, Y: 300,
			HTML: page("NSPanel (floating, non-activating)", "#9333ea"),
			Mac: application.MacWindow{
				WindowClass:      application.MacWindowClassPanel,
				PanelPreferences: application.MacPanelPreferences{FloatingPanel: true, BecomesKeyOnlyIfNeeded: true},
			},
		})
		return func() { panel.Close(); normal.Close() }
	}},
	{name: "notch", os: "darwin", run: func(app *application.App) func() {
		n := app.Window.NewNotchWindow(application.NotchWindowOptions{
			WindowOptions: application.WebviewWindowOptions{HTML: page("Notch window", "#000000")},
		})
		n.Show()
		return n.Close
	}},
	{name: "notch-animated", os: "darwin", run: func(app *application.App) func() {
		// Shown with the slide-in animation; the screenshot is taken after it
		// completes (420 ms, well inside the settle time).
		n := app.Window.NewNotchWindow(application.NotchWindowOptions{
			Animated:      true,
			WindowOptions: application.WebviewWindowOptions{HTML: page("Notch animated", "#000000")},
		})
		n.Show()
		return n.Close
	}},
	{name: "notch-hidden", os: "darwin", run: func(app *application.App) func() {
		// Shown, then hidden with the slide-out animation about a second later.
		// The screenshot is taken after the hide completes (280 ms), so the
		// notch must be off screen. Hide is not called synchronously here: it
		// would pre-empt the show animation.
		n := app.Window.NewNotchWindow(application.NotchWindowOptions{
			Animated:      true,
			WindowOptions: application.WebviewWindowOptions{HTML: page("Notch hidden", "#000000")},
		})
		n.Show()
		go func() {
			// Kept: the slide-in (420 ms) must finish before Hide, or Hide pre-empts
			// it. The probe has no signal for the end of the animation.
			time.Sleep(time.Second)
			application.InvokeSync(func() { n.Hide() })
			// Wait for the slide-out (280 ms) to take the notch off screen, up to
			// 2 s. On timeout the line below still prints the real state.
			waitFor(2*time.Second, func() bool { return !n.Visibility() })
			application.InvokeSync(func() {
				fmt.Printf("NOTCH visible=%v (want false after hide)\n", n.Visibility())
			})
		}()
		return n.Close
	}},
	{name: "tray-menu", os: "darwin", run: func(app *application.App) func() {
		menu := app.NewMenu()
		menu.Add("First item")
		menu.AddCheckbox("Checked item", true)
		menu.AddSeparator()
		menu.Add("Quit")
		tray := app.SystemTray.New()
		tray.SetLabel("Wails purego")
		tray.SetMenu(menu)
		// Opening the menu starts AppKit's menu-tracking loop on the main
		// thread; the screenshot is taken from this goroutine while it is open.
		// Kept: lets the status item be placed before it is opened. The tray has
		// no readiness signal to poll.
		go func() { time.Sleep(time.Second); tray.OpenMenu() }()
		return func() {}
	}},
}

func backend() string {
	if runtime.GOOS == "darwin" {
		return "purego"
	}
	return "purego GTK"
}

func screenshot(name string) error {
	file := filepath.Join(os.TempDir(), name+".png")
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("screencapture", "-x", file)
	case "windows":
		// The whole virtual screen, through GDI+ (needs an interactive session).
		cmd = exec.Command("powershell", "-NoProfile", "-Command", `Add-Type -AssemblyName System.Windows.Forms,System.Drawing
$b = [System.Windows.Forms.SystemInformation]::VirtualScreen
$i = New-Object System.Drawing.Bitmap $b.Width, $b.Height
$g = [System.Drawing.Graphics]::FromImage($i)
$g.CopyFromScreen($b.Location, [System.Drawing.Point]::Empty, $b.Size)
$i.Save('`+file+`', [System.Drawing.Imaging.ImageFormat]::Png)`)
	default:
		cmd = exec.Command("import", "-window", "root", file)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %v: %s", cmd.Path, err, bytes.TrimSpace(out))
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	fmt.Printf("EVIDENCE-BEGIN %s.png\n", name)
	enc := base64.StdEncoding.EncodeToString(data)
	for len(enc) > 0 {
		n := min(len(enc), 76)
		fmt.Println(enc[:n])
		enc = enc[n:]
	}
	fmt.Println("EVIDENCE-END")
	return nil
}

// waitFor polls cond every 20 ms until it holds or timeout passes, and reports
// whether it held.
func waitFor(timeout time.Duration, cond func() bool) bool {
	deadline := time.Now().Add(timeout)
	for {
		if cond() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func selected() []scenario {
	only := map[string]bool{}
	for _, s := range strings.Split(os.Getenv("VISUAL_SCENARIOS"), ",") {
		if s = strings.TrimSpace(s); s != "" {
			only[s] = true
		}
	}
	var out []scenario
	for _, s := range scenarios {
		if (s.os == "" || slices.Contains(strings.Split(s.os, ","), runtime.GOOS)) && (only[s.name] || len(only) == 0 && !s.alone) {
			out = append(out, s)
		}
	}
	return out
}

func main() {
	app := application.New(application.Options{
		Name: "Visual evidence probe",
		Mac:  application.MacOptions{ActivationPolicy: application.ActivationPolicyRegular},
		// Each scenario closes its windows before the next opens.
		Linux:   application.LinuxOptions{DisableQuitOnLastWindowClosed: true},
		Windows: application.WindowsOptions{DisableQuitOnLastWindowClosed: true},
	})
	run := selected()
	for _, s := range run {
		if s.prepare != nil {
			s.prepare(app)
		}
	}
	failed := false
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		go func() {
			for _, s := range run {
				var cleanup func()
				application.InvokeSync(func() { cleanup = s.run(app) })
				time.Sleep(settle)
				prefix := runtime.GOOS + "-" + runtime.GOARCH + "-" + s.name
				if err := screenshot(prefix); err != nil {
					fmt.Fprintf(os.Stderr, "SCENARIO %s screenshot failed: %v\n", s.name, err)
					failed = true
				} else {
					fmt.Printf("SCENARIO %s ok\n", s.name)
				}
				if s.name == "tray-menu" {
					// The open menu holds the main thread in its tracking loop;
					// it is the last scenario, so finish here.
					break
				}
				application.InvokeSync(cleanup)
				// Kept: lets the closed window leave the screen before the next
				// scenario. Nothing observable marks that point.
				time.Sleep(500 * time.Millisecond)
			}
			if failed {
				os.Exit(1)
			}
			os.Exit(0)
		}()
	})
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
