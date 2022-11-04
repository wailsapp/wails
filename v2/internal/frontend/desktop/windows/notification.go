//go:build windows
// +build windows

package windows

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/notification"
)

func (f *Frontend) SendNotification(options frontend.NotificationOptions) error {
	return notification.SendNotification(options)
}
