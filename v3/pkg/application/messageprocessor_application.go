package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	ApplicationHide = 0
	ApplicationShow = 1
	ApplicationQuit = 2
)

var applicationMethodNames = map[int]string{
	ApplicationQuit: "Quit",
	ApplicationHide: "Hide",
	ApplicationShow: "Show",
}

func (m *MessageProcessor) processApplicationMethod(
	req *RuntimeRequest,
) (any, error) {
	switch req.Method {
	case ApplicationQuit:
		globalApplication.Quit()
		return unit, nil
	case ApplicationHide:
		globalApplication.Hide()
		return unit, nil
	case ApplicationShow:
		globalApplication.Show()
		return unit, nil
	default:
		return nil, errs.NewInvalidApplicationCallErrorf("unknown method %d", req.Method)
	}
}
