package other

import "github.com/wailsapp/wails/v3/pkg/application"

type Service13 int

var LocalService = application.NewService(new(Service13))
