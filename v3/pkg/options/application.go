package options

import "github.com/wailsapp/wails/v3/pkg/logger"

type Application struct {
	Name        string
	Description string
	Icon        []byte
	Mac         Mac
	Bind        []interface{}
	Logger      struct {
		Silent        bool
		CustomLoggers []logger.Output
	}
}
