//go:build mcp && !ios && !android

package application

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestMCPServer() *mcpServer {
	m := &mcpServer{
		logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		evalTimeout: time.Second,
		addr:        "127.0.0.1:0",
		pending:     make(map[string]chan mcpEvalResult),
	}
	m.registerTools()
	return m
}

func postMCP(t *testing.T, m *mcpServer, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	m.handleMCP(rec, req)
	return rec
}

func decodeMCPResponse(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var response map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response %q: %v", rec.Body.String(), err)
	}
	return response
}

func TestMCPInitialize(t *testing.T) {
	m := newTestMCPServer()
	rec := postMCP(t, m, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"0.0.1"}}}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	response := decodeMCPResponse(t, rec)
	result, ok := response["result"].(map[string]any)
	if !ok {
		t.Fatalf("expected result object, got %v", response)
	}
	if result["protocolVersion"] != "2025-06-18" {
		t.Errorf("unexpected protocolVersion: %v", result["protocolVersion"])
	}
	serverInfo, _ := result["serverInfo"].(map[string]any)
	if serverInfo["name"] != "wails-mcp" {
		t.Errorf("unexpected serverInfo: %v", serverInfo)
	}
}

func TestMCPInitializeEchoesOlderProtocol(t *testing.T) {
	m := newTestMCPServer()
	rec := postMCP(t, m, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26"}}`)
	result := decodeMCPResponse(t, rec)["result"].(map[string]any)
	if result["protocolVersion"] != "2025-03-26" {
		t.Errorf("expected echoed protocol version, got %v", result["protocolVersion"])
	}
}

func TestMCPToolsList(t *testing.T) {
	m := newTestMCPServer()
	rec := postMCP(t, m, `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`)
	result := decodeMCPResponse(t, rec)["result"].(map[string]any)
	tools, ok := result["tools"].([]any)
	if !ok || len(tools) == 0 {
		t.Fatalf("expected non-empty tools list, got %v", result)
	}
	// Verify the registry exposes exactly the expected 16 tools.
	names := make(map[string]bool)
	for _, tool := range tools {
		toolMap, _ := tool.(map[string]any)
		if name, ok := toolMap["name"].(string); ok {
			if names[name] {
				t.Errorf("tool %q duplicated in tools/list", name)
			}
			names[name] = true
		}
	}
	expected := []string{
		"app_info", "windows_list", "window_control", "js_eval",
		"dom_html", "dom_query",
		"mouse_move", "mouse_click", "mouse_drag", "mouse_scroll",
		"keyboard_type", "keyboard_press",
		"call_bound_method", "emit_event", "wait_for_event", "screenshot_dom",
	}
	if len(tools) != len(expected) {
		t.Errorf("expected exactly %d tools, got %d", len(expected), len(tools))
	}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("tool %q missing from tools/list", name)
		}
	}
}

func TestMCPPing(t *testing.T) {
	m := newTestMCPServer()
	rec := postMCP(t, m, `{"jsonrpc":"2.0","id":3,"method":"ping"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	response := decodeMCPResponse(t, rec)
	if response["result"] == nil {
		t.Error("expected non-nil result for ping")
	}
}

func TestMCPMethodNotFound(t *testing.T) {
	m := newTestMCPServer()
	rec := postMCP(t, m, `{"jsonrpc":"2.0","id":4,"method":"unknown/method"}`)
	response := decodeMCPResponse(t, rec)
	errObj, ok := response["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object, got %v", response)
	}
	if errObj["code"].(float64) != mcpCodeMethodNotFound {
		t.Errorf("expected method-not-found code, got %v", errObj["code"])
	}
}

func TestMCPNotificationNoResponse(t *testing.T) {
	m := newTestMCPServer()
	// Notifications have no id field — server returns 202 Accepted with no body.
	rec := postMCP(t, m, `{"jsonrpc":"2.0","method":"notifications/initialized"}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 for notification, got %d", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body for notification, got %q", rec.Body.String())
	}
}

