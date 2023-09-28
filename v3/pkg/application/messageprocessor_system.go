package application

import (
	"net/http"
)

const (
	SystemIsDarkMode = 0
)

var systemMethodNames = map[int]string{
	SystemIsDarkMode: "IsDarkMode",
}

func (m *MessageProcessor) processSystemMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {

	switch method {
	case SystemIsDarkMode:
		m.json(rw, globalApplication.IsDarkMode())
	default:
		m.httpError(rw, "Unknown system method: %d", method)
	}

	m.Info("Runtime:", "method", "System."+systemMethodNames[method])

}
