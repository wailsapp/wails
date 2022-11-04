//go:build darwin
// +build darwin

package notification

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
	notify "github.com/willdot/gomacosnotify"
)

func SendNotification(options frontend.NotificationOptions) error {
	var err error
	notifyOpts := notify.Notification{

		Title:    options.Title,
		SubTitle: options.MacOptions.SubTitle,
		Message:  options.Message,
	}

	if len(options.MacOptions.Actions) > 0 {
		notifyOpts.Actions = []string{}
		for i := 0; i < len(options.MacOptions.Actions); i++ {
			notifyOpts.Actions = append(notifyOpts.Actions, options.MacOptions.Actions[i].Label)
		}
	}

	if options.MacOptions.ContentImage != nil {
		i, err := ContentImagePath(options.MacOptions.ContentImage)
		if err == nil {
			notifyOpts.ContentImage = i
		}
	}

	if options.MacOptions.CloseText != "" {
		notifyOpts.CloseText = options.MacOptions.CloseText
	}

	notifyOpts.SetTimeout(int(options.Timeout.Seconds()))
	n, err := notify.New()
	if err != nil {
		return err
	}

	go func() {
		resp, err := n.Send(notifyOpts)
		if err != nil {
			return
		}

		for i := 0; i < len(options.MacOptions.Actions); i++ {
			if options.MacOptions.Actions[i].OnAction != nil && options.MacOptions.Actions[i].Label == resp.ActivationValue {
				options.MacOptions.Actions[i].OnAction(resp.ActivationType, resp.ActivationValue)
			}
		}

	}()

	return err
}
