//go:build !darwin

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func checkPlatform(*application.App, *application.WebviewWindow) error { return nil }
