package dispatcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

const systemCallPrefix = ":wails:"

type position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type size struct {
	W int `json:"w"`
	H int `json:"h"`
}

func (d *Dispatcher) processSystemCall(payload callMessage, sender frontend.Frontend) (interface{}, error) {
	// Strip prefix
	name := strings.TrimPrefix(payload.Name, systemCallPrefix)

	switch name {
	case "WindowGetPos":
		x, y := sender.WindowGetPosition()
		return &position{x, y}, nil
	case "WindowGetSize":
		w, h := sender.WindowGetSize()
		return &size{w, h}, nil
	case "ScreenGetAll":
		return sender.ScreenGetAll()
	case "WindowIsMaximised":
		return sender.WindowIsMaximised(), nil
	case "WindowIsMinimised":
		return sender.WindowIsMinimised(), nil
	case "WindowIsNormal":
		return sender.WindowIsNormal(), nil
	case "WindowIsFullscreen":
		return sender.WindowIsFullscreen(), nil
	case "Environment":
		return runtime.Environment(d.ctx), nil
	case "ClipboardGetText":
		t, err := sender.ClipboardGetText()
		return t, err
	case "ClipboardSetText":
		if len(payload.Args) < 1 {
			return false, errors.New("empty argument, cannot set clipboard")
		}
		var arg string
		if err := json.Unmarshal(payload.Args[0], &arg); err != nil {
			return false, err
		}
		if err := sender.ClipboardSetText(arg); err != nil {
			return false, err
		}
		return true, nil
	case "InitializeNotifications":
		err := sender.InitializeNotifications()
		return nil, err
	case "CleanupNotifications":
		sender.CleanupNotifications()
		return nil, nil
	case "IsNotificationAvailable":
		return sender.IsNotificationAvailable(), nil
	case "RequestNotificationAuthorization":
		authorized, err := sender.RequestNotificationAuthorization()
		if err != nil {
			return nil, err
		}
		return authorized, nil
	case "CheckNotificationAuthorization":
		authorized, err := sender.CheckNotificationAuthorization()
		if err != nil {
			return nil, err
		}
		return authorized, nil
	case "SendNotification":
		if len(payload.Args) < 1 {
			return nil, errors.New("empty argument, cannot send notification")
		}
		var options frontend.NotificationOptions
		if err := json.Unmarshal(payload.Args[0], &options); err != nil {
			return nil, err
		}
		err := sender.SendNotification(options)
		return nil, err
	case "SendNotificationWithActions":
		if len(payload.Args) < 1 {
			return nil, errors.New("empty argument, cannot send notification")
		}
		var options frontend.NotificationOptions
		if err := json.Unmarshal(payload.Args[0], &options); err != nil {
			return nil, err
		}
		err := sender.SendNotificationWithActions(options)
		return nil, err
	case "RegisterNotificationCategory":
		if len(payload.Args) < 1 {
			return nil, errors.New("empty argument, cannot register category")
		}
		var category frontend.NotificationCategory
		if err := json.Unmarshal(payload.Args[0], &category); err != nil {
			return nil, err
		}
		err := sender.RegisterNotificationCategory(category)
		return nil, err
	case "RemoveNotificationCategory":
		if len(payload.Args) < 1 {
			return nil, errors.New("empty argument, cannot remove category")
		}
		var categoryId string
		if err := json.Unmarshal(payload.Args[0], &categoryId); err != nil {
			return nil, err
		}
		err := sender.RemoveNotificationCategory(categoryId)
		return nil, err
	case "RemoveAllPendingNotifications":
		err := sender.RemoveAllPendingNotifications()
		return nil, err
	case "RemovePendingNotification":
		if len(payload.Args) < 1 {
			return nil, errors.New("empty argument, cannot remove notification")
		}
		var identifier string
		if err := json.Unmarshal(payload.Args[0], &identifier); err != nil {
			return nil, err
		}
		err := sender.RemovePendingNotification(identifier)
		return nil, err
	case "RemoveAllDeliveredNotifications":
		err := sender.RemoveAllDeliveredNotifications()
		return nil, err
	case "RemoveDeliveredNotification":
		if len(payload.Args) < 1 {
			return nil, errors.New("empty argument, cannot remove notification")
		}
		var identifier string
		if err := json.Unmarshal(payload.Args[0], &identifier); err != nil {
			return nil, err
		}
		err := sender.RemoveDeliveredNotification(identifier)
		return nil, err
	case "RemoveNotification":
		if len(payload.Args) < 1 {
			return nil, errors.New("empty argument, cannot remove notification")
		}
		var identifier string
		if err := json.Unmarshal(payload.Args[0], &identifier); err != nil {
			return nil, err
		}
		err := sender.RemoveNotification(identifier)
		return nil, err
	default:
		return nil, fmt.Errorf("unknown systemcall message: %s", payload.Name)
	}
}
