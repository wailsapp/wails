package application

import (
	"net/http"
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

func (m *MessageProcessor) processScreensMethod(method int, rw http.ResponseWriter, _ *http.Request, _ Window, _ QueryParams) {

	switch method {
	case ScreensGetAll:
		screens, err := globalApplication.GetScreens()
		if err != nil {
			m.Error("GetAll: %s", err.Error())
			return
		}
		m.json(rw, screens)
	case ScreensGetPrimary:
		screen, err := globalApplication.GetPrimaryScreen()
		if err != nil {
			m.Error("GetPrimary: %s", err.Error())
			return
		}
		m.json(rw, screen)
	case ScreensGetCurrent:
		screen, err := globalApplication.CurrentWindow().GetScreen()
		if err != nil {
			m.Error("GetCurrent: %s", err.Error())
			return
		}
		m.json(rw, screen)
	default:
		m.httpError(rw, "Unknown screens method: %d", method)
	}

	m.Info("Runtime Call:", "method", "Screens."+screensMethodNames[method])

}
