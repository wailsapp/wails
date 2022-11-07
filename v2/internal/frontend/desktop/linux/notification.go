//go:build linux
// +build linux

package linux

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
)

func (f *Frontend) SendNotification(options frontend.NotificationOptions) error {
	return f.notifier.SendNotification(options)
}
