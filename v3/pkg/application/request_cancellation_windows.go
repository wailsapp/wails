//go:build windows && !server

package application

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

//go:embed request_keepalive_windows.js
var requestKeepaliveScript string

const keepaliveBinding = "__wailsAssetKeepalive"
const keepaliveFragment = "&__wails_keepalive="

type requestCancellationDevTools interface {
	CallDevTools(string, string, func(string, error)) error
	CallDevToolsForSession(string, string, string, func(string, error)) error
	OnDevTools(string, func(string, string, error)) (func() error, error)
}

// windowsRequestCancellation owns the CDP sessions for one WebView.
// Fetch provides a network ID before WebResourceRequested; Network reports aborts
// while the HTTP handler is still running. Requests are matched by identity, never URL.
// Except for tracker operations, all access is confined to the UI thread.
type windowsRequestCancellation struct {
	chromium        requestCancellationDevTools
	tracker         *webview.RequestTracker
	unsubscribe     []func() error
	closed          bool
	owners          map[string]string
	scriptID        string
	keepaliveOwners map[string]keepaliveOwner
}

func (w *windowsWebviewWindow) startRequestCancellation(ready func()) {
	if globalApplication.assets == nil {
		ready()
		return
	}
	c := &windowsRequestCancellation{chromium: w.chromium, tracker: webview.NewRequestTracker(globalApplication.Context())}
	w.requestCancellation = c
	done := func(err error) {
		if c.closed {
			return
		}
		if err != nil {
			globalApplication.error("Windows request cancellation unavailable: %v", err)
			c.close()
		}
		if !w.parent.isDestroyed() {
			ready()
		}
	}
	if err := c.subscribe(); err != nil {
		done(err)
		return
	}
	c.enableSession("", done)
}

type keepaliveOwner struct {
	session string
	context int64
}

type cancellationCommand struct {
	method string
	params any
}

func (c *windowsRequestCancellation) enableSession(session string, done func(error)) {
	commands := []cancellationCommand{
		{"Network.enable", struct{}{}},
		{"Runtime.enable", struct{}{}},
		{"Runtime.addBinding", map[string]string{"name": keepaliveBinding}},
	}
	if session == "" {
		commands = append(commands,
			cancellationCommand{"Page.enable", struct{}{}},
			cancellationCommand{"Page.addScriptToEvaluateOnNewDocument", map[string]string{"source": requestKeepaliveScript}},
			cancellationCommand{"Fetch.enable", map[string]any{"patterns": []map[string]string{
				{"urlPattern": "http://wails.localhost/*", "requestStage": "Request"},
				{"urlPattern": "http://wails.localhost:*/*", "requestStage": "Request"},
			}}},
		)
	} else {
		// WebView2 intercepts worker requests on the root Fetch session. Only the
		// worker's Network domain must be enabled to observe its terminal events.
		commands = append(commands, cancellationCommand{"Runtime.evaluate", map[string]string{"expression": requestKeepaliveScript}})
	}
	commands = append(commands, cancellationCommand{"Target.setAutoAttach", map[string]any{"autoAttach": true, "waitForDebuggerOnStart": true, "flatten": true}})
	c.configure(session, commands, done)
}

func (c *windowsRequestCancellation) configure(session string, commands []cancellationCommand, done func(error)) {
	if len(commands) == 0 {
		done(nil)
		return
	}
	c.call(session, commands[0].method, commands[0].params, func(err error) {
		if err != nil {
			done(err)
			return
		}
		c.configure(session, commands[1:], done)
	})
}

func (c *windowsRequestCancellation) subscribe() error {
	for _, event := range []string{"Network.requestWillBeSent", "Runtime.bindingCalled", "Runtime.executionContextDestroyed", "Runtime.executionContextsCleared", "Fetch.requestPaused", "Network.loadingFailed", "Network.loadingFinished", "Target.attachedToTarget", "Target.detachedFromTarget"} {
		unsubscribe, err := c.chromium.OnDevTools(event, func(session, data string, err error) {
			if c.closed {
				return
			}
			if err == nil {
				err = c.event(session, event, data)
			}
			if err != nil {
				globalApplication.error("Windows request cancellation (%s): %v", event, err)
			}
		})
		if err != nil {
			return err
		}
		c.unsubscribe = append(c.unsubscribe, unsubscribe)
	}
	return nil
}

