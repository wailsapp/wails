package devserver

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"net/http"
)

//go:embed runtime_debug_windows.js
var wailsFS embed.FS

type DevServer struct {
	server *http.Server
}

func (d DevServer) Run(ctx context.Context) error {
	wailsfs := http.FileServer(http.FS(wailsFS))
	http.Handle("/wails/", http.StripPrefix("/wails", wailsfs))

	assetdir := ctx.Value("assetdir")
	if assetdir != nil {
		println("Serving assets from:", assetdir.(string))
		fs := http.FileServer(http.Dir(assetdir.(string)))
		http.Handle("/", interceptor(fs))
	}
	err := d.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (d DevServer) Quit() {
	panic("implement me")
}

func (d DevServer) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	panic("implement me")
}

func (d DevServer) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	panic("implement me")
}

func (d DevServer) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	panic("implement me")
}

func (d DevServer) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	panic("implement me")
}

func (d DevServer) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	panic("implement me")
}

func (d DevServer) WindowSetTitle(title string) {
	panic("implement me")
}

func (d DevServer) WindowShow() {
	panic("implement me")
}

func (d DevServer) WindowHide() {
	panic("implement me")
}

func (d DevServer) WindowCenter() {
	panic("implement me")
}

func (d DevServer) WindowMaximise() {
	panic("implement me")
}

func (d DevServer) WindowUnmaximise() {
	panic("implement me")
}

func (d DevServer) WindowMinimise() {
	panic("implement me")
}

func (d DevServer) WindowUnminimise() {
	panic("implement me")
}

func (d DevServer) WindowSetPos(x int, y int) {
	panic("implement me")
}

func (d DevServer) WindowGetPos() (int, int) {
	panic("implement me")
}

func (d DevServer) WindowSetSize(width int, height int) {
	panic("implement me")
}

func (d DevServer) WindowGetSize() (int, int) {
	panic("implement me")
}

func (d DevServer) WindowSetMinSize(width int, height int) {
	panic("implement me")
}

func (d DevServer) WindowSetMaxSize(width int, height int) {
	panic("implement me")
}

func (d DevServer) WindowFullscreen() {
	panic("implement me")
}

func (d DevServer) WindowUnFullscreen() {
	panic("implement me")
}

func (d DevServer) WindowSetColour(colour int) {
	panic("implement me")
}

func (d DevServer) SetApplicationMenu(menu *menu.Menu) {
	panic("implement me")
}

func (d DevServer) UpdateApplicationMenu() {
	panic("implement me")
}

func (d DevServer) Notify(name string, data ...interface{}) {
	panic("implement me")
}

func interceptor(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("Serving file from disk:", r.RequestURI)
		nextHandler.ServeHTTP(w, r)
	})
}

func NewFrontend(appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *DevServer {
	result := &DevServer{}
	result.server = &http.Server{Addr: ":34115"}
	return result
}
