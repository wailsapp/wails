//go:build darwin

package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
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

	newPanel := func() *application.WebviewWindow {
		return app.Window.NewWithOptions(application.WebviewWindowOptions{
			Name:          "non-activating-panel",
			Title:         "Non-activating panel",
			URL:           "/panel",
			Width:         560,
			Height:        180,
			X:             480,
			Y:             90,
			Frameless:     true,
			DisableResize: true,
			// Let the panel page announce that it is ready, then exercise the
			// explicit non-activating Show path below.
			Hidden: true,
			KeyBindings: map[string]func(application.Window){
				"escape": func(window application.Window) {
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
	}

	panel = newPanel()

	app.Event.On("panel-ready", func(*application.CustomEvent) {
		panelMu.RLock()
		current := panel
		panelMu.RUnlock()
		if current != nil {
			current.Show()
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
<script type="module">
import * as wails from "/wails/runtime.js"; window.wails = wails;
wails.Events.Emit('panel-ready');
document.querySelector('#action').addEventListener('click', () => {
  document.querySelector('#status').textContent = 'Action handled at ' + new Date().toLocaleTimeString();
});
</script></body></html>`
