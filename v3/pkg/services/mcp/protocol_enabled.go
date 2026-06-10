//go:build mcp

package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// protocolVersion is the latest MCP protocol revision this server knows.
// Older revisions are accepted as-is since the subset we implement
// (initialize, ping, tools/list, tools/call) is identical across them.
const protocolVersion = "2025-06-18"

const serverVersion = "1.0.0"

// JSON-RPC 2.0 error codes.
const (
	codeParseError     = -32700
	codeInvalidRequest = -32600
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
	codeInternalError  = -32603
)

type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

func errorResponse(id json.RawMessage, code int, format string, args ...any) *jsonrpcResponse {
	return &jsonrpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &jsonrpcError{Code: code, Message: fmt.Sprintf(format, args...)},
	}
}

func resultResponse(id json.RawMessage, result any) *jsonrpcResponse {
	return &jsonrpcResponse{JSONRPC: "2.0", ID: id, Result: result}
}

// originAllowed guards against DNS-rebinding attacks as required by the MCP
// streamable HTTP spec. Requests from non-browser clients carry no Origin
// header and are allowed; browser contexts must be local.
func originAllowed(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" || origin == "null" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := u.Hostname()
	return host == "localhost" || host == "127.0.0.1" || host == "::1" ||
		strings.HasSuffix(host, ".localhost")
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Mcp-Session-Id, MCP-Protocol-Version, Last-Event-ID")
	w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id")
}

// handleMCP implements the server side of the MCP streamable HTTP transport.
// Responses are always plain JSON (the spec allows this in place of SSE).
func (m *mcpServer) handleMCP(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if !originAllowed(r) {
		http.Error(w, "forbidden origin", http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
		return
	case http.MethodGet:
		// No server-initiated streams.
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	case http.MethodDelete:
		// Stateless server: session termination is a no-op.
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodPost:
		// Handled below.
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 16*1024*1024))
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	requests, batch, parseErr := parseMessages(body)
	if parseErr != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse(nil, codeParseError, "parse error: %v", parseErr))
		return
	}

	var responses []*jsonrpcResponse
	for _, req := range requests {
		if response := m.handleMessage(req); response != nil {
			responses = append(responses, response)
		}
	}

	if len(responses) == 0 {
		// Notifications and responses only.
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if batch {
		writeJSON(w, http.StatusOK, responses)
		return
	}
	writeJSON(w, http.StatusOK, responses[0])
}

func parseMessages(body []byte) (requests []*jsonrpcRequest, batch bool, err error) {
	trimmed := strings.TrimLeftFunc(string(body), func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\r' || r == '\n'
	})
	if strings.HasPrefix(trimmed, "[") {
		var reqs []*jsonrpcRequest
		if err := json.Unmarshal(body, &reqs); err != nil {
			return nil, true, err
		}
		return reqs, true, nil
	}
	var req jsonrpcRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, false, err
	}
	return []*jsonrpcRequest{&req}, false, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

// handleMessage processes a single JSON-RPC message. It returns nil for
// notifications, which receive no response.
func (m *mcpServer) handleMessage(req *jsonrpcRequest) *jsonrpcResponse {
	isNotification := len(req.ID) == 0 || string(req.ID) == "null"

	switch req.Method {
	case "initialize":
		var params struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		_ = json.Unmarshal(req.Params, &params)
		version := protocolVersion
		// Echo back known earlier revisions for maximum compatibility.
		switch params.ProtocolVersion {
		case "2024-11-05", "2025-03-26":
			version = params.ProtocolVersion
		}
		return resultResponse(req.ID, map[string]any{
			"protocolVersion": version,
			"capabilities": map[string]any{
				"tools": map[string]any{"listChanged": false},
			},
			"serverInfo": map[string]any{
				"name":    "wails-mcp",
				"title":   "Wails Application Control",
				"version": serverVersion,
			},
			"instructions": "This MCP server controls a live Wails desktop application. " +
				"Use the tools to inspect windows, query the DOM, simulate user input " +
				"(mouse input is shown with an animated cursor inside the app), evaluate " +
				"JavaScript and call bound Go methods. Coordinates are CSS pixels relative " +
				"to the window's viewport.",
		})

	case "ping":
		return resultResponse(req.ID, map[string]any{})

	case "tools/list":
		tools := make([]map[string]any, 0, len(m.tools))
		for _, t := range m.tools {
			tools = append(tools, map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"inputSchema": t.Schema,
			})
		}
		return resultResponse(req.ID, map[string]any{"tools": tools})

	case "tools/call":
		if isNotification {
			return nil
		}
		return m.handleToolCall(req)

	default:
		if strings.HasPrefix(req.Method, "notifications/") || isNotification {
			return nil
		}
		return errorResponse(req.ID, codeMethodNotFound, "method not found: %s", req.Method)
	}
}

func (m *mcpServer) handleToolCall(req *jsonrpcRequest) (response *jsonrpcResponse) {
	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return errorResponse(req.ID, codeInvalidParams, "invalid params: %v", err)
	}

	var selected *tool
	for _, t := range m.tools {
		if t.Name == params.Name {
			selected = t
			break
		}
	}
	if selected == nil {
		return errorResponse(req.ID, codeInvalidParams, "unknown tool: %s", params.Name)
	}
	if params.Arguments == nil {
		params.Arguments = map[string]any{}
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			m.logger.Error("mcp: tool panicked", "tool", params.Name, "panic", recovered)
			response = resultResponse(req.ID, toolError(fmt.Sprintf("tool %s panicked: %v", params.Name, recovered)))
		}
	}()

	result, err := selected.Handler(params.Arguments)
	if err != nil {
		return resultResponse(req.ID, toolError(err.Error()))
	}
	return resultResponse(req.ID, toolResult(result))
}

// toolResult converts a tool's return value into an MCP CallToolResult.
func toolResult(value any) map[string]any {
	var text string
	switch v := value.(type) {
	case nil:
		text = "ok"
	case string:
		text = v
	default:
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			text = fmt.Sprintf("%v", v)
		} else {
			text = string(data)
		}
	}
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": text}},
		"isError": false,
	}
}

func toolError(message string) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": message}},
		"isError": true,
	}
}
