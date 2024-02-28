package server

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type Window struct {
	id     uint
	server *Server
}

// formatJS ensures the 'data' provided marshals to valid json or panics
func (w Window) formatJS(f string, callID string, data string) string {
	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(f, callID, j)
}

func (w Window) AbsolutePosition() (x, y int) {
	return 0, 0
}

func (w Window) CallError(callID string, result string) {
	w.ExecJS(callID, w.formatJS("_wails.callErrorCallback('%s', %s);", callID, result))
}

func (w Window) CallResponse(callID string, result string) {
	ids := strings.Split(callID, "|")
	j, err := json.Marshal(callback{
		ID:     ids[1],
		Result: result,
	})
	if err != nil {
		fmt.Println("Failed to build CallResponse data", result)
	}

	w.server.sendToClient(ids[0], message{Type: "cb", Data: string(j)})
}

func (w Window) DialogError(dialogID string, result string) {
	w.ExecJS(dialogID, w.formatJS("_wails.dialogErrorCallback('%s', %s);", dialogID, result))
}

func (w Window) DialogResponse(dialogID string, result string, isJSON bool) {
	cmd := "_wails.dialogResultCallback('%s', %s, true);"
	if !isJSON {
		cmd = "_wails.dialogResultCallback('%s', %s, false);"
	}
	w.ExecJS(dialogID, w.formatJS(cmd, dialogID, result))
}

func (w Window) ID() uint {
	return w.id
}

func (w Window) Center() {

}

func (w Window) Close() {}

func (w Window) Destroy() {}

func (w Window) ExecJS(callID, js string) {
	w.server.sendToClient(callID, message{
		Type: "javascript",
		Data: js,
	})
}

func (w Window) Focus() {}

func (w Window) ForceReload() {}

func (w Window) Fullscreen() application.Window {
	return w
}

func (w Window) GetScreen() (*application.Screen, error) {
	return nil, fmt.Errorf("can't return screen for external window")
}

func (w Window) GetZoom() float64 {
	return 1.0
}

func (w Window) Height() int {
	return 0
}

func (w Window) Hide() application.Window {
	return w
}

func (w Window) IsFullscreen() bool {
	return false
}

func (w Window) IsMaximised() bool {
	return false
}

func (w Window) IsMinimised() bool {
	return false
}

func (w Window) Maximise() application.Window {
	return w
}

func (w Window) Minimise() application.Window {
	return w
}

func (w Window) Minimize() {}

func (w Window) Name() string {
	return "external window"
}

func (w Window) On(eventType events.WindowEventType, callback func(ctx *application.WindowEvent)) func() {
	return func() {
		fmt.Printf("server.Window.On(%v)\n", eventType)
	}
}

func (w Window) Position() (int, int) {
	return 0, 0
}

func (w Window) RegisterContextMenu(name string, menu *application.Menu) {}

func (w Window) RelativePosition() (x, y int) {
	return 0, 0
}

func (w Window) Reload() {}

func (w Window) Resizable() bool {
	return true
}

func (w Window) Restore() {}

func (w Window) SetAbsolutePosition(x, y int) {}

func (w Window) SetAlwaysOnTop(b bool) application.Window {
	return w
}

func (w Window) SetBackgroundColour(colour application.RGBA) application.Window {
	return w
}

func (w Window) SetFrameless(frameless bool) application.Window {
	return w
}

func (w Window) SetFullscreenButtonEnabled(enabled bool) application.Window {
	return w
}

func (w Window) SetHTML(html string) application.Window {
	return w
}

func (w Window) SetMaxSize(maxWidth, maxHeight int) application.Window {
	return w
}

func (w Window) SetMinSize(minWidth, minHeight int) application.Window {
	return w
}

func (w Window) SetRelativePosition(x, y int) application.Window {
	return w
}

func (w Window) SetResizable(b bool) application.Window {
	return w
}

func (w Window) SetSize(width, height int) application.Window {
	return w
}

func (w Window) SetTitle(title string) application.Window {
	return w
}

func (w Window) SetURL(s string) application.Window {
	return w
}

func (w Window) SetZoom(magnification float64) application.Window {
	return w
}

func (w Window) Show() application.Window {
	return w
}

func (w Window) Size() (width int, height int) {
	return 0, 0
}

func (w Window) ToggleDevTools() {
}

func (w Window) ToggleFullscreen() {}

func (w Window) ToggleMaximise() {}

func (w Window) UnFullscreen() {}

func (w Window) UnMaximise() {}

func (w Window) UnMinimise() {}

func (w Window) Width() int {
	return 0
}

func (w Window) Zoom() {}

func (w Window) ZoomIn() {}

func (w Window) ZoomOut() {}

func (w Window) ZoomReset() application.Window {
	return w
}

func (w Window) DisableSizeConstraints() {}

func (w Window) DispatchWailsEvent(event *application.WailsEvent) {
	w.server.sendToAllClients(
		message{
			Type: "wailsevent",
			Data: event.ToJSON(),
		})
}

func (w Window) EnableSizeConstraints() {}

func (w Window) Error(message string, args ...any) {}

func (w Window) HandleDragAndDropMessage(filenames []string) {

}

func (w Window) HandleKeyEvent(acceleratorString string) {

}

func (w Window) HandleMessage(message string) {
	log.Println("HandleMessage", message)
}

func (w Window) HandleWindowEvent(id uint) {}

func (w Window) Info(message string, args ...any) {}

func (w Window) OpenContextMenu(data *application.ContextMenuData) {}

func (w Window) Run() {}
