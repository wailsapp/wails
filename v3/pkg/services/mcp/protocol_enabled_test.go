//go:build mcp

package mcp

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestServer() *mcpServer {
	m := &mcpServer{
		config:  Config{Host: "127.0.0.1", Port: 0, EvalTimeout: time.Second},
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		addr:    "127.0.0.1:0",
		pending: make(map[string]chan evalResult),
	}
	m.registerTools()
	return m
}

func postMCP(t *testing.T, m *mcpServer, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	m.handleMCP(recorder, request)
	return recorder
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response %q: %v", recorder.Body.String(), err)
	}
	return response
}

func TestInitialize(t *testing.T) {
	m := newTestServer()
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"0.0.1"}}}`)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	response := decodeResponse(t, recorder)
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

func TestInitializeEchoesOlderProtocolVersion(t *testing.T) {
	m := newTestServer()
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26"}}`)
	result := decodeResponse(t, recorder)["result"].(map[string]any)
	if result["protocolVersion"] != "2025-03-26" {
		t.Errorf("expected echoed protocol version, got %v", result["protocolVersion"])
	}
}

func TestToolsList(t *testing.T) {
	m := newTestServer()
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`)
	result := decodeResponse(t, recorder)["result"].(map[string]any)
	tools, ok := result["tools"].([]any)
	if !ok || len(tools) == 0 {
		t.Fatalf("expected tools array, got %v", result)
	}
	names := map[string]bool{}
	for _, raw := range tools {
		toolEntry := raw.(map[string]any)
		name := toolEntry["name"].(string)
		names[name] = true
		if toolEntry["description"] == "" {
			t.Errorf("tool %s has no description", name)
		}
		schema, ok := toolEntry["inputSchema"].(map[string]any)
		if !ok || schema["type"] != "object" {
			t.Errorf("tool %s has invalid inputSchema: %v", name, toolEntry["inputSchema"])
		}
	}
	for _, expected := range []string{
		"app_info", "windows_list", "window_control", "js_eval", "dom_html", "dom_query",
		"mouse_move", "mouse_click", "mouse_drag", "mouse_scroll",
		"keyboard_type", "keyboard_press", "call_bound_method",
		"emit_event", "wait_for_event", "screenshot_dom",
	} {
		if !names[expected] {
			t.Errorf("missing expected tool %q", expected)
		}
	}
}

func TestToolCall(t *testing.T) {
	m := newTestServer()
	m.tools = append(m.tools, &tool{
		Name:        "echo",
		Description: "test tool",
		Schema:      objectSchema(nil, map[string]any{}),
		Handler: func(args map[string]any) (any, error) {
			return map[string]any{"echo": args["value"]}, nil
		},
	})
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"echo","arguments":{"value":"hello"}}}`)
	result := decodeResponse(t, recorder)["result"].(map[string]any)
	if result["isError"] != false {
		t.Fatalf("expected isError=false, got %v", result)
	}
	content := result["content"].([]any)[0].(map[string]any)
	if !strings.Contains(content["text"].(string), "hello") {
		t.Errorf("unexpected content: %v", content)
	}
}

func TestToolCallError(t *testing.T) {
	m := newTestServer()
	m.tools = append(m.tools, &tool{
		Name:    "boom",
		Schema:  objectSchema(nil, map[string]any{}),
		Handler: func(args map[string]any) (any, error) { return nil, errors.New("it broke") },
	})
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"boom"}}`)
	result := decodeResponse(t, recorder)["result"].(map[string]any)
	if result["isError"] != true {
		t.Fatalf("expected isError=true, got %v", result)
	}
	content := result["content"].([]any)[0].(map[string]any)
	if content["text"] != "it broke" {
		t.Errorf("unexpected error text: %v", content["text"])
	}
}

func TestToolCallPanicIsCaught(t *testing.T) {
	m := newTestServer()
	m.tools = append(m.tools, &tool{
		Name:    "panic",
		Schema:  objectSchema(nil, map[string]any{}),
		Handler: func(args map[string]any) (any, error) { panic("kaboom") },
	})
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"panic"}}`)
	result := decodeResponse(t, recorder)["result"].(map[string]any)
	if result["isError"] != true {
		t.Fatalf("expected isError=true, got %v", result)
	}
	content := result["content"].([]any)[0].(map[string]any)
	if !strings.Contains(content["text"].(string), "kaboom") {
		t.Errorf("expected panic message, got %v", content["text"])
	}
}