func (c *windowsRequestCancellation) call(session, method string, params any, done func(error)) {
	if c.closed {
		return
	}
	data, err := json.Marshal(params)
	if err != nil {
		done(err)
		return
	}
	err = c.chromium.CallDevToolsForSession(session, method, string(data), func(result string, err error) {
		if method == "Page.addScriptToEvaluateOnNewDocument" && err == nil {
			var script struct{ Identifier string }
			err = json.Unmarshal([]byte(result), &script)
			c.scriptID = script.Identifier
		}
		if !c.closed {
			done(err)
		}
	})
	if err != nil {
		done(err)
	}
}

type pausedWebRequest struct {
	RequestID string `json:"requestId"`
	NetworkID string `json:"networkId"`
	Request   struct {
		Headers     map[string]string `json:"headers"`
		URLFragment string            `json:"urlFragment"`
	} `json:"request"`
}

type requestHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *windowsRequestCancellation) event(session, event, data string) error {
	switch event {
	case "Runtime.bindingCalled":
		return c.keepaliveEvent(session, data)
	case "Target.attachedToTarget", "Target.detachedFromTarget":
		return c.targetEvent(event, data)
	case "Runtime.executionContextDestroyed":
		var context struct {
			ExecutionContextID int64 `json:"executionContextId"`
		}
		if err := json.Unmarshal([]byte(data), &context); err != nil {
			return err
		}
		c.retireKeepalives(session, context.ExecutionContextID)
		return nil
	case "Runtime.executionContextsCleared":
		c.retireKeepalives(session, 0)
		return nil
	case "Fetch.requestPaused":
		return c.resumeRequest(session, data)
	}
	var request struct {
		RequestID string `json:"requestId"`
	}
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return err
	}
	if event == "Network.requestWillBeSent" {
		if session != "" {
			if c.owners == nil {
				c.owners = make(map[string]string)
			}
			c.owners[request.RequestID] = session
		}
	} else {
		// WebView2's root Fetch and child Network sessions share network IDs.
		// Keepalive contexts have a separate identity and ignore consumer teardown.
		c.tracker.Cancel(request.RequestID)
		delete(c.owners, request.RequestID)
	}
	return nil
}

func (c *windowsRequestCancellation) requestToken(paused pausedWebRequest) string {
	pos := strings.LastIndex(paused.Request.URLFragment, keepaliveFragment)
	if pos < 0 {
		return c.tracker.Start(paused.NetworkID)
	}
	id := paused.Request.URLFragment[pos+len(keepaliveFragment):]
	token := c.tracker.Token("keepalive:" + id)
	if token == "" {
		return "cancelled-keepalive"
	}
	return token
}

func (c *windowsRequestCancellation) resumeRequest(session, data string) error {
	var paused pausedWebRequest
	if err := json.Unmarshal([]byte(data), &paused); err != nil {
		return err
	}
	if paused.RequestID == "" {
		return fmt.Errorf("missing Fetch request ID")
	}
	token := c.requestToken(paused)
	headers := make([]requestHeader, 0, len(paused.Request.Headers)+1)
	for name, value := range paused.Request.Headers {
		if !strings.EqualFold(name, webview.RequestIDHeader) {
			headers = append(headers, requestHeader{name, value})
		}
	}
	if token != "" {
		headers = append(headers, requestHeader{webview.RequestIDHeader, token})
	}
	c.call(session, "Fetch.continueRequest", map[string]any{"requestId": paused.RequestID, "headers": headers}, func(err error) {
		if err != nil {
			c.tracker.Cancel(paused.NetworkID)
			globalApplication.error("Resuming Windows asset request: %v", err)
		}
	})
	return nil
}

