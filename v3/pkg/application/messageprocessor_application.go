package application

import (
	"net/http"
)

const (
	ApplicationQuit = 0
	ApplicationHide = 1
	ApplicationShow = 2
)

var applicationMethodNames = map[int]string{
	ApplicationQuit: "Quit",
	ApplicationHide: "Hide",
	ApplicationShow: "Show",
}

func (m *MessageProcessor) processApplicationMethod(method int, rw http.ResponseWriter, r *http.Request, window *WebviewWindow, params QueryParams) {

	switch method {
	case ApplicationQuit:
		globalApplication.Quit()
		m.ok(rw)
	case ApplicationHide:
		globalApplication.Hide()
		m.ok(rw)
	case ApplicationShow:
		globalApplication.Show()
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown application method: %d", method)
	}

	m.Info("Runtime Call:", "method", "Application."+applicationMethodNames[method])

}
