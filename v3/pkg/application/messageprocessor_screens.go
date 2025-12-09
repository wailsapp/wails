package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	ScreensGetAll     = 0
	ScreensGetPrimary = 1
	ScreensGetCurrent = 2
)

var screensMethodNames = map[int]string{
	ScreensGetAll:     "GetAll",
	ScreensGetPrimary: "GetPrimary",
	ScreensGetCurrent: "GetCurrent",
}

func (m *MessageProcessor) processScreensMethod(req *RuntimeRequest) (any, error) {
	switch req.Method {
	case ScreensGetAll:
		return globalApplication.Screen.GetAll(), nil
	case ScreensGetPrimary:
		return globalApplication.Screen.GetPrimary(), nil
	case ScreensGetCurrent:
		screen, err := globalApplication.Window.Current().GetScreen()
		if err != nil {
			return nil, errs.WrapInvalidScreensCallErrorf(err, "Window.GetScreen failed")
		}
		return screen, nil
	default:
		return nil, errs.NewInvalidScreensCallErrorf("Unknown method: %d", req.Method)
	}
}
