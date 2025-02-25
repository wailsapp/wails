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

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	notificationService := notifications.New()

	app := application.New(application.Options{
		Name:        "notifications_demo",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(notificationService),
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

	go func() {
		granted, err := notificationService.RequestUserNotificationAuthorization()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
			return
		}

		if granted {
			notificationService.OnNotificationResponse(func(response notifications.NotificationResponse) {
				data, _ := json.Marshal(response)
				fmt.Printf("%s\n", string(data))
				app.EmitEvent("notification:response", response)
			})
			time.Sleep(time.Second * 2)

			var uuid1 string = "Wails Notification Demo"
			if application.Get().Environment().OS == "darwin" {
				uuid1 = "uuid1"
			}

			notificationService.SendNotification(notifications.NotificationOptions{
				ID:    uuid1,
				Title: "Title!",
				Body:  "Body!",
				Data: map[string]interface{}{
					"messageId": "msg-123",
					"senderId":  "user-123",
					"timestamp": time.Now().String(),
				},
			})

			time.Sleep(time.Second * 2)

			var uuid2 string = "Wails Notification Demo"
			if application.Get().Environment().OS == "darwin" {
				uuid2 = "uuid2"
			}

			notificationService.RegisterNotificationCategory(notifications.NotificationCategory{
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

			notificationService.SendNotificationWithActions(notifications.NotificationOptions{
				ID:         uuid2,
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
		}
	}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