func (c *windowsRequestCancellation) close() {
	if c == nil || c.closed {
		return
	}
	c.closed = true
	c.tracker.Close()
	if c.scriptID != "" {
		params, _ := json.Marshal(map[string]string{"identifier": c.scriptID})
		_ = c.chromium.CallDevTools("Page.removeScriptToEvaluateOnNewDocument", string(params), func(string, error) {})
	}
	_ = c.chromium.CallDevTools("Runtime.removeBinding", `{"name":"__wailsAssetKeepalive"}`, func(string, error) {})
	// Detaching also releases workers paused during setup.
	_ = c.chromium.CallDevTools("Target.setAutoAttach", `{"autoAttach":false,"waitForDebuggerOnStart":false,"flatten":true}`, func(string, error) {})
	// Disabling Fetch releases paused requests even when setup failed halfway.
	if err := c.chromium.CallDevTools("Fetch.disable", "{}", func(_ string, err error) {
		if err != nil {
			globalApplication.error("Disabling Windows request interception: %v", err)
		}
	}); err != nil {
		globalApplication.error("Disabling Windows request interception: %v", err)
	}
	for _, unsubscribe := range c.unsubscribe {
		if err := unsubscribe(); err != nil {
			globalApplication.error("Removing Windows request cancellation callback: %v", err)
		}
	}
	c.unsubscribe = nil
	clear(c.owners)
	clear(c.keepaliveOwners)
}

func (c *windowsRequestCancellation) targetEvent(event, data string) error {
	var target struct {
		SessionID string `json:"sessionId"`
	}
	if err := json.Unmarshal([]byte(data), &target); err != nil {
		return err
	}
	if target.SessionID == "" {
		return fmt.Errorf("missing target session ID")
	}
	if event == "Target.detachedFromTarget" {
		c.retireKeepalives(target.SessionID, 0)
		for id, session := range c.owners {
			if session == target.SessionID {
				c.tracker.Cancel(id)
				delete(c.owners, id)
			}
		}
		return nil
	}
	c.enableSession(target.SessionID, func(err error) {
		if err != nil {
			globalApplication.error("Worker request cancellation setup: %v", err)
		}
		// Always release the startup pause, even if a runtime does not support an API.
		c.call(target.SessionID, "Runtime.runIfWaitingForDebugger", struct{}{}, func(err error) {
			if err != nil {
				globalApplication.error("Resuming attached WebView target: %v", err)
				params, _ := json.Marshal(map[string]string{"sessionId": target.SessionID})
				// The root session can detach and release a paused worker even when
				// the runtime cannot address that child session directly.
				_ = c.chromium.CallDevTools("Target.detachFromTarget", string(params), func(string, error) {})
			}
		})
	})
	return nil
}

func (c *windowsRequestCancellation) keepaliveEvent(session, data string) error {
	var event struct {
		Name, Payload      string
		ExecutionContextID int64 `json:"executionContextId"`
	}
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return err
	}
	if event.Name != keepaliveBinding {
		return nil
	}
	var state struct{ ID, State string }
	if err := json.Unmarshal([]byte(event.Payload), &state); err != nil {
		return err
	}
	if len(state.ID) != 36 {
		return fmt.Errorf("invalid keepalive request identity")
	}
	id := "keepalive:" + state.ID
	switch state.State {
	case "start":
		if c.tracker.StartKeepalive(id) != "" {
			if c.keepaliveOwners == nil {
				c.keepaliveOwners = make(map[string]keepaliveOwner)
			}
			c.keepaliveOwners[id] = keepaliveOwner{session, event.ExecutionContextID}
		}
	case "failed":
		c.tracker.FailKeepalive(id)
		delete(c.keepaliveOwners, id)
	case "abort", "end":
		c.tracker.Cancel(id)
		delete(c.keepaliveOwners, id)
	default:
		return fmt.Errorf("invalid keepalive request state")
	}
	return nil
}

// Losing a JavaScript context ends delivery, not an active keepalive handler.
func (c *windowsRequestCancellation) retireKeepalives(session string, context int64) {
	for id, owner := range c.keepaliveOwners {
		if owner.session == session && (context == 0 || owner.context == context) {
			c.tracker.FailKeepalive(id)
			delete(c.keepaliveOwners, id)
		}
	}
}
