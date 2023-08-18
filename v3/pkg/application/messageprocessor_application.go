package application

import (
	"net/http"
)

func (m *MessageProcessor) processApplicationMethod(method string, rw http.ResponseWriter, r *http.Request, window *WebviewWindow, params QueryParams) {

	switch method {
	case "Quit":
		globalApplication.Quit()
		m.ok(rw)
	case "Hide":
		globalApplication.Hide()
		m.ok(rw)
	case "Show":
		globalApplication.Show()
		m.ok(rw)
	case "IsDarkMode":
		m.json(rw, globalApplication.IsDarkMode())
	default:
		m.httpError(rw, "Unknown application method: %s", method)
	}

	m.Info("Runtime:", "method", "Application."+method)

}
