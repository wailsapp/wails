//go:build linux
// +build linux

package linux

import (
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/notification"
)

func (f *Frontend) SendNotification(options frontend.NotificationOptions) error {
	switch f.notifier.GetMethod() {
	case notification.MethodDbus:
		ID, err := f.notifier.SendViaDbus(options)
		if options.LinuxOptions.OnShow != nil {
			options.LinuxOptions.OnShow(ID)
		}
		return err
	case notification.MethodNotifySend:
		ID, err := f.notifier.SendViaNotifySend(options)
		if options.LinuxOptions.OnShow != nil {
			options.LinuxOptions.OnShow(ID)
		}
		return err
	case notification.MethodKdialog:
		ID, err := f.notifier.SendViaKnotify(options)
		if options.LinuxOptions.OnShow != nil {
			options.LinuxOptions.OnShow(ID)
		}
		return err
	}
	return errors.New("no notification method is available")
}
