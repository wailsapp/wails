//go:build mcp && !ios && !android

package application

import (
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

// mcpInjectJS is the in-page support library: animated cursor overlay, input
// simulation and result callback harness. Idempotent, prepended to every eval.
//
//go:embed mcp_inject.js
var mcpInjectJS string

// mcpEvalResult is the payload posted back by the page when an evaluation
// finishes.
type mcpEvalResult struct {
	ID    string          `json:"id"`
	Ok    bool            `json:"ok"`
	Value json.RawMessage `json:"value"`
	Error string          `json:"error"`
}

// callbackURL returns the URL the in-page script POSTs evaluation results to.
// Wildcard bind addresses (0.0.0.0, ::) are normalised to localhost so the
// URL is reachable from within the webview.
func (m *mcpServer) callbackURL() string {
	host, port, err := net.SplitHostPort(m.addr)
	if err != nil {
		return fmt.Sprintf("http://%s/eval-result", m.addr)
	}
	if host == "0.0.0.0" || host == "::" || host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("http://%s/eval-result", net.JoinHostPort(host, port))
}

// eval runs JavaScript in the window and waits for its result. The body is
// wrapped in an async function receiving the support library as `mcp`; use
// `return` to produce a value.
func (m *mcpServer) eval(window Window, body string, timeout time.Duration) (json.RawMessage, error) {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return nil, fmt.Errorf("failed to generate evaluation id: %w", err)
	}
	id := hex.EncodeToString(idBytes)

	ch := make(chan mcpEvalResult, 1)
	m.pendingMu.Lock()
	m.pending[id] = ch
	m.pendingMu.Unlock()
	defer func() {
		m.pendingMu.Lock()
		delete(m.pending, id)
		m.pendingMu.Unlock()
	}()

	script := mcpInjectJS + "\n" + fmt.Sprintf(
		"window.__wailsMCP.run(%s, %s, %s, async (mcp) => {\n%s\n});",
		strconv.Quote(id),
		strconv.Quote(m.callbackURL()),
		strconv.FormatBool(!m.hideCursor),
		body,
	)
	window.ExecJS(script)

	select {
	case result := <-ch:
		if !result.Ok {
			return nil, fmt.Errorf("javascript error: %s", result.Error)
		}
		return result.Value, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timed out after %s waiting for the page to respond; "+
			"the window may be reloading, busy or showing a native dialog", timeout)
	}
}

// evalInto runs JavaScript and unmarshals the result into out.
func (m *mcpServer) evalInto(window Window, body string, timeout time.Duration, out any) error {
	value, err := m.eval(window, body, timeout)
	if err != nil {
		return err
	}
	if len(value) == 0 || string(value) == "null" {
		return nil
	}
	return json.Unmarshal(value, out)
}

// handleEvalResult receives evaluation results posted by the in-page script.
func (m *mcpServer) handleEvalResult(w http.ResponseWriter, r *http.Request) {
	mcpSetCORSHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
		return
	case http.MethodPost:
		// Handled below.
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var result mcpEvalResult
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64*1024*1024)).Decode(&result); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	m.pendingMu.Lock()
	ch, ok := m.pending[result.ID]
	if ok {
		delete(m.pending, result.ID)
	}
	m.pendingMu.Unlock()

	if ok {
		ch <- result
	}
	w.WriteHeader(http.StatusNoContent)
}
