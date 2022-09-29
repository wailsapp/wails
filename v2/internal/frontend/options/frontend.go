package options

import "github.com/wailsapp/wails/v2/pkg/options"

// Frontend contains options for creating the Frontend
type Frontend struct {
	options.App
	HasMainWindow bool
}
