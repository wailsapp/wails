package main

import (
	"embed"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
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
	// Create a new Notification Service
	ns := notifications.New()

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "Notifications Demo",
		Description: "A demo of using desktop notifications with Wails",
		Services: []application.Service{
			application.NewService(ns),
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

	app.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		// Create a goroutine that spawns desktop notifications from Go
		go func() {
			var authorized bool
			var err error
			authorized, err = ns.CheckNotificationAuthorization()
			if err != nil {
				println(fmt.Errorf("checking app notification authorization failed: %s", err))
			}

			if !authorized {
				authorized, err = ns.RequestNotificationAuthorization()
				if err != nil {
					println(fmt.Errorf("requesting app notification authorization failed: %s", err))
				}
			}

			if authorized {
				ns.OnNotificationResponse(func(result notifications.NotificationResult) {
					if result.Error != nil {
						println(fmt.Errorf("parsing notification result failed: %s", result.Error))
					} else {
						fmt.Printf("Response: %+v\n", result.Response)
						println("Sending response to frontend...")
						app.EmitEvent("notification:action", result.Response)
					}
				})

				err = ns.SendNotification(notifications.NotificationOptions{
					ID:       "uuid-basic-1",
					Title:    "Notification Title",
					Subtitle: "Subtitle on macOS and Linux",
					Body:     "Body text of notification.",
					Data: map[string]interface{}{
						"user-id":    "user-123",
						"message-id": "msg-123",
						"timestamp":  time.Now().Unix(),
					},
				})
				if err != nil {
					println(fmt.Errorf("sending basic notification failed: %s", err))
				}

				// Delay before sending next notification
				time.Sleep(time.Second * 2)

				const CategoryID = "backend-notification-id"

				err = ns.RegisterNotificationCategory(notifications.NotificationCategory{
					ID: CategoryID,
					Actions: []notifications.NotificationAction{
						{ID: "VIEW", Title: "View"},
						{ID: "MARK_READ", Title: "Mark as read"},
						{ID: "DELETE", Title: "Delete", Destructive: true},
					},
					HasReplyField:    true,
					ReplyPlaceholder: "Message...",
					ReplyButtonTitle: "Reply",
				})
				if err != nil {
					println(fmt.Errorf("creating notification category failed: %s", err))
				}

				err = ns.SendNotificationWithActions(notifications.NotificationOptions{
					ID:         "uuid-with-actions-1",
					Title:      "Actions Notification Title",
					Subtitle:   "Subtitle on macOS and Linux",
					Body:       "Body text of notification with actions.",
					CategoryID: CategoryID,
					Data: map[string]interface{}{
						"user-id":    "user-123",
						"message-id": "msg-123",
						"timestamp":  time.Now().Unix(),
					},
				})
				if err != nil {
					println(fmt.Errorf("sending notification with actions failed: %s", err))
				}
			}
		}()
	})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
