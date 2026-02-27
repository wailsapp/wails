package main

import (
	"fmt"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type WindowService struct{}

func (w *WindowService) OpenTabbedWindow() {
	w.openWindow("Tabbed Window", application.MacWindowTabbingModePreferred)
}

func (w *WindowService) OpenNonTabbedWindow() {
	w.openWindow("Non-Tabbed Window", application.MacWindowTabbingModeDisallowed)
}

func (w *WindowService) openWindow(titlePrefix string, tabbingMode application.MacWindowTabbingMode) {
	app := application.Get()
	if app == nil {
		return
	}

	timestamp := time.Now().Format("15:04:05")
	windowTitle := fmt.Sprintf("%s (%s)", titlePrefix, timestamp)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: windowTitle,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
			TabbingMode:             tabbingMode,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})
}
