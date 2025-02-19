package application

import (
	"fmt"
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
		m.httpError(rw, "Invalid system call:", fmt.Errorf("unknown method: %d", method))
		return
	}

	m.Info("Runtime call:", "method", "System."+systemMethodNames[method])
}
