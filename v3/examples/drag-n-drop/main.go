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

// FileDropInfo defines the payload for the file drop event sent to the frontend.
type FileDropInfo struct {
	Files         []string          `json:"files"`
	TargetID      string            `json:"targetID"`
	TargetClasses []string          `json:"targetClasses"`
	DropX         float64           `json:"dropX"`
	DropY         float64           `json:"dropY"`
	Attributes    map[string]string `json:"attributes,omitempty"`
}

// FilesDroppedOnTarget is called when files are dropped onto a registered drop target
// or the window if no specific target is hit.
func FilesDroppedOnTarget(
	files []string,
	targetID string,
	targetClasses []string,
	dropX float64,
	dropY float64,
	isTargetDropzone bool, // This parameter is kept for logging but not sent to frontend in this event
	attributes map[string]string,
) {
	log.Println("=============== Go: FilesDroppedOnTarget Debug Info ===============")
	log.Println(fmt.Sprintf("  Files: %v", files))
	log.Println(fmt.Sprintf("  Target ID: '%s'", targetID))
	log.Println(fmt.Sprintf("  Target Classes: %v", targetClasses))
	log.Println(fmt.Sprintf("  Drop X: %f, Drop Y: %f", dropX, dropY))
	log.Println(
		fmt.Sprintf(
			"  Drop occurred on a designated dropzone (runtime validated before this Go event): %t",
			isTargetDropzone,
		),
	)
	log.Println(fmt.Sprintf("  Element Attributes: %v", attributes))
	log.Println("================================================================")

	payload := FileDropInfo{
		Files:         files,
		TargetID:      targetID,
		TargetClasses: targetClasses,
		DropX:         dropX,
		DropY:         dropY,
		Attributes:    attributes,
	}

	log.Println("Go: Emitted 'frontend:FileDropInfo' event with payload:", payload)
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
	})

	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Drag-n-drop Demo",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
		EnableDragAndDrop: true,
	})

	log.Println("Setting up event listener for 'WindowDropZoneFilesDropped'...")
	win.OnWindowEvent(
		events.Common.WindowDropZoneFilesDropped,
		func(event *application.WindowEvent) {

			droppedFiles := event.Context().DroppedFiles()
			details := event.Context().DropZoneDetails()

			log.Printf("Dropped files count: %d", len(droppedFiles))
			log.Printf("Event context: %+v", event.Context())

			if details != nil {
				log.Printf("DropZone details found:")
				log.Printf("  ElementID: '%s'", details.ElementID)
				log.Printf("  ClassList: %v", details.ClassList)
				log.Printf("  X: %d, Y: %d", details.X, details.Y)
				log.Printf("  Attributes: %+v", details.Attributes)

				// Call the App method with the extracted data
				FilesDroppedOnTarget(
					droppedFiles,
					details.ElementID,
					details.ClassList,
					float64(details.X),
					float64(details.Y),
					details.ElementID != "", // isTargetDropzone based on whether an ID was found
					details.Attributes,
				)
			} else {
				log.Println("DropZone details are nil - drop was not on a specific registered zone")
				// This case might occur if DropZoneDetails are nil, meaning the drop was not on a specific registered zone
				// or if the context itself was problematic.
				FilesDroppedOnTarget(droppedFiles, "", nil, 0, 0, false, nil)
			}

			payload := FileDropInfo{
				Files:         droppedFiles,
				TargetID:      details.ElementID,
				TargetClasses: details.ClassList,
				DropX:         float64(details.X),
				DropY:         float64(details.Y),
				Attributes:    details.Attributes, // Add the attributes
			}

			log.Printf("Emitting event payload: %+v", payload)
			application.Get().Event.Emit("frontend:FileDropInfo", payload)
			log.Println(
				"=============== End WindowDropZoneFilesDropped Event Debug ===============",
			)
		},
	)

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
