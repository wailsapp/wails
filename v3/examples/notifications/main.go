package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	NotificationService := notifications.New()

	app := application.New(application.Options{
		Name:        "notifications",
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

	app.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
		// Request pemission to send notifications
		granted, err := NotificationService.RequestUserNotificationAuthorization()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}

		if granted {
			// Send notification with no actions
			err = NotificationService.SendNotification("some-uuid", "Title", "", "body!")
			if err != nil {
				log.Default().Printf("WARNING: %s\n", err)
			}

			err = NotificationService.RegisterNotificationCategory(notifications.NotificationCategory{
				ID: "MESSAGE_CATEGORY",
				Actions: []notifications.NotificationAction{
					{ID: "VIEW_ACTION", Title: "View"},
					{ID: "MARK_READ_ACTION", Title: "Mark as Read"},
					{ID: "DELETE_ACTION", Title: "Delete", Destructive: true},
				},
				HasReplyField:    true,
				ReplyPlaceholder: "Type your reply...",
				ReplyButtonTitle: "Reply",
			})
			if err != nil {
				log.Default().Printf("WARNING: %s\n", err)
			}

			err = NotificationService.SendNotificationWithActions(notifications.NotificationOptions{
				ID:         "some-other-uuid",
				Title:      "New Message",
				Subtitle:   "From: Jane Doe",
				Body:       "Is it raining today where you are?",
				CategoryID: "MESSAGE_CATEGORY",
				Data: map[string]interface{}{
					"messageId": "msg-123",
					"senderId":  "user-123",
					"timestamp": time.Now().Unix(),
				},
			})
			if err != nil {
				log.Default().Printf("WARNING: %s\n", err)
			}
		}
	})

	app.OnEvent("notificationResponse", func(event *application.CustomEvent) {
		data, _ := json.Marshal(event.Data)
		fmt.Printf("%s\n", string(data))
	})

	go func() {
		time.Sleep(time.Second * 5)
		// Sometime later check if you are still authorized
		granted, err := NotificationService.CheckNotificationAuthorization()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}
		println(granted)

		// Sometime later remove delivered notification by identifier
		err = NotificationService.RemoveDeliveredNotification("some-uuid")
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}
	}()

	go func() {
		time.Sleep(time.Second * 10)
		// Sometime later remove all pending and delivered notifications
		err := NotificationService.RemoveAllPendingNotifications()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}
		err = NotificationService.RemoveAllDeliveredNotifications()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
