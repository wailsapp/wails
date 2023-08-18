package application

import (
	"net/http"
)

func (m *MessageProcessor) processSystemMethod(method string, rw http.ResponseWriter, r *http.Request, window *WebviewWindow, params QueryParams) {

	switch method {
	case "IsDarkMode":
		m.json(rw, globalApplication.IsDarkMode())
	default:
		m.httpError(rw, "Unknown system method: %s", method)
	}

	m.Info("Runtime:", "method", "System."+method)

}
