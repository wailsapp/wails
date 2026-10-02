// This smoke test requires a graphical desktop (or Xvfb on Linux).
// Run with: go run ./test/manual/webview-panel
package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	loaded := make(chan string, 16)
	fail := make(chan error, 1)
	report := func(err error) {
		select {
		case fail <- err:
		default:
		}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/initial" {
			if r.Header.Get("X-Panel-Test") != "initial-only" {
				report(fmt.Errorf("initial header missing"))
			}
			if r.UserAgent() != "WailsPanelSmoke/1.0" {
				report(fmt.Errorf("user agent: %q", r.UserAgent()))
			}
			loaded <- "external"
		}
		fmt.Fprint(w, "<!doctype html><title>External panel</title><input autofocus placeholder='Panel input'>")
	}))
	defer server.Close()

	assets := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		switch r.URL.Path {
		case "/panel.html":
			fmt.Fprint(w, `<!doctype html><title>Local panel</title><script>fetch('/panel-ready')</script><p>Local asset panel</p>`)
		case "/panel-ready":
			loaded <- "local"
		case "/parent-ready":
			loaded <- "runtime"
		case "/failure":
			report(fmt.Errorf("frontend: %s", r.URL.RawQuery))
		default:
			fmt.Fprint(w, `<!doctype html><title>Panel smoke test</title><script type="module">
import {Panel} from '/wails/runtime.js';
try {
 const panel = Panel.Get('external');
 if (typeof await panel.IsFocused() !== 'boolean') throw Error('focus result');
 if (await panel.Name() !== 'external') throw Error('panel destroyed by focus query');
 if ('SetURL' in panel || 'ExecJS' in panel || 'SetHTML' in panel) throw Error('content IPC exposed');
 await fetch('/parent-ready');
} catch (error) { await fetch('/failure?'+encodeURIComponent(String(error))); }
</script><p>Testing embedded native webviews…</p>`)
		}
	})
	app := application.New(application.Options{Name: "Panel smoke test", Assets: application.AssetOptions{Handler: assets}})
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{Width: 800, Height: 600, URL: "/"})
	hidden := false
	panel := window.NewPanel(application.WebviewPanelOptions{
		Name: "external", URL: server.URL + "/initial", X: 20, Y: 80, Width: 300, Height: 200, Visible: &hidden,
		Headers: map[string]string{"X-Panel-Test": "initial-only"}, UserAgent: "WailsPanelSmoke/1.0",
	})
	local := window.NewPanel(application.WebviewPanelOptions{Name: "local", URL: "/panel.html", X: 340, Y: 80, Width: 300, Height: 200})
	var result error
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer app.Quit()
		seen := map[string]bool{}
		deadline := time.After(45 * time.Second)
		for len(seen) < 3 {
			select {
			case name := <-loaded:
				seen[name] = true
			case err := <-fail:
				result = err
				return
			case <-deadline:
				result = fmt.Errorf("load timeout: completed %v", seen)
				return
			}
		}
		// Flush past the native creation calls before querying initial state.
		application.InvokeSync(func() {})
		if panel.IsVisible() {
			result = fmt.Errorf("initially hidden panel is visible")
			return
		}
		panel.Show()
		if !panel.IsVisible() {
			result = fmt.Errorf("Show failed")
			return
		}
		panel.Hide()
		if panel.IsVisible() {
			result = fmt.Errorf("Hide failed")
			return
		}
		bounds := application.Rect{X: 30, Y: 90, Width: 280, Height: 180}
		panel.SetBounds(bounds)
		if actual := panel.Bounds(); actual != bounds {
			result = fmt.Errorf("bounds: %+v, want %+v", actual, bounds)
			return
		}
		panel.SetZoom(1.25)
		if math.Abs(panel.GetZoom()-1.25) > 0.01 {
			result = fmt.Errorf("zoom did not change")
			return
		}
		panel.SetZIndex(10)
		local.SetZIndex(2)
		extra := window.NewPanel(application.WebviewPanelOptions{Name: "dynamic", Width: 100, Height: 80})
		if extra == nil {
			result = fmt.Errorf("dynamic creation failed")
			return
		}
		extra.Destroy()
		var wg sync.WaitGroup
		for range 8 {
			wg.Add(1)
			go func() { defer wg.Done(); panel.Destroy() }()
		}
		wg.Wait()
		if window.GetPanel("external") != nil || len(window.GetPanels()) != 1 {
			result = fmt.Errorf("panel removal failed")
			return
		}
		local.Destroy()
	}()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
	<-done
	if result != nil {
		log.Fatal(result)
	}
	fmt.Println("PASS: panel runtime, local assets, headers, user agent, visibility, bounds, zoom, dynamic creation and concurrent destruction")
}
