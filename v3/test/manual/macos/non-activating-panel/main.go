//go:build darwin

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
	app := application.New(application.Options{
		Name: "NSPanel manual test",
		Assets: application.AssetOptions{
			Handler: http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				response.Header().Set("Content-Type", "text/html; charset=utf-8")
				if request.URL.Path == "/panel" {
					_, _ = response.Write([]byte(panelHTML))
					return
				}
				_, _ = response.Write([]byte(controlHTML))
			}),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	var panelMu sync.RWMutex
	var panel *application.WebviewWindow
	var actionCount atomic.Uint64
	var escapeCount atomic.Uint64
	var hideEventCount atomic.Uint64
	var mainEventCount atomic.Uint64

	publishCounters := func(window application.Window) {
		window.ExecJS(fmt.Sprintf(
			"window.updatePanelCounters?.(%d, %d, %d, %d)",
			actionCount.Load(),
			escapeCount.Load(),
			hideEventCount.Load(),
			mainEventCount.Load(),
		))
	}

	newPanel := func() *application.WebviewWindow {
		result := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Name:          "non-activating-panel",
			Title:         "Non-activating panel",
			URL:           "/panel",
			Width:         560,
			Height:        240,
			X:             480,
			Y:             90,
			Frameless:     true,
			DisableResize: true,
			// Let the panel page announce that it is ready, then exercise the
			// explicit non-activating Show path below.
			Hidden: true,
			KeyBindings: map[string]func(application.Window){
				"escape": func(window application.Window) {
					escapeCount.Add(1)
					publishCounters(window)
					window.Hide()
				},
			},
			Mac: application.MacWindow{
				Backdrop:    application.MacBackdropTranslucent,
				WindowClass: application.MacWindowClassPanel,
				PanelPreferences: application.MacPanelPreferences{
					NonActivating: true,
				},
				WindowLevel: application.MacWindowLevelFloating,
				CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
					application.MacWindowCollectionBehaviorFullScreenAuxiliary |
					application.MacWindowCollectionBehaviorStationary,
			},
		})
		result.OnWindowEvent(events.Common.WindowHide, func(*application.WindowEvent) {
			hideEventCount.Add(1)
			publishCounters(result)
		})
		result.OnWindowEvent(events.Mac.WindowWillBecomeMain, func(*application.WindowEvent) {
			mainEventCount.Add(1)
			publishCounters(result)
		})
		result.OnWindowEvent(events.Mac.WindowDidBecomeMain, func(*application.WindowEvent) {
			mainEventCount.Add(1)
			publishCounters(result)
		})
		return result
	}

	panel = newPanel()

	app.Event.On("panel-ready", func(*application.CustomEvent) {
		panelMu.RLock()
		current := panel
		panelMu.RUnlock()
		if current != nil {
			current.Show()
			publishCounters(current)
		}
	})
	app.Event.On("panel-action", func(*application.CustomEvent) {
		actionCount.Add(1)
		panelMu.RLock()
		current := panel
		panelMu.RUnlock()
		if current != nil {
			publishCounters(current)
		}
	})
	app.Event.On("panel-show", func(*application.CustomEvent) {
		panelMu.RLock()
		current := panel
		panelMu.RUnlock()
		if current != nil {
			current.Show()
		}
	})
	app.Event.On("panel-focus", func(*application.CustomEvent) {
		panelMu.RLock()
		current := panel
		panelMu.RUnlock()
		if current != nil {
			current.Focus()
		}
	})
	app.Event.On("panel-hide", func(*application.CustomEvent) {
		panelMu.RLock()
		current := panel
		panelMu.RUnlock()
		if current != nil {
			current.Hide()
		}
	})
	app.Event.On("panel-close", func(*application.CustomEvent) {
		panelMu.Lock()
		current := panel
		panel = nil
		panelMu.Unlock()
		if current != nil {
			current.Close()
		}
	})
	app.Event.On("panel-recreate", func(*application.CustomEvent) {
		panelMu.Lock()
		if panel == nil {
			panel = newPanel()
		}
		current := panel
		panelMu.Unlock()
		current.Show()
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "normal-window-control",
		Title:  "Normal NSWindow control",
		URL:    "/",
		Width:  420,
		Height: 360,
		X:      40,
		Y:      90,
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

const sharedStyle = `<style>
* { box-sizing: border-box; }
body { margin: 0; padding: 22px; font: 14px -apple-system, BlinkMacSystemFont, sans-serif; color: #f5f5f5; background: #17181c; }
h1 { margin: 0 0 8px; font-size: 20px; }
p { color: #aeb3be; line-height: 1.45; }
.actions { display: flex; flex-wrap: wrap; gap: 8px; margin: 16px 0; }
button, input { border: 1px solid #454957; border-radius: 7px; color: #fff; background: #292c34; padding: 9px 12px; }
button:hover { background: #383c47; }
input { width: 100%; background: #111216; }
#status { color: #7dd3fc; }
.metrics { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 14px; }
.metric { border: 1px solid #383c47; border-radius: 6px; padding: 5px 8px; color: #aeb3be; }
.metric strong { color: #f5f5f5; }
</style>`

const controlHTML = `<!doctype html><html><head><meta charset="utf-8"><title>NSWindow control</title>` + sharedStyle + `</head>
<body><h1>Normal NSWindow control</h1>
<p>This is the unchanged NSWindow path. Use it to show, focus, close, and recreate the NSPanel.</p>
<div class="actions">
<button onclick="wails.Events.Emit('panel-show')">Show panel</button>
<button onclick="wails.Events.Emit('panel-focus')">Focus panel</button>
<button onclick="wails.Events.Emit('panel-hide')">Hide panel</button>
<button onclick="wails.Events.Emit('panel-close')">Close panel</button>
<button onclick="wails.Events.Emit('panel-recreate')">Recreate panel</button>
</div>
<p>Keep another application active while using the panel. Its menu bar should remain active.</p>
<script type="module">import * as wails from "/wails/runtime.js"; window.wails = wails;</script>
</body></html>`

const panelHTML = `<!doctype html><html><head><meta charset="utf-8"><title>NSPanel</title>` + sharedStyle + `</head>
<body><h1>Non-activating NSPanel</h1>
<input autofocus placeholder="Type here; the other app should remain active">
<div class="actions"><button id="action">Run panel action</button><button onclick="wails.Events.Emit('panel-hide')">Dismiss</button></div>
<span id="status">No action yet</span>
<div class="metrics">
<span class="metric">Actions <strong id="action-count">0</strong></span>
<span class="metric">Escape callbacks <strong id="escape-count">0</strong></span>
<span class="metric">Hide events <strong id="hide-count">0</strong></span>
<span class="metric">Main events <strong id="main-count">0</strong></span>
</div>
<script type="module">
import * as wails from "/wails/runtime.js"; window.wails = wails;
window.updatePanelCounters = (actions, escapes, hides, mainEvents) => {
  document.querySelector('#action-count').textContent = actions;
  document.querySelector('#escape-count').textContent = escapes;
  document.querySelector('#hide-count').textContent = hides;
  document.querySelector('#main-count').textContent = mainEvents;
};
wails.Events.Emit('panel-ready');
document.querySelector('#action').addEventListener('click', () => {
  wails.Events.Emit('panel-action');
  document.querySelector('#status').textContent = 'Action handled at ' + new Date().toLocaleTimeString();
});
</script></body></html>`
