package application

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	httpFetch = 0
)

func (m *MessageProcessor) processHTTPMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	switch method {
	case httpFetch:
		m.httpFetch(rw, r, window)
	default:
		m.httpError(rw, "Unknown HTTP method: %d", method)
	}
}

func (m *MessageProcessor) httpFetch(rw http.ResponseWriter, r *http.Request, window Window) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		m.httpError(rw, "Failed to read request body: %v", err)
		return
	}

	// Parse the HTTP request
	request, err := parseHTTPRequest(string(body))
	if err != nil {
		m.httpError(rw, "Invalid HTTP request: %v", err)
		return
	}

	// Perform the HTTP request
	response := PerformHTTPRequest(*request)

	// Marshal response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		m.httpError(rw, "Failed to marshal response: %v", err)
		return
	}

	// Send response
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(responseJSON)
}