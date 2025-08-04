package main

import (
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "JumpList Example",
		Description: "A Wails application demonstrating Windows Jump Lists",
		Assets: application.AssetOptions{
			FS: assets,
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: false,
		},
	})

	// Create window
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Jump List Example",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	// Create jump list (Windows only - no-op on other platforms)
	jumpList := app.CreateJumpList()
	if jumpList != nil {
		// Add recent documents category
		recentCategory := application.JumpListCategory{
			Name: "Recent Documents",
			Items: []application.JumpListItem{
				{
					Type:        application.JumpListItemTypeTask,
					Title:       "Open Document 1",
					Description: "Open the first document",
					FilePath:    os.Args[0], // Using current executable for demo
					Arguments:   "--open doc1.txt",
					IconPath:    os.Args[0],
					IconIndex:   0,
				},
				{
					Type:        application.JumpListItemTypeTask,
					Title:       "Open Document 2",
					Description: "Open the second document",
					FilePath:    os.Args[0],
					Arguments:   "--open doc2.txt",
					IconPath:    os.Args[0],
					IconIndex:   0,
				},
			},
		}
		jumpList.AddCategory(recentCategory)

		// Add tasks (appears at the bottom of the jump list)
		tasksCategory := application.JumpListCategory{
			Name: "", // Empty name means tasks
			Items: []application.JumpListItem{
				{
					Type:        application.JumpListItemTypeTask,
					Title:       "New Document",
					Description: "Create a new document",
					FilePath:    os.Args[0],
					Arguments:   "--new",
					IconPath:    os.Args[0],
					IconIndex:   0,
				},
				{
					Type:        application.JumpListItemTypeTask,
					Title:       "Open Settings",
					Description: "Open application settings",
					FilePath:    os.Args[0],
					Arguments:   "--settings",
					IconPath:    os.Args[0],
					IconIndex:   0,
				},
			},
		}
		jumpList.AddCategory(tasksCategory)

		// Apply the jump list
		err := jumpList.Apply()
		if err != nil {
			log.Printf("Failed to apply jump list: %v", err)
		} else {
			log.Println("Jump list applied successfully")
		}

		// You can also clear and update the jump list at runtime
		window.OnWindowEvent(application.WindowEventReady, func(event *application.WindowEvent) {
			log.Println("Window ready - Jump list can be updated at any time")
		})
	}

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}