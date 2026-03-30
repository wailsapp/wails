package webcontentsview

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type mockAutomationTarget struct {
	mu          sync.Mutex
	info        TargetInfo
	caps        automationNativeCapabilities
	observers   map[uint64]automationObserver
	nextID      uint64
	console     []AutomationConsoleMessage
	cookies     []AutomationCookie
	lastEval    string
	lastWorld   automationExecutionWorld
	lastAwait   bool
	lastInvoke  string
	lastPayload json.RawMessage
}

func newMockAutomationTarget() *mockAutomationTarget {
	return &mockAutomationTarget{
		info: TargetInfo{
			TargetID: "wv-test",
			Type:     "webcontentsview",
			Name:     "Test WebContentsView",
			URL:      "https://example.com",
			Platform: automationPlatform(),
		},
		caps: automationNativeCapabilities{
			PageRuntime:       true,
			AutomationRuntime: true,
			AsyncRuntime:      true,
			DOM:               true,
			Storage:           true,
			Cookies:           true,
			NetworkBasic:      true,
			Accessibility:     true,
			Inspection:        true,
			PDF:               true,
		},
		observers: make(map[uint64]automationObserver),
		console: []AutomationConsoleMessage{
			{
				Level:     "log",
				Text:      "buffered",
				Timestamp: time.Now().UnixMilli(),
			},
		},
		cookies: []AutomationCookie{
			{
				Name:   "session",
				Value:  "abc123",
				Domain: "example.com",
				Path:   "/",
			},
		},
	}
}

func (m *mockAutomationTarget) targetID() string {
	return m.info.TargetID
}

func (m *mockAutomationTarget) targetInfo() TargetInfo {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.info
}

func (m *mockAutomationTarget) targetCapabilities() automationNativeCapabilities {
	return m.caps
}

func (m *mockAutomationTarget) addAutomationObserver(observer automationObserver) uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	m.observers[m.nextID] = observer
	return m.nextID
}

func (m *mockAutomationTarget) removeAutomationObserver(id uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.observers, id)
}

func (m *mockAutomationTarget) ensureAutomationReady() error {
	return nil
}

func (m *mockAutomationTarget) navigate(url string) error {
	m.mu.Lock()
	m.info.URL = url
	m.mu.Unlock()
	return nil
}

func (m *mockAutomationTarget) captureScreenshot() (string, error) {
	return "data:image/png;base64,ZmFrZQ==", nil
}

func (m *mockAutomationTarget) printToPDF() (string, error) {
	return "data:application/pdf;base64,ZmFrZQ==", nil
}

func (m *mockAutomationTarget) evaluate(expression string, world automationExecutionWorld, awaitPromise bool) (automationRemoteObject, error) {
	m.mu.Lock()
	m.lastEval = expression
	m.lastWorld = world
	m.lastAwait = awaitPromise
	m.mu.Unlock()

	return automationRemoteObject{
		Type:  "string",
		Value: "evaluated:" + expression,
	}, nil
}

func (m *mockAutomationTarget) invoke(method string, params json.RawMessage) (any, error) {
	m.mu.Lock()
	m.lastInvoke = method
	m.lastPayload = append(json.RawMessage(nil), params...)
	m.mu.Unlock()
	return map[string]any{
		"ok":     true,
		"method": method,
	}, nil
}

func (m *mockAutomationTarget) getCookies() ([]AutomationCookie, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]AutomationCookie, len(m.cookies))
	copy(result, m.cookies)
	return result, nil
}

func (m *mockAutomationTarget) setCookie(cookie AutomationCookie) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for index, existing := range m.cookies {
		if existing.Name == cookie.Name && existing.Domain == cookie.Domain && existing.Path == cookie.Path {
			m.cookies[index] = cookie
			return nil
		}
	}
	m.cookies = append(m.cookies, cookie)
	return nil
}

func (m *mockAutomationTarget) deleteCookie(name, domain, path string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	deleted := false
	filtered := m.cookies[:0]
	for _, cookie := range m.cookies {
		if cookie.Name != name {
			filtered = append(filtered, cookie)
			continue
		}
		if domain != "" && cookie.Domain != domain {
			filtered = append(filtered, cookie)
			continue
		}
		if path != "" && cookie.Path != path {
			filtered = append(filtered, cookie)
			continue
		}
		deleted = true
	}
	m.cookies = append([]AutomationCookie(nil), filtered...)
	return deleted, nil
}

func (m *mockAutomationTarget) clearCookies() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cookies = nil
	return nil
}

func (m *mockAutomationTarget) setInspectable(enabled bool) error {
	m.mu.Lock()
	m.info.InspectionEnabled = enabled
	m.mu.Unlock()
	return nil
}

