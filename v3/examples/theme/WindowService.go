package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type WindowService struct {
	app *application.App
}

func (s *WindowService) SetAppTheme(theme string) {
	s.app.SetTheme(application.AppTheme(theme))
}

func (s *WindowService) GetAppTheme() string {
	return s.app.GetTheme()
}

func (s *WindowService) SetWinTheme(theme string) {
	win := s.app.Window.Current()
	win.SetTheme((application.WinTheme(theme)))
}

func (s *WindowService) GetWinTheme() string {
	win := s.app.Window.Current()
	theme := win.GetTheme()
	return theme
}