func TestToolCallUnknownTool(t *testing.T) {
	m := newTestServer()
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"nope"}}`)
	response := decodeResponse(t, recorder)
	errorObject, ok := response["error"].(map[string]any)
	if !ok || errorObject["code"].(float64) != codeInvalidParams {
		t.Fatalf("expected invalid params error, got %v", response)
	}
}

func TestMethodNotFound(t *testing.T) {
	m := newTestServer()
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","id":7,"method":"resources/list"}`)
	response := decodeResponse(t, recorder)
	errorObject, ok := response["error"].(map[string]any)
	if !ok || errorObject["code"].(float64) != codeMethodNotFound {
		t.Fatalf("expected method not found error, got %v", response)
	}
}

func TestNotificationGets202(t *testing.T) {
	m := newTestServer()
	recorder := postMCP(t, m, `{"jsonrpc":"2.0","method":"notifications/initialized"}`)
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", recorder.Code)
	}
}

func TestBatchRequest(t *testing.T) {
	m := newTestServer()
	recorder := postMCP(t, m, `[{"jsonrpc":"2.0","id":1,"method":"ping"},{"jsonrpc":"2.0","id":2,"method":"ping"}]`)
	var responses []map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &responses); err != nil {
		t.Fatalf("expected array response, got %q", recorder.Body.String())
	}
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
}

func TestParseError(t *testing.T) {
	m := newTestServer()
	recorder := postMCP(t, m, `{not json`)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	response := decodeResponse(t, recorder)
	errorObject := response["error"].(map[string]any)
	if errorObject["code"].(float64) != codeParseError {
		t.Errorf("expected parse error code, got %v", errorObject)
	}
}

func TestOriginValidation(t *testing.T) {
	m := newTestServer()

	request := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"ping"}`))
	request.Header.Set("Origin", "http://evil.example.com")
	recorder := httptest.NewRecorder()
	m.handleMCP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Errorf("expected 403 for foreign origin, got %d", recorder.Code)
	}

	for _, origin := range []string{"http://localhost:3000", "http://127.0.0.1:9099", "http://wails.localhost"} {
		request = httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"ping"}`))
		request.Header.Set("Origin", origin)
		recorder = httptest.NewRecorder()
		m.handleMCP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Errorf("expected 200 for origin %s, got %d", origin, recorder.Code)
		}
	}
}

func TestUnsupportedHTTPMethods(t *testing.T) {
	m := newTestServer()

	request := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	recorder := httptest.NewRecorder()
	m.handleMCP(recorder, request)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for GET, got %d", recorder.Code)
	}

	request = httptest.NewRequest(http.MethodDelete, "/mcp", nil)
	recorder = httptest.NewRecorder()
	m.handleMCP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 for DELETE, got %d", recorder.Code)
	}
}

func TestEvalResultPlumbing(t *testing.T) {
	m := newTestServer()
	ch := make(chan evalResult, 1)
	m.pendingMu.Lock()
	m.pending["abc123"] = ch
	m.pendingMu.Unlock()

	request := httptest.NewRequest(http.MethodPost, "/eval-result", strings.NewReader(`{"id":"abc123","ok":true,"value":{"answer":42}}`))
	recorder := httptest.NewRecorder()
	m.handleEvalResult(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}

	select {
	case result := <-ch:
		if !result.Ok || !strings.Contains(string(result.Value), "42") {
			t.Errorf("unexpected result: %+v", result)
		}
	default:
		t.Fatal("result was not delivered")
	}

	m.pendingMu.Lock()
	_, stillPending := m.pending["abc123"]
	m.pendingMu.Unlock()
	if stillPending {
		t.Error("pending entry was not removed")
	}

	// Unknown IDs are ignored without error.
	request = httptest.NewRequest(http.MethodPost, "/eval-result", strings.NewReader(`{"id":"unknown","ok":true}`))
	recorder = httptest.NewRecorder()
	m.handleEvalResult(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Errorf("expected 204 for unknown id, got %d", recorder.Code)
	}
}

func TestCallbackURL(t *testing.T) {
	m := newTestServer()
	m.addr = "127.0.0.1:9099"
	if got := m.callbackURL(); got != "http://127.0.0.1:9099/eval-result" {
		t.Errorf("unexpected fallback callback URL: %s", got)
	}
	m.route = "/wails-mcp"
	if got := m.callbackURL(); got != "/wails-mcp/eval-result" {
		t.Errorf("unexpected same-origin callback URL: %s", got)
	}
	m.route = "/wails-mcp/"
	if got := m.callbackURL(); got != "/wails-mcp/eval-result" {
		t.Errorf("unexpected same-origin callback URL with trailing slash: %s", got)
	}
}

func TestInjectJSEmbedded(t *testing.T) {
	if !strings.Contains(injectJS, "window.__wailsMCP") {
		t.Fatal("inject.js does not define window.__wailsMCP")
	}
	for _, api := range []string{"run", "move", "click", "drag", "scroll", "typeText", "press", "query", "snapshot"} {
		if !strings.Contains(injectJS, api) {
			t.Errorf("inject.js is missing API %q", api)
		}
	}
}
