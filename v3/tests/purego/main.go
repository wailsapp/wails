// The native parity probe runs unchanged with either the cgo or purego backend.
// Run under a desktop session or dbus-run-session -- xvfb-run -a ./probe.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Probe struct{}

func (*Probe) Echo(value string) string { return value }

func main() {
	var app *application.App
	var window *application.WebviewWindow
	var completed atomic.Bool
	var once sync.Once
	cancelled := make(chan struct{})
	app = application.New(application.Options{
		Name:     "Native parity probe",
		Services: []application.Service{application.NewService(&Probe{})},
		Assets: application.AssetOptions{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/form":
				// FormData: the server must see a real multipart form.
				if err := r.ParseMultipartForm(1 << 20); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				file, _, err := r.FormFile("f")
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				defer file.Close()
				io.Copy(w, file)
			case "/echo":
				w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
				io.Copy(w, r.Body)
			case "/stream":
				w.Header().Set("Content-Type", "text/plain")
				io.WriteString(w, "ready\n")
				w.(http.Flusher).Flush()
				select {
				case <-r.Context().Done():
					once.Do(func() { close(cancelled) })
				case <-time.After(15 * time.Second):
				}
			case "/done":
				if message := r.URL.Query().Get("error"); message != "" {
					fmt.Fprintln(os.Stderr, message)
					exit(1)
				}
				select {
				case <-cancelled:
				case <-time.After(5 * time.Second):
					fmt.Fprintln(os.Stderr, "request cancellation timed out")
					exit(1)
				}
				app.Clipboard.SetText("purego clipboard")
				if value, ok := app.Clipboard.Text(); !ok || value != "purego clipboard" {
					fmt.Fprintln(os.Stderr, "clipboard round trip failed")
					exit(1)
				}
				window.SetZoom(1.25)
				if window.GetZoom() != 1.25 {
					fmt.Fprintln(os.Stderr, "zoom round trip failed")
					exit(1)
				}
				if err := checkNativeAPIs(app, window); err != nil {
					fmt.Fprintln(os.Stderr, err)
					exit(1)
				}
				completed.Store(true)
				fmt.Println("PASS: assets, binary POST, JavaScript bridge, request cancellation, clipboard, zoom, menus, screens, windows")
				flushCoverage()
				app.Quit()
			default:
				w.Header().Set("Content-Type", "text/html")
				io.WriteString(w, `<!doctype html><html><body><h1>Native parity probe</h1>
<script type="module">
import {Call} from '/wails/runtime.js';
try {
 const result = await Call.ByName('main.Probe.Echo', 'hello native');
 if (result !== 'hello native') throw new Error('bridge round trip failed: ' + result);
 const bytes = new Uint8Array([0, 1, 128, 255, 10]);
 const echo = async (name, send) => {
   const result = new Uint8Array(await send());
   if (String(result) !== String(bytes)) throw new Error('binary POST mismatch (' + name + '): ' + result);
 };
 const post = (body) => fetch('/echo', {method: 'POST', body}).then((r) => r.arrayBuffer());
 await echo('Uint8Array', () => post(bytes));
 await echo('ArrayBuffer', () => post(bytes.slice().buffer));
 await echo('DataView', () => post(new DataView(bytes.slice().buffer)));
 await echo('Blob', () => post(new Blob([bytes], {type: 'application/octet-stream'})));
 await echo('Request', () => fetch(new Request('/echo', {method: 'POST', body: bytes})).then((r) => r.arrayBuffer()));
 await echo('FormData', () => {
   const form = new FormData();
   form.append('f', new Blob([bytes]), 'f.bin');
   return fetch('/form', {method: 'POST', body: form}).then((r) => r.arrayBuffer());
 });
 await echo('XMLHttpRequest (sync)', () => {
   const xhr = new XMLHttpRequest();
   xhr.open('POST', '/echo', false);
   xhr.overrideMimeType('text/plain; charset=x-user-defined');
   xhr.send(bytes);
   return Uint8Array.from(xhr.responseText, (c) => c.charCodeAt(0) & 0xff).buffer;
 });
 await echo('XMLHttpRequest', () => new Promise((resolve, reject) => {
   const xhr = new XMLHttpRequest();
   xhr.open('POST', '/echo');
   xhr.responseType = 'arraybuffer';
   xhr.onload = () => resolve(xhr.response);
   xhr.onerror = () => reject(new Error('XMLHttpRequest failed'));
   xhr.send(bytes);
 }));
 const controller = new AbortController();
 const stream = await fetch('/stream', {signal: controller.signal});
 const reader = stream.body.getReader();
 const chunk = await reader.read();
 if (chunk.done) throw new Error('stream ended before cancellation');
 controller.abort();
 await fetch('/done');
} catch (error) { await fetch('/done?error=' + encodeURIComponent(String(error))); }
</script></body></html>`)
			}
		})},
	})
	window = app.Window.NewWithOptions(application.WebviewWindowOptions{Title: "Native parity probe", Width: 640, Height: 480, URL: "/"})
	go func() {
		<-time.After(45 * time.Second)
		fmt.Fprintln(os.Stderr, "native parity probe timed out")
		exit(1)
	}()
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exit(1)
	}
	if !completed.Load() {
		fmt.Fprintln(os.Stderr, "application exited before completing probe")
		exit(1)
	}
}

