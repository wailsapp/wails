//go:build openharmony && cgo && !server

package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

// OHOS asset/runtime/message handling mirrors the android helpers: the ArkWeb
// shell forwards asset requests, runtime calls and window messages into this
// .so through the WailsBridge, and these functions service them.

var (
	openharmonyMessageProcessorOnce sync.Once
	openharmonyMessageProcessor     *MessageProcessor

	// openharmonyRuntimeReadyWindows tracks windows for which a synthetic
	// runtime-ready message has already been emitted (mirrors android).
	openharmonyRuntimeReadyWindows sync.Map
)

func openharmonyProcessor(app *App) *MessageProcessor {
	openharmonyMessageProcessorOnce.Do(func() {
		openharmonyMessageProcessor = NewMessageProcessor(app.Logger)
	})
	return openharmonyMessageProcessor
}

func openharmonyFirstWindowID(app *App) uint {
	if app == nil {
		return 0
	}
	windows := app.Window.GetAll()
	if len(windows) == 0 {
		return 0
	}
	return windows[0].ID()
}

type openharmonyRuntimeCallPayload struct {
	Object     *int            `json:"object"`
	Method     *int            `json:"method"`
	WindowName string          `json:"windowName"`
	Args       json.RawMessage `json:"args"`
	ClientID   string          `json:"clientId"`
}

func serveAssetForOpenharmony(app *App, method string, path string) ([]byte, error) {
	isRuntimeCall := strings.HasPrefix(path, "/wails/runtime")

	if !isRuntimeCall {
		if path == "" || path == "/" {
			path = "/index.html"
		}
	}

	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	if app.assets == nil {
		return nil, fmt.Errorf("asset server not initialized")
	}

	if method == "" {
		method = http.MethodGet
	}

	req, err := http.NewRequest(method, "https://wails.localhost"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// http.NewRequest leaves Body nil for client requests, but handlers reached
	// via ServeHTTP expect the server guarantee of a non-nil Body.
	req.Body = http.NoBody

	// Runtime calls need the window ID and name headers so the MessageProcessor
	// can route the call to the right window.
	if isRuntimeCall {
		windows := app.Window.GetAll()
		if len(windows) > 0 {
			windowID := windows[0].ID()
			req.Header.Set("x-wails-window-id", fmt.Sprintf("%d", windowID))
			req.Header.Set("x-wails-window-name", windows[0].Name())

			// A call to /wails/runtime proves the runtime is mounted, so treat the
			// first one as an implicit ready signal as a fallback.
			if _, alreadyReady := openharmonyRuntimeReadyWindows.LoadOrStore(windowID, true); !alreadyReady {
				windowMessageBuffer <- &windowMessage{
					windowId: windowID,
					message:  "wails:runtime:ready",
				}
			}
		}
	}

	recorder := httptest.NewRecorder()
	app.assets.ServeHTTP(recorder, req)

	result := recorder.Result()
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if isRuntimeCall && result.StatusCode != http.StatusOK {
		envelope, err := json.Marshal(map[string]any{"ok": false, "error": strings.TrimSpace(string(body))})
		if err != nil {
			return body, nil
		}
		return envelope, nil
	}

	if result.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("asset not found: %s", path)
	}

	if result.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("asset server error: status %d for %s", result.StatusCode, path)
	}

	return body, nil
}

func handleRuntimeCallForOpenharmony(app *App, payload string) string {
	fail := func(msg string) string {
		b, _ := json.Marshal(map[string]any{"ok": false, "error": msg})
		return string(b)
	}

	var call openharmonyRuntimeCallPayload
	if err := json.Unmarshal([]byte(payload), &call); err != nil {
		return fail("unable to parse runtime call: " + err.Error())
	}
	if call.Object == nil {
		return fail("missing object value")
	}
	if call.Method == nil {
		return fail("missing method value")
	}

	resp, err := openharmonyProcessor(app).HandleRuntimeCallWithIDs(context.Background(), &RuntimeRequest{
		Object:            *call.Object,
		Method:            *call.Method,
		Args:              &Args{call.Args},
		WebviewWindowName: call.WindowName,
		ClientID:          call.ClientID,
	})
	if err != nil {
		return fail(err.Error())
	}

	var envelope map[string]any
	if text, ok := resp.(string); ok {
		envelope = map[string]any{"ok": true, "text": text}
	} else {
		envelope = map[string]any{"ok": true, "data": resp}
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		return fail("unable to marshal response: " + err.Error())
	}
	return string(b)
}

func handleMessageForOpenharmony(app *App, message string) {
	if message == "" {
		return
	}

	openharmonyDebugLogf("[JS message] %s", message)

	msg := message
	var msgData map[string]interface{}
	if err := json.Unmarshal([]byte(message), &msgData); err == nil && msgData != nil {
		if name, ok := msgData["name"].(string); ok && name != "" {
			msg = name
		} else if name, ok := msgData["message"].(string); ok && name != "" {
			msg = name
		}
	}

	windowMessageBuffer <- &windowMessage{
		windowId: openharmonyFirstWindowID(app),
		message:  msg,
	}
}

func getMimeTypeForPath(path string) string {
	switch {
	// Server-sent event endpoints carry no file extension; EventSource rejects
	// the response unless the MIME type is text/event-stream.
	case endsWith(path, "/stream"), endsWith(path, "/events"):
		return "text/event-stream"
	case endsWith(path, ".html"), endsWith(path, ".htm"):
		return "text/html"
	case endsWith(path, ".js"), endsWith(path, ".mjs"):
		return "application/javascript"
	case endsWith(path, ".css"):
		return "text/css"
	case endsWith(path, ".json"), endsWith(path, ".map"):
		return "application/json"
	case endsWith(path, ".png"):
		return "image/png"
	case endsWith(path, ".jpg"), endsWith(path, ".jpeg"):
		return "image/jpeg"
	case endsWith(path, ".gif"):
		return "image/gif"
	case endsWith(path, ".svg"):
		return "image/svg+xml"
	case endsWith(path, ".ico"):
		return "image/x-icon"
	case endsWith(path, ".webp"):
		return "image/webp"
	case endsWith(path, ".woff"):
		return "font/woff"
	case endsWith(path, ".woff2"):
		return "font/woff2"
	case endsWith(path, ".ttf"):
		return "font/ttf"
	case endsWith(path, ".otf"):
		return "font/otf"
	case endsWith(path, ".wasm"):
		return "application/wasm"
	case endsWith(path, ".mp3"):
		return "audio/mpeg"
	case endsWith(path, ".mp4"):
		return "video/mp4"
	case endsWith(path, ".webm"):
		return "video/webm"
	case endsWith(path, ".xml"):
		return "application/xml"
	case endsWith(path, ".txt"):
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

func endsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

var _ = slog.Debug
