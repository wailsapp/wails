package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPRequest represents an HTTP request from the frontend
type HTTPRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Timeout int               `json:"timeout,omitempty"` // timeout in seconds
}

// HTTPResponse represents an HTTP response to send back to the frontend
type HTTPResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Error      string            `json:"error,omitempty"`
}

// PerformHTTPRequest performs an HTTP request on behalf of the frontend
func PerformHTTPRequest(request HTTPRequest) HTTPResponse {
	// Validate method
	method := strings.ToUpper(request.Method)
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" && method != "HEAD" && method != "OPTIONS" {
		return HTTPResponse{
			StatusCode: 0,
			Error:      fmt.Sprintf("Invalid HTTP method: %s", request.Method),
		}
	}

	// Create HTTP client with timeout
	timeout := 30 * time.Second
	if request.Timeout > 0 {
		timeout = time.Duration(request.Timeout) * time.Second
	}
	client := &http.Client{
		Timeout: timeout,
	}

	// Create request
	var body io.Reader
	if request.Body != "" {
		body = bytes.NewBufferString(request.Body)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, request.URL, body)
	if err != nil {
		return HTTPResponse{
			StatusCode: 0,
			Error:      fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	// Set headers
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	// Set default User-Agent if not provided
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Wails/3.0")
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return HTTPResponse{
			StatusCode: 0,
			Error:      fmt.Sprintf("Request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTTPResponse{
			StatusCode: resp.StatusCode,
			Error:      fmt.Sprintf("Failed to read response body: %v", err),
		}
	}

	// Build response headers
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	return HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    responseHeaders,
		Body:       string(bodyBytes),
	}
}

// parseHTTPRequest parses the HTTP request from JSON
func parseHTTPRequest(data string) (*HTTPRequest, error) {
	var request HTTPRequest
	err := json.Unmarshal([]byte(data), &request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTTP request: %v", err)
	}

	// Validate required fields
	if request.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}
	if request.Method == "" {
		request.Method = "GET"
	}

	return &request, nil
}