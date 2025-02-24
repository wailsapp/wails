package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	NotificationService := notifications.New()
	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "notifications_demo",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(NotificationService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	go func() {
		app.OnEvent("notificationResponse", func(event *application.CustomEvent) {
			data, _ := json.Marshal(event.Data)
			fmt.Printf("%s\n", string(data))
		})

		NotificationService.RequestUserNotificationAuthorization()
		NotificationService.CheckNotificationAuthorization()

		time.Sleep(time.Second * 2)
		NotificationService.SendNotification(notifications.NotificationOptions{
			ID:    "Wails Notification Demo",
			Title: "Title!",
			Body:  "Body!",
			Data: map[string]interface{}{
				"messageId": "msg-123",
				"senderId":  "user-123",
				"timestamp": time.Now().String(),
			},
		})

		time.Sleep(time.Second * 2)
		NotificationService.RegisterNotificationCategory(notifications.NotificationCategory{
			ID: "BACKEND_NOTIF",
			Actions: []notifications.NotificationAction{
				{ID: "VIEW_ACTION", Title: "View"},
				{ID: "MARK_READ_ACTION", Title: "Mark as Read"},
				{ID: "DELETE_ACTION", Title: "Delete", Destructive: true},
			},
			HasReplyField:    true,
			ReplyButtonTitle: "Reply",
			ReplyPlaceholder: "Reply to backend...",
		})

		NotificationService.SendNotificationWithActions(notifications.NotificationOptions{
			ID:         "Wails Notification Demo",
			Title:      "Complex Backend Notification",
			Subtitle:   "Should not show on Windows",
			Body:       "Is it raining today where you are?",
			CategoryID: "BACKEND_NOTIF",
			Data: map[string]interface{}{
				"messageId": "msg-456",
				"senderId":  "user-456",
				"timestamp": time.Now().String(),
			},
		})
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
