package application

import (
	"net/http"
)

func (m *MessageProcessor) processScreensMethod(method string, rw http.ResponseWriter, _ *http.Request, _ *WebviewWindow, _ QueryParams) {

	switch method {
	case "GetAll":
		screens, err := globalApplication.GetScreens()
		if err != nil {
			m.Error("GetAll: %s", err.Error())
			return
		}
		m.json(rw, screens)
	case "GetPrimary":
		screen, err := globalApplication.GetPrimaryScreen()
		if err != nil {
			m.Error("GetPrimary: %s", err.Error())
			return
		}
		m.json(rw, screen)
	case "GetCurrent":
		screen, err := globalApplication.CurrentWindow().GetScreen()
		if err != nil {
			m.Error("GetCurrent: %s", err.Error())
			return
		}
		m.json(rw, screen)
	default:
		m.httpError(rw, "Unknown clipboard method: %s", method)
	}

}