func (m *mockAutomationTarget) bufferedConsoleMessages() []AutomationConsoleMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]AutomationConsoleMessage, len(m.console))
	copy(result, m.console)
	return result
}

func (m *mockAutomationTarget) emit(event automationTargetEvent) {
	m.mu.Lock()
	observers := make([]automationObserver, 0, len(m.observers))
	for _, observer := range m.observers {
		observers = append(observers, observer)
	}
	m.mu.Unlock()

	for _, observer := range observers {
		observer(event)
	}
}

func TestAutomationServerCapabilitiesAndCommands(t *testing.T) {
	server := NewAutomationServer(AutomationServerOptions{
		Address: "127.0.0.1:0",
	})
	if err := server.Start(); err != nil {
		t.Fatalf("start automation server: %v", err)
	}
	defer server.Stop(context.Background())

	target := newMockAutomationTarget()
	server.registerTarget(target)

	conn := dialAutomationServer(t, server.Endpoint())
	defer conn.Close(websocket.StatusNormalClosure, "")

	connected := readEnvelope(t, conn)
	if connected.Method != "App.connected" {
		t.Fatalf("expected App.connected event, got %#v", connected)
	}

	writeRequest(t, conn, 1, "App.getCapabilities", "", nil)
	capabilities := readEnvelope(t, conn)
	if capabilities.Error != nil {
		t.Fatalf("unexpected capabilities error: %#v", capabilities.Error)
	}

	resultMap := capabilities.Result.(map[string]any)
	if resultMap["protocolVersion"] != automationProtocolVersion {
		t.Fatalf("unexpected protocol version: %#v", resultMap["protocolVersion"])
	}

	writeRequest(t, conn, 2, "Target.attachToTarget", "", map[string]any{
		"targetId": target.targetID(),
	})
	attach := readEnvelope(t, conn)
	attachResult := attach.Result.(map[string]any)
	sessionID := attachResult["sessionId"].(string)
	if sessionID == "" {
		t.Fatal("expected session id")
	}

	writeRequest(t, conn, 3, "Runtime.evaluate", sessionID, map[string]any{
		"expression":   "document.title",
		"awaitPromise": true,
	})
	evaluate := readEnvelope(t, conn)
	if evaluate.Error != nil {
		t.Fatalf("unexpected evaluate error: %#v", evaluate.Error)
	}

	evaluateResult := evaluate.Result.(map[string]any)
	runtimeResult := evaluateResult["result"].(map[string]any)
	if runtimeResult["value"] != "evaluated:document.title" {
		t.Fatalf("unexpected runtime value: %#v", runtimeResult["value"])
	}

	writeRequest(t, conn, 4, "DOM.querySelector", sessionID, map[string]any{
		"selector": "#main",
	})
	dom := readEnvelope(t, conn)
	if dom.Error != nil {
		t.Fatalf("unexpected dom error: %#v", dom.Error)
	}

	target.mu.Lock()
	lastInvoke := target.lastInvoke
	target.mu.Unlock()
	if lastInvoke != "DOM.querySelector" {
		t.Fatalf("expected DOM.querySelector invoke, got %s", lastInvoke)
	}
}

func TestAutomationServerRoutesPageAndConsoleEvents(t *testing.T) {
	server := NewAutomationServer(AutomationServerOptions{
		Address: "127.0.0.1:0",
	})
	if err := server.Start(); err != nil {
		t.Fatalf("start automation server: %v", err)
	}
	defer server.Stop(context.Background())

	target := newMockAutomationTarget()
	server.registerTarget(target)

	conn := dialAutomationServer(t, server.Endpoint())
	defer conn.Close(websocket.StatusNormalClosure, "")
	_ = readEnvelope(t, conn)

	writeRequest(t, conn, 1, "Target.attachToTarget", "", map[string]any{
		"targetId": target.targetID(),
	})
	attach := readEnvelope(t, conn)
	sessionID := attach.Result.(map[string]any)["sessionId"].(string)

	writeRequest(t, conn, 2, "Console.enable", sessionID, nil)
	consoleEnable := readEnvelope(t, conn)
	if consoleEnable.Error != nil {
		t.Fatalf("unexpected console enable error: %#v", consoleEnable.Error)
	}

	target.emit(automationTargetEvent{
		TargetID: target.targetID(),
		Method:   "Page.loadEventFired",
		Params: map[string]any{
			"targetId": target.targetID(),
			"url":      "https://example.com/loaded",
		},
		Scope: automationEventScopeAttached,
	})
	loadEvent := readEnvelope(t, conn)
	if loadEvent.Method != "Page.loadEventFired" || loadEvent.SessionID != sessionID {
		t.Fatalf("unexpected page event: %#v", loadEvent)
	}

	target.emit(automationTargetEvent{
		TargetID: target.targetID(),
		Method:   "Console.messageAdded",
		Params: map[string]any{
			"targetId": target.targetID(),
			"message": map[string]any{
				"level": "log",
				"text":  "hello from target",
			},
		},
		Scope: automationEventScopeConsole,
	})
	consoleEvent := readEnvelope(t, conn)
	if consoleEvent.Method != "Console.messageAdded" || consoleEvent.SessionID != sessionID {
		t.Fatalf("unexpected console event: %#v", consoleEvent)
	}
}

