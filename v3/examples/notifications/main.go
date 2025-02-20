package main

import (
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
	app := application.New(application.Options{
		Name:        "Notifications Demo",
		Description: "A test of macOS notifications",
		Assets:      application.AlphaAssets,
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:  500,
		Height: 800,
	})

	app.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
		// Request pemission to send notifications
		granted, err := application.RequestUserNotificationAuthorization()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}

		if granted {
			// Send notification with no actions
			err = application.SendNotification("some-uuid", "Title", "", "body!")
			if err != nil {
				log.Default().Printf("WARNING: %s\n", err)
			}

			err = application.RegisterNotificationCategory(application.NotificationCategory{
				ID: "MESSAGE_CATEGORY",
				Actions: []application.NotificationAction{
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

			err = application.SendNotificationWithActions(application.NotificationOptions{
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

	app.OnApplicationEvent(events.Mac.DidReceiveNotificationResponse, func(event *application.ApplicationEvent) {
		data := event.Context().GetData()

		// Parse data received
		if data != nil {
			if identifier, ok := data["identifier"].(string); ok {
				fmt.Printf("Notification identifier: %s\n", identifier)
			}

			if actionIdentifier, ok := data["actionIdentifier"].(string); ok {
				fmt.Printf("Action Identifier: %s\n", actionIdentifier)
			}

			if userText, ok := data["userText"].(string); ok {
				fmt.Printf("User replied: %s\n", userText)
			}

			if userInfo, ok := data["userInfo"].(map[string]interface{}); ok {
				fmt.Printf("Custom data: %+v\n", userInfo)
			}

			// Send notification to JS
			app.EmitEvent("notification", data)
		}
	})

	go func() {
		time.Sleep(time.Second * 5)
		// Sometime later check if you are still authorized
		granted, err := application.CheckNotificationAuthorization()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}
		println(granted)

		// Sometime later remove delivered notification by identifier
		err = application.RemoveDeliveredNotification("some-uuid")
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}
	}()

	go func() {
		time.Sleep(time.Second * 10)
		// Sometime later remove all pending and delivered notifications
		err := application.RemoveAllPendingNotifications()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}
		err = application.RemoveAllDeliveredNotifications()
		if err != nil {
			log.Default().Printf("WARNING: %s\n", err)
		}
	}()

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
