//go:build windows
// +build windows

package notification

import (
	"log"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"gopkg.in/toast.v1"
)

func SendNotification(options frontend.NotificationOptions) error {

	actions := []toast.Action{}
	if options.WindowsOptions.Actions != nil {
		for _, action := range options.WindowsOptions.Actions {
			actions = append(actions, toast.Action{
				Type:      action.Type,
				Label:     action.Label,
				Arguments: action.Arguments,
			})
		}
	}

	noty := toast.Notification{
		AppID:   options.AppID,
		Title:   options.Title,
		Message: options.Message,
		Actions: actions,
	}

	if options.AppIcon != nil {
		if p, err := AppIconPath(options.AppIcon); err == nil {
			noty.Icon = p
		}
	}

	/*
		if options.WindowsOptions != nil && options.WindowsOptions.Sound != "" {
			a, err := toast.Audio(options.WindowsOptions.Sound); err == nil {
				noty.Audio = a
			}
		}
	*/

	err := noty.Push()
	if err != nil {
		log.Fatalln(err)
	}

	return err
}
