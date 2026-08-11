//go:build mcp && !ios && !android

package application

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// protocolVersion is the latest MCP protocol revision this server supports.
// Older revisions are accepted; the subset we implement (initialize, ping,
// tools/list, tools/call) is identical across them.
const mcpProtocolVersion = "2025-06-18"

const mcpServerVersion = "1.0.0"

// JSON-RPC 2.0 error codes.
const (
	mcpCodeParseError     = -32700
	mcpCodeInvalidRequest = -32600
	mcpCodeMethodNotFound = -32601
	mcpCodeInvalidParams  = -32602
	mcpCodeInternalError  = -32603
)

type mcpJSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type mcpJSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type mcpJSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *mcpJSONRPCError `json:"error,omitempty"`
}

func mcpErrorResponse(id json.RawMessage, code int, format string, args ...any) *mcpJSONRPCResponse {
	return &mcpJSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &mcpJSONRPCError{Code: code, Message: fmt.Sprintf(format, args...)},
	}
}

func mcpResultResponse(id json.RawMessage, result any) *mcpJSONRPCResponse {
	return &mcpJSONRPCResponse{JSONRPC: "2.0", ID: id, Result: result}
}

// mcpOriginAllowed guards against DNS-rebinding attacks as required by the MCP
// streamable HTTP spec. Requests with no Origin header (non-browser clients)
// are always allowed; browser contexts must originate from localhost.
func mcpOriginAllowed(r *http.Request) bool {
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

func mcpSetCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Mcp-Session-Id, MCP-Protocol-Version, Last-Event-ID")
	w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id")
}

// handleMCP implements the server side of the MCP streamable HTTP transport.
// Responses are always plain JSON (the spec allows this in place of SSE).
func (m *mcpServer) handleMCP(w http.ResponseWriter, r *http.Request) {
	mcpSetCORSHeaders(w)
	if !mcpOriginAllowed(r) {
		http.Error(w, "forbidden origin", http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
		return
	case http.MethodGet:
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

	requests, batch, parseErr := mcpParseMessages(body)
	if parseErr != nil {
		mcpWriteJSON(w, http.StatusBadRequest, mcpErrorResponse(nil, mcpCodeParseError, "parse error: %v", parseErr))
		return
	}
	if batch && len(requests) == 0 {
		mcpWriteJSON(w, http.StatusBadRequest, mcpErrorResponse(nil, mcpCodeInvalidRequest, "invalid request: empty batch"))
		return
	}

	var responses []*mcpJSONRPCResponse
	for _, req := range requests {
		if response := m.handleMessage(req); response != nil {
			responses = append(responses, response)
		}
	}

	if len(responses) == 0 {
		// Notifications only — acknowledged but no body.
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if batch {
		mcpWriteJSON(w, http.StatusOK, responses)
		return
	}
	mcpWriteJSON(w, http.StatusOK, responses[0])
}

func mcpParseMessages(body []byte) (requests []*mcpJSONRPCRequest, batch bool, err error) {
	trimmed := strings.TrimLeftFunc(string(body), func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\r' || r == '\n'
	})
	if strings.HasPrefix(trimmed, "[") {
		var reqs []*mcpJSONRPCRequest
		if err := json.Unmarshal(body, &reqs); err != nil {
			return nil, true, err
		}
		return reqs, true, nil
	}
	var req mcpJSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, false, err
	}
	return []*mcpJSONRPCRequest{&req}, false, nil
}

func mcpWriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

// handleMessage processes a single JSON-RPC message. Returns nil for
// notifications, which receive no response.
func (m *mcpServer) handleMessage(req *mcpJSONRPCRequest) *mcpJSONRPCResponse {
	isNotification := len(req.ID) == 0 || string(req.ID) == "null"
	if isNotification || strings.HasPrefix(req.Method, "notifications/") {
		return nil
	}

	switch req.Method {
	case "initialize":
		var params struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		_ = json.Unmarshal(req.Params, &params)
		version := mcpProtocolVersion
		switch params.ProtocolVersion {
		case "2024-11-05", "2025-03-26":
			version = params.ProtocolVersion
		}
		return mcpResultResponse(req.ID, map[string]any{
			"protocolVersion": version,
			"capabilities": map[string]any{
				"tools": map[string]any{"listChanged": false},
			},
			"serverInfo": map[string]any{
				"name":    "wails-mcp",
				"title":   "Wails Application Control",
				"version": mcpServerVersion,
			},
			"instructions": "This MCP server controls a live Wails desktop application. " +
				"Use the tools to inspect windows, query the DOM, simulate user input " +
				"(mouse input is shown with an animated cursor inside the app), evaluate " +
				"JavaScript and call bound Go methods. Coordinates are CSS pixels relative " +
				"to the window's viewport.",
		})

	case "ping":
		return mcpResultResponse(req.ID, map[string]any{})

	case "tools/list":
		tools := make([]map[string]any, 0, len(m.tools))
		for _, t := range m.tools {
			tools = append(tools, map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"inputSchema": t.Schema,
			})
		}
		return mcpResultResponse(req.ID, map[string]any{"tools": tools})

	case "tools/call":
		return m.handleToolCall(req)

	default:
		return mcpErrorResponse(req.ID, mcpCodeMethodNotFound, "method not found: %s", req.Method)
	}
}

func (m *mcpServer) handleToolCall(req *mcpJSONRPCRequest) (response *mcpJSONRPCResponse) {
	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return mcpErrorResponse(req.ID, mcpCodeInvalidParams, "invalid params: %v", err)
	}

	var selected *mcpTool
	for _, t := range m.tools {
		if t.Name == params.Name {
			selected = t
			break
		}
	}
	if selected == nil {
		return mcpErrorResponse(req.ID, mcpCodeInvalidParams, "unknown tool: %s", params.Name)
	}
	if params.Arguments == nil {
		params.Arguments = map[string]any{}
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			m.logger.Error("mcp: tool panicked", "tool", params.Name, "panic", recovered)
			response = mcpResultResponse(req.ID, mcpToolError(fmt.Sprintf("tool %s panicked: %v", params.Name, recovered)))
		}
	}()

	result, err := selected.Handler(params.Arguments)
	if err != nil {
		return mcpResultResponse(req.ID, mcpToolError(err.Error()))
	}
	return mcpResultResponse(req.ID, mcpToolResult(result))
}

// mcpToolResult converts a tool's return value into an MCP CallToolResult.
func mcpToolResult(value any) map[string]any {
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

func mcpToolError(message string) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": message}},
		"isError": true,
	}
}
