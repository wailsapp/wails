package application

import (
	"fmt"
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
			m.httpError(rw, "GetScreens failed:", err)
			return
		}
		m.json(rw, screens)
	case ScreensGetPrimary:
		screen, err := globalApplication.GetPrimaryScreen()
		if err != nil {
			m.httpError(rw, "GetPrimary failed:", err)
			return
		}
		m.json(rw, screen)
	case ScreensGetCurrent:
		screen, err := globalApplication.CurrentWindow().GetScreen()
		if err != nil {
			m.httpError(rw, "Window.GetScreen failed:", err)
			return
		}
		m.json(rw, screen)
	default:
		m.httpError(rw, "Invalid screens call:", fmt.Errorf("unknown method: %d", method))
		return
	}

	m.Info("Runtime call:", "method", "Screens."+screensMethodNames[method])

}
