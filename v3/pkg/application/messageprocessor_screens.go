package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	ScreensGetAll     = 0
	ScreensGetPrimary = 1
	ScreensGetCurrent = 2
	ScreensGetByID    = 3
	ScreensGetByIndex = 4
)

var screensMethodNames = map[int]string{
	ScreensGetAll:     "GetAll",
	ScreensGetPrimary: "GetPrimary",
	ScreensGetCurrent: "GetCurrent",
	ScreensGetByID:    "GetByID",
	ScreensGetByIndex: "GetByIndex",
}

func (m *MessageProcessor) processScreensMethod(req *RuntimeRequest) (any, error) {
	args := req.Args.AsMap()

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
	case ScreensGetByID:
		id := args.String("id")
		if id == nil {
			return nil, errs.NewInvalidScreensCallErrorf("missing or invalid argument 'id'")
		}
		screen := globalApplication.Screen.GetByID(*id)
		if screen == nil {
			return nil, errs.NewInvalidScreensCallErrorf("screen not found: %s", *id)
		}
		return screen, nil
	case ScreensGetByIndex:
		index := args.Int("index")
		if index == nil {
			return nil, errs.NewInvalidScreensCallErrorf("missing or invalid argument 'index'")
		}
		screen := globalApplication.Screen.GetByIndex(*index)
		if screen == nil {
			return nil, errs.NewInvalidScreensCallErrorf("screen not found at index: %d", *index)
		}
		return screen, nil
	default:
		return nil, errs.NewInvalidScreensCallErrorf("Unknown method: %d", req.Method)
	}
}
