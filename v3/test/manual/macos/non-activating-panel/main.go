// Manual test for MacWindow.NonActivatingPanel.
//
// Opens two windows side by side:
//   - "Panel"  - NonActivatingPanel: true (the feature under test)
//   - "Normal" - default WebviewWindow (the control)
//
// See ../README.md for the verification checklist.
//
//go:build darwin

package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const panelHTML = `<!DOCTYPE html>
<html><head><style>
  body { margin: 0; padding: 20px; font: 14px -apple-system;
         color: #fff; background: rgba(30,30,30,0.85);
         -webkit-app-region: drag; height: 100vh; box-sizing: border-box; }
  h1 { margin: 0 0 6px; font-size: 15px; }
  p  { margin: 0 0 12px; opacity: 0.7; font-size: 12px; }
  input { -webkit-app-region: no-drag; width: 100%; padding: 8px;
          background: rgba(255,255,255,0.1); color: #fff;
          border: 1px solid rgba(255,255,255,0.25); border-radius: 4px; font: 13px monospace; }
</style></head><body>
  <h1>Non-Activating Panel</h1>
  <p>Click the input. Wails should NOT take focus from your other app.</p>
  <input placeholder="Type to verify text input still works" autofocus />
</body></html>`

const normalHTML = `<!DOCTYPE html>
<html><head><style>
  body { margin: 0; padding: 20px; font: 14px -apple-system;
         color: #222; background: #f5f5f5; height: 100vh; box-sizing: border-box; }
  h1 { margin: 0 0 6px; font-size: 15px; }
  p  { margin: 0 0 12px; color: #555; font-size: 12px; }
  input { width: 100%; padding: 8px; border: 1px solid #ccc;
          border-radius: 4px; font: 13px monospace; }
</style></head><body>
  <h1>Normal Window (control)</h1>
  <p>Click the input. Wails SHOULD activate (compare with the panel).</p>
  <input placeholder="Type to verify text input works" />
</body></html>`

func main() {
	app := application.New(application.Options{
		Name: "NonActivatingPanel manual test",
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Panel",
		Width:     360,
		Height:    180,
		X:         400,
		Y:         300,
		Frameless: true,
		HTML:      panelHTML,
		Mac: application.MacWindow{
			NonActivatingPanel: true,
			WindowLevel:        application.MacWindowLevelFloating,
			Backdrop:           application.MacBackdropTranslucent,
			CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
				application.MacWindowCollectionBehaviorFullScreenAuxiliary,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Normal",
		Width:  360,
		Height: 180,
		X:      790,
		Y:      300,
		HTML:   normalHTML,
	})

	log.Println("Two windows opened. See README.md for the verification checklist.")
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
