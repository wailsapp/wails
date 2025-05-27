package main

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
	app *application.App
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.app = application.Get()
}

// FilesDroppedOnTarget is called from JavaScript when files are dropped.
func (a *App) FilesDroppedOnTarget(files []string, targetID string, targetClasses []string, dropX float64, dropY float64, isTargetDropzone bool) {
	a.app.Logger.Info("Go: Received 'FilesDroppedOnTarget' call from frontend")
	a.app.Logger.Info(fmt.Sprintf("  Files: %v", files))
	a.app.Logger.Info(fmt.Sprintf("  Target ID: %s", targetID))
	a.app.Logger.Info(fmt.Sprintf("  Target Classes: %v", targetClasses))
	a.app.Logger.Info(fmt.Sprintf("  Drop X: %f, Drop Y: %f", dropX, dropY))
	a.app.Logger.Info(fmt.Sprintf("  Is Target Dropzone: %t", isTargetDropzone))
}

func main() {
	appInstance := NewApp()

	app := application.New(application.Options{
		Name:        "Drag-n-drop Demo",
		Description: "A demo of the Drag-n-drop API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Services: []application.Service{
			application.NewService(appInstance),
		},
		// The Startup(ctx context.Context) method on appInstance should be called by convention
		// if application.Service is an interface that *App implements.
	})

	win := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Drag-n-drop Demo",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
		EnableDragAndDrop: true,
	})

	// Set up the event listener for enriched drag-and-drop events on this window
	win.OnWindowEvent(events.Common.WindowDropZoneFilesDropped, func(event *application.WindowEvent) {
		appInstance.app.Logger.Info("Go (Window Event): Received 'WindowDropZoneFilesDropped'")
		droppedFiles := event.Context().DroppedFiles()
		details := event.Context().DropZoneDetails()

		if details != nil {
			appInstance.app.Logger.Info(fmt.Sprintf("  Ctx X: %d, Y: %d, ID: '%s', Classes: %v", details.X, details.Y, details.ElementID, details.ClassList))
			appInstance.FilesDroppedOnTarget(droppedFiles, details.ElementID, details.ClassList, float64(details.X), float64(details.Y), details.ElementID != "")
		} else {
			appInstance.app.Logger.Info("  Ctx DropZoneDetails were nil")
			appInstance.FilesDroppedOnTarget(droppedFiles, "", nil, 0, 0, false)
		}
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
