package application

import (
	"net/http"
)

const (
	SystemIsDarkMode = 0
	Environment      = 1
)

var systemMethodNames = map[int]string{
	SystemIsDarkMode: "IsDarkMode",
	Environment:      "Environment",
}

func (m *MessageProcessor) processSystemMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {

	switch method {
	case SystemIsDarkMode:
		m.json(rw, globalApplication.IsDarkMode())
	case Environment:
		m.json(rw, globalApplication.Environment())
	default:
		m.httpError(rw, "Unknown system method: %d", method)
	}

	m.Info("Runtime Call:", "method", "System."+systemMethodNames[method])

}