func TestMCPBatchRequest(t *testing.T) {
	m := newTestMCPServer()
	batch := `[
		{"jsonrpc":"2.0","id":1,"method":"ping"},
		{"jsonrpc":"2.0","id":2,"method":"ping"}
	]`
	rec := postMCP(t, m, batch)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var responses []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &responses); err != nil {
		t.Fatalf("failed to decode batch response: %v", err)
	}
	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}
}

func TestMCPUnknownTool(t *testing.T) {
	m := newTestMCPServer()
	rec := postMCP(t, m, `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"nonexistent"}}`)
	response := decodeMCPResponse(t, rec)
	// Unknown tools produce a JSON-RPC error response (invalid params).
	errObj, ok := response["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object for unknown tool, got %v", response)
	}
	if errObj["code"].(float64) != mcpCodeInvalidParams {
		t.Errorf("expected invalid-params code, got %v", errObj["code"])
	}
}

func TestMCPPanicRecovery(t *testing.T) {
	m := newTestMCPServer()
	m.tools = append(m.tools, &mcpTool{
		Name:   "panic_tool",
		Schema: mcpObjectSchema(nil, nil),
		Handler: func(args map[string]any) (any, error) {
			panic("intentional panic")
		},
	})
	rec := postMCP(t, m, `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"panic_tool"}}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 after panic recovery, got %d", rec.Code)
	}
	response := decodeMCPResponse(t, rec)
	result, _ := response["result"].(map[string]any)
	if result["isError"] != true {
		t.Errorf("expected isError=true after panic, got %v", result)
	}
}

func TestMCPOriginAllowed(t *testing.T) {
	cases := []struct {
		origin  string
		allowed bool
	}{
		{"", true},
		{"null", true},
		{"http://localhost:3000", true},
		{"http://127.0.0.1:8080", true},
		{"http://[::1]:0", true},
		{"http://evil.com", false},
		{"https://attacker.com", false},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}
		if got := mcpOriginAllowed(req); got != tc.allowed {
			t.Errorf("origin %q: expected allowed=%v, got %v", tc.origin, tc.allowed, got)
		}
	}
}

func TestMCPForbiddenOrigin(t *testing.T) {
	m := newTestMCPServer()
	req := httptest.NewRequest(http.MethodPost, "/mcp",
		strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"ping"}`))
	req.Header.Set("Origin", "http://evil.example.com")
	rec := httptest.NewRecorder()
	m.handleMCP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 for bad origin, got %d", rec.Code)
	}
}

func TestMCPEvalResultDelivery(t *testing.T) {
	m := newTestMCPServer()
	id := "test-eval-id-12345"
	ch := make(chan mcpEvalResult, 1)
	m.pendingMu.Lock()
	m.pending[id] = ch
	m.pendingMu.Unlock()

	payload := `{"id":"test-eval-id-12345","ok":true,"value":42}`
	req := httptest.NewRequest(http.MethodPost, "/eval-result", strings.NewReader(payload))
	rec := httptest.NewRecorder()
	m.handleEvalResult(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	select {
	case result := <-ch:
		if !result.Ok {
			t.Errorf("expected ok=true, got %v", result)
		}
		var val float64
		if err := json.Unmarshal(result.Value, &val); err != nil || val != 42 {
			t.Errorf("unexpected value: %v (err=%v)", result.Value, err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for eval result delivery")
	}
}

func TestMCPEnvVarDefaults(t *testing.T) {
	// Verify defaults match constants.
	if mcpDefaultHost != "127.0.0.1" {
		t.Errorf("unexpected default host: %s", mcpDefaultHost)
	}
	if mcpDefaultPort != 9099 {
		t.Errorf("unexpected default port: %d", mcpDefaultPort)
	}
	if mcpDefaultTimeout != 30*time.Second {
		t.Errorf("unexpected default timeout: %v", mcpDefaultTimeout)
	}
}
