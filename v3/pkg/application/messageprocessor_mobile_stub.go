//go:build !ios && !android

package application

import (
	"fmt"
	"net/http"
)

// processIOSMethod is a stub for non-mobile platforms
func (m *MessageProcessor) processIOSMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	m.httpError(rw, "iOS methods not available:", fmt.Errorf("unknown method: %d", method))
}

// processAndroidMethod is a stub for non-mobile platforms
func (m *MessageProcessor) processAndroidMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	m.httpError(rw, "Android methods not available:", fmt.Errorf("unknown method: %d", method))
}
