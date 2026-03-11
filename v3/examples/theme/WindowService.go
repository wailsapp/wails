package main

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type WindowService struct {
	app *application.App
}

func (s *WindowService) SetAppTheme(theme string) {
	s.app.SetTheme(application.AppTheme(theme))
}

func (s *WindowService) GetAppTheme() string {
	return s.app.GetTheme().String()
}

func (s *WindowService) SetWinTheme(ctx context.Context, theme string) {
	win := ctx.Value(application.WindowKey).(application.Window)
	if win == nil {
		return
	}
	win.SetTheme((application.WinTheme(theme)))
}

func (s *WindowService) GetWinTheme(ctx context.Context) string {
	win := ctx.Value(application.WindowKey).(application.Window)
	if win == nil {
		return ""
	}
	return win.GetTheme().String()
}