func checkNativeAPIs(app *application.App, window *application.WebviewWindow) error {
	if width, height := window.Size(); width <= 0 || height <= 0 {
		return fmt.Errorf("invalid window size: %dx%d", width, height)
	}
	if screen, err := window.GetScreen(); err != nil || screen == nil {
		return fmt.Errorf("window screen: %v", err)
	}
	window.SetTitle("Native API probe")
	window.SetResizable(false)
	window.SetResizable(true)
	window.SetBackgroundColour(application.RGBA{Red: 40, Green: 50, Blue: 60, Alpha: 255})
	menu := application.NewMenu()
	items := menu.AddSubmenu("Probe")
	check := items.AddCheckbox("Checked", true)
	items.AddRadio("First", true)
	items.AddRadio("Second", false)
	items.AddSeparator()
	items.Add("Action").SetAccelerator("Ctrl+M")
	application.InvokeSync(menu.Update)
	window.SetMenu(menu)
	check.SetChecked(false)
	check.SetLabel("Updated")
	check.SetHidden(true)
	check.SetHidden(false)
	application.InvokeSync(menu.Update)
	window.SetMenu(nil)
	// Hide and Show are queued to the main thread on macOS (both backends), so
	// the state is polled rather than read straight away.
	window.Hide()
	if !eventually(func() bool { return !window.IsVisible() }) {
		return fmt.Errorf("hidden window remains visible")
	}
	window.Show()
	if !eventually(window.IsVisible) {
		return fmt.Errorf("shown window is not visible")
	}
	if err := checkPlatform(app, window); err != nil {
		return err
	}
	// Checking a radio item unchecks the rest of its group (pg-n5z).
	radios := application.NewMenu()
	first := radios.AddRadio("First", true)
	second := radios.AddRadio("Second", false)
	application.InvokeSync(radios.Update)
	second.SetChecked(true)
	if first.Checked() || !second.Checked() {
		return fmt.Errorf("radio SetChecked: first=%v second=%v", first.Checked(), second.Checked())
	}
	if os.Getenv("WAILS_TEST_DIALOGS") == "1" {
		if _, err := exec.LookPath("xdotool"); err != nil {
			return fmt.Errorf("dialog automation requires xdotool: %w", err)
		}
		result := make(chan struct{}, 1)
		dialog := app.Dialog.Info().SetTitle("Parity message").SetMessage("Native dialog callback")
		dialog.AddButton("Accept").SetAsDefault().OnClick(func() { result <- struct{}{} })
		keys := make(chan error, 1)
		go func() { keys <- pressDialogKey("Parity message", "Return") }()
		dialog.Show()
		if err := <-keys; err != nil {
			return err
		}
		select {
		case <-result:
		case <-time.After(5 * time.Second):
			return fmt.Errorf("message dialog callback timed out")
		}
		go func() { keys <- pressDialogKey("Parity file", "Escape") }()
		file, err := app.Dialog.OpenFile().SetTitle("Parity file").AttachToWindow(window).AddFilter("Text", "*.txt").PromptForSingleSelection()
		if err != nil || file != "" {
			return fmt.Errorf("cancelled file dialog returned %q, %v", file, err)
		}
		if err := <-keys; err != nil {
			return err
		}
		fmt.Println("PASS: message dialog callback and file dialog cancellation")
	}
	return nil
}

// Only enable dialog automation inside an isolated X server.
func pressDialogKey(title, key string) error {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		output, err := exec.Command("xdotool", "search", "--onlyvisible", "--name", "^"+title+"$").Output()
		ids := strings.Fields(string(output))
		// GTK4's message dialog renders the title as a heading inside an
		// otherwise untitled window. Find that application's other window.
		if title == "Parity message" && (err != nil || len(ids) == 0) {
			output, err = exec.Command("xdotool", "search", "--onlyvisible", "--pid", strconv.Itoa(os.Getpid())).Output()
			for _, id := range strings.Fields(string(output)) {
				name, nameErr := exec.Command("xdotool", "getwindowname", id).Output()
				if nameErr == nil && strings.TrimSpace(string(name)) != "Native API probe" {
					ids = append(ids, id)
				}
			}
		}
		if err == nil && len(ids) != 0 {
			id := ids[len(ids)-1]
			if err := exec.Command("xdotool", "windowfocus", "--sync", id).Run(); err != nil {
				return err
			}
			return exec.Command("xdotool", "key", "--window", id, key).Run()
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("dialog %q did not appear", title)
}

// eventually polls cond for up to two seconds.
func eventually(cond func() bool) bool {
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		if cond() {
			return true
		}
	}
	return cond()
}