func TestAutomationServerCookiesAndNetworkFlow(t *testing.T) {
	server := NewAutomationServer(AutomationServerOptions{
		Address: "127.0.0.1:0",
	})
	if err := server.Start(); err != nil {
		t.Fatalf("start automation server: %v", err)
	}
	defer server.Stop(context.Background())

	target := newMockAutomationTarget()
	server.registerTarget(target)

	conn := dialAutomationServer(t, server.Endpoint())
	defer conn.Close(websocket.StatusNormalClosure, "")
	_ = readEnvelope(t, conn)

	writeRequest(t, conn, 1, "Target.attachToTarget", "", map[string]any{
		"targetId": target.targetID(),
	})
	attach := readEnvelope(t, conn)
	sessionID := attach.Result.(map[string]any)["sessionId"].(string)

	writeRequest(t, conn, 2, "App.getCapabilities", "", nil)
	capabilities := readEnvelope(t, conn)
	result := capabilities.Result.(map[string]any)
	if result["networkCaptureMode"] != "basic" {
		t.Fatalf("unexpected network capture mode: %#v", result["networkCaptureMode"])
	}

	writeRequest(t, conn, 3, "Storage.getCookies", sessionID, nil)
	cookies := readEnvelope(t, conn)
	cookieResult := cookies.Result.(map[string]any)
	list := cookieResult["cookies"].([]any)
	if len(list) != 1 {
		t.Fatalf("expected one cookie, got %d", len(list))
	}

	writeRequest(t, conn, 4, "Storage.setCookie", sessionID, map[string]any{
		"cookie": map[string]any{
			"name":   "auth",
			"value":  "token",
			"domain": "example.com",
			"path":   "/",
		},
	})
	stored := readEnvelope(t, conn)
	if stored.Error != nil {
		t.Fatalf("unexpected set cookie error: %#v", stored.Error)
	}

	writeRequest(t, conn, 5, "Storage.deleteCookie", sessionID, map[string]any{
		"name":   "auth",
		"domain": "example.com",
		"path":   "/",
	})
	deleted := readEnvelope(t, conn)
	if deleted.Result.(map[string]any)["deleted"] != true {
		t.Fatalf("expected deleteCookie to delete the cookie: %#v", deleted.Result)
	}

	writeRequest(t, conn, 6, "Network.enable", sessionID, nil)
	networkEnable := readEnvelope(t, conn)
	if networkEnable.Result.(map[string]any)["captureMode"] != "basic" {
		t.Fatalf("unexpected network enable result: %#v", networkEnable.Result)
	}

	target.emit(automationTargetEvent{
		TargetID: target.targetID(),
		Method:   "Network.requestWillBeSent",
		Params: map[string]any{
			"targetId": target.targetID(),
			"event": map[string]any{
				"requestId": "net-1",
				"url":       "https://example.com/api",
				"method":    "GET",
				"type":      "fetch",
			},
		},
		Scope: automationEventScopeNetwork,
	})
	networkEvent := readEnvelope(t, conn)
	if networkEvent.Method != "Network.requestWillBeSent" || networkEvent.SessionID != sessionID {
		t.Fatalf("unexpected network event: %#v", networkEvent)
	}
}

func dialAutomationServer(t *testing.T, endpoint AutomationEndpoint) *websocket.Conn {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, endpoint.URL+"?token="+endpoint.Token, nil)
	if err != nil {
		t.Fatalf("dial automation server: %v", err)
	}
	return conn
}

func writeRequest(t *testing.T, conn *websocket.Conn, id int, method, sessionID string, params any) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := map[string]any{
		"id":     id,
		"method": method,
	}
	if sessionID != "" {
		payload["sessionId"] = sessionID
	}
	if params != nil {
		payload["params"] = params
	}
	if err := wsjson.Write(ctx, conn, payload); err != nil {
		t.Fatalf("write request %s: %v", method, err)
	}
}

func readEnvelope(t *testing.T, conn *websocket.Conn) automationProtocolEnvelope {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var envelope automationProtocolEnvelope
	if err := wsjson.Read(ctx, conn, &envelope); err != nil {
		t.Fatalf("read envelope: %v", err)
	}
	return envelope
}
