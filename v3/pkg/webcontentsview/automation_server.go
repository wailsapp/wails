package webcontentsview

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// AutomationServer exposes registered WebContentsView instances over a loopback WebSocket endpoint.
type AutomationServer struct {
	options  AutomationServerOptions
	endpoint AutomationEndpoint
	listener net.Listener
	server   *http.Server

	mu       sync.RWMutex
	targets  map[string]*registeredAutomationTarget
	sessions map[*automationSession]struct{}
}

type registeredAutomationTarget struct {
	target     automationTarget
	observerID uint64
}

type automationSession struct {
	server *AutomationServer
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
	send   chan any

	mu                   sync.RWMutex
	autoAttach           bool
	targetSessions       map[string]string
	sessionTargets       map[string]string
	consoleSubscriptions map[string]bool
	networkSubscriptions map[string]bool
	idleTimer            *time.Timer
}

type automationProtocolRequest struct {
	ID        json.RawMessage `json:"id,omitempty"`
	Method    string          `json:"method"`
	Params    json.RawMessage `json:"params,omitempty"`
	SessionID string          `json:"sessionId,omitempty"`
}

type automationProtocolEnvelope struct {
	ID        json.RawMessage          `json:"id,omitempty"`
	Method    string                   `json:"method,omitempty"`
	Params    any                      `json:"params,omitempty"`
	Result    any                      `json:"result,omitempty"`
	Error     *automationProtocolError `json:"error,omitempty"`
	SessionID string                   `json:"sessionId,omitempty"`
}

type automationProtocolError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// NewAutomationServer creates a new automation server with sensible loopback-only defaults.
func NewAutomationServer(options AutomationServerOptions) *AutomationServer {
	if options.Address == "" {
		options.Address = "127.0.0.1:0"
	}
	if options.Path == "" {
		options.Path = "/automation"
	}
	if options.Token == "" {
		options.Token = newAutomationToken()
	}
	if options.IdleTimeout <= 0 {
		options.IdleTimeout = 5 * time.Minute
	}
	if options.MaxPayloadBytes <= 0 {
		options.MaxPayloadBytes = 1 << 20
	}

	return &AutomationServer{
		options:  options,
		targets:  make(map[string]*registeredAutomationTarget),
		sessions: make(map[*automationSession]struct{}),
	}
}

// Start begins listening for automation clients on the configured loopback address.
func (s *AutomationServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		return nil
	}

	host, _, err := net.SplitHostPort(s.options.Address)
	if err != nil {
		return err
	}
	if host == "" || !isLoopbackHost(host) {
		return fmt.Errorf("automation server address must be loopback-only: %s", s.options.Address)
	}

	listener, err := net.Listen("tcp", s.options.Address)
	if err != nil {
		return err
	}

	s.listener = listener
	s.endpoint = AutomationEndpoint{
		URL:             "ws://" + listener.Addr().String() + s.options.Path,
		Token:           s.options.Token,
		ProtocolVersion: automationProtocolVersion,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(s.options.Path, s.handleWebSocket)
	s.server = &http.Server{
		Handler: mux,
	}

	go func(server *http.Server, ln net.Listener) {
		_ = server.Serve(ln)
	}(s.server, listener)

	return nil
}

// Stop shuts down the automation server and disconnects active sessions.
func (s *AutomationServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	server := s.server
	listener := s.listener
	s.server = nil
	s.listener = nil
	s.mu.Unlock()

	s.mu.RLock()
	sessions := make([]*automationSession, 0, len(s.sessions))
	for session := range s.sessions {
		sessions = append(sessions, session)
	}
	s.mu.RUnlock()

	for _, session := range sessions {
		session.close(websocket.StatusNormalClosure, "server stopped")
	}

	if listener != nil {
		_ = listener.Close()
	}
	if server == nil {
		return nil
	}
	return server.Shutdown(ctx)
}

// Endpoint returns the active connection details for the started server.
func (s *AutomationServer) Endpoint() AutomationEndpoint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.endpoint
}

// RegisterView adds a WebContentsView to the automation target registry.
func (s *AutomationServer) RegisterView(view *WebContentsView) error {
	if view == nil {
		return errors.New("webcontentsview is nil")
	}
	if err := view.ensureAutomationReady(); err != nil && !errors.Is(err, ErrAutomationNotSupported) {
		return err
	}
	view.syncAutomationState()
	s.registerTarget(view)
	return nil
}

// UnregisterView removes a WebContentsView from the automation target registry.
func (s *AutomationServer) UnregisterView(view *WebContentsView) {
	if view == nil {
		return
	}
	s.unregisterTarget(view.targetID())
}

// Capabilities reports the server's current domains, features, and registered targets.
func (s *AutomationServer) Capabilities() AutomationCapabilities {
	return s.capabilities()
}

func (s *AutomationServer) registerTarget(target automationTarget) {
	targetID := target.targetID()
	info := target.targetInfo()

	s.mu.Lock()
	if existing, ok := s.targets[targetID]; ok {
		target.removeAutomationObserver(existing.observerID)
	}
	observerID := target.addAutomationObserver(s.dispatchTargetEvent)
	s.targets[targetID] = &registeredAutomationTarget{
		target:     target,
		observerID: observerID,
	}
	sessions := make([]*automationSession, 0, len(s.sessions))
	for session := range s.sessions {
		sessions = append(sessions, session)
	}
	s.mu.Unlock()

	s.broadcastAll(automationProtocolEnvelope{
		Method: "Target.targetCreated",
		Params: map[string]any{
			"targetInfo": info,
		},
	})

	for _, session := range sessions {
		if !session.shouldAutoAttach() {
			continue
		}
		sessionID := session.attachTarget(targetID)
		session.trySend(automationProtocolEnvelope{
			Method:    "Target.attachedToTarget",
			SessionID: sessionID,
			Params: map[string]any{
				"sessionId":  sessionID,
				"targetInfo": info,
			},
		})
	}
}

func (s *AutomationServer) unregisterTarget(targetID string) {
	s.mu.Lock()
	registered, ok := s.targets[targetID]
	if ok {
		delete(s.targets, targetID)
	}
	sessions := make([]*automationSession, 0, len(s.sessions))
	for session := range s.sessions {
		sessions = append(sessions, session)
	}
	s.mu.Unlock()

	if !ok {
		return
	}

	registered.target.removeAutomationObserver(registered.observerID)
	for _, session := range sessions {
		session.detachTarget(targetID)
	}
	s.broadcastAll(automationProtocolEnvelope{
		Method: "Target.targetDestroyed",
		Params: map[string]any{
			"targetId": targetID,
		},
	})
}

func (s *AutomationServer) handleWebSocket(rw http.ResponseWriter, req *http.Request) {
	if !isLoopbackRemote(req.RemoteAddr) {
		http.Error(rw, "automation server only accepts loopback connections", http.StatusForbidden)
		return
	}
	if !s.authorized(req) {
		http.Error(rw, "invalid automation token", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(rw, req, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	conn.SetReadLimit(s.options.MaxPayloadBytes)

	ctx, cancel := context.WithCancel(req.Context())
	session := &automationSession{
		server:               s,
		conn:                 conn,
		ctx:                  ctx,
		cancel:               cancel,
		send:                 make(chan any, 256),
		targetSessions:       make(map[string]string),
		sessionTargets:       make(map[string]string),
		consoleSubscriptions: make(map[string]bool),
		networkSubscriptions: make(map[string]bool),
	}
	if s.options.IdleTimeout > 0 {
		session.idleTimer = time.AfterFunc(s.options.IdleTimeout, func() {
			session.close(websocket.StatusPolicyViolation, "idle timeout")
		})
	}

	s.mu.Lock()
	s.sessions[session] = struct{}{}
	s.mu.Unlock()

	defer func() {
		session.close(websocket.StatusNormalClosure, "")
		s.mu.Lock()
		delete(s.sessions, session)
		s.mu.Unlock()
	}()

	go session.writeLoop()

	session.trySend(automationProtocolEnvelope{
		Method: "App.connected",
		Params: map[string]any{
			"capabilities": s.capabilities(),
			"endpoint":     s.Endpoint(),
		},
	})

	for {
		session.touch()
		var message automationProtocolRequest
		if err := wsjson.Read(ctx, conn, &message); err != nil {
			return
		}
		session.touch()

		response := s.handleMessage(session, message)
		if response != nil {
			if err := session.sendBlocking(*response); err != nil {
				return
			}
		}
	}
}

func (s *AutomationServer) handleMessage(session *automationSession, message automationProtocolRequest) *automationProtocolEnvelope {
	if message.Method == "" {
		return automationErrorResponse(message.ID, -32600, "missing method", nil)
	}

	switch message.Method {
	case "App.getVersion":
		return automationResultResponse(message.ID, map[string]any{
			"protocolVersion": automationProtocolVersion,
			"appName":         s.options.AppName,
			"appVersion":      s.options.AppVersion,
			"appBuild":        s.options.AppBuild,
			"platform":        automationPlatform(),
		})

	case "App.getCapabilities":
		return automationResultResponse(message.ID, s.capabilities())

	case "App.listTargets", "Target.getTargets":
		return automationResultResponse(message.ID, map[string]any{
			"targets": s.listTargets(),
		})

	case "Target.attachToTarget":
		var params struct {
			TargetID string `json:"targetId"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.findTarget(params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		sessionID := session.attachTarget(target.targetID())
		return automationResultResponse(message.ID, map[string]any{
			"sessionId":  sessionID,
			"targetInfo": target.targetInfo(),
		})

	case "Target.detachFromTarget":
		var params struct {
			SessionID string `json:"sessionId"`
			TargetID  string `json:"targetId"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		if params.SessionID != "" {
			session.detachSession(params.SessionID)
		} else if params.TargetID != "" {
			session.detachTarget(params.TargetID)
		} else {
			return automationErrorResponse(message.ID, -32602, "missing sessionId or targetId", nil)
		}
		return automationResultResponse(message.ID, map[string]any{
			"detached": true,
		})

	case "Target.setAutoAttach":
		var params struct {
			AutoAttach bool `json:"autoAttach"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		session.setAutoAttach(params.AutoAttach)
		if params.AutoAttach {
			for _, target := range s.listAutomationTargets() {
				sessionID := session.attachTarget(target.targetID())
				session.trySend(automationProtocolEnvelope{
					Method:    "Target.attachedToTarget",
					SessionID: sessionID,
					Params: map[string]any{
						"sessionId":  sessionID,
						"targetInfo": target.targetInfo(),
					},
				})
			}
		}
		return automationResultResponse(message.ID, map[string]any{
			"autoAttach": params.AutoAttach,
		})

	case "Page.navigate":
		var params struct {
			TargetID string `json:"targetId"`
			URL      string `json:"url"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.resolveTarget(session, message, params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		if params.URL == "" {
			return automationErrorResponse(message.ID, -32602, "missing url", nil)
		}
		if err := target.navigate(params.URL); err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		return automationResultResponse(message.ID, map[string]any{
			"targetId": target.targetID(),
			"url":      params.URL,
		})

	case "Page.captureScreenshot":
		var params struct {
			TargetID string `json:"targetId"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.resolveTarget(session, message, params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		screenshot, err := target.captureScreenshot()
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		return automationResultResponse(message.ID, map[string]any{
			"data": screenshot,
		})

	case "Page.printToPDF":
		var params struct {
			TargetID string `json:"targetId"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.resolveTarget(session, message, params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		pdf, err := target.printToPDF()
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		return automationResultResponse(message.ID, map[string]any{
			"data": pdf,
		})

	case "Page.getNavigationHistory":
		var params struct {
			TargetID string `json:"targetId"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.resolveTarget(session, message, params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		info := target.targetInfo()
		return automationResultResponse(message.ID, map[string]any{
			"currentIndex": 0,
			"entries": []map[string]any{
				{
					"id":    0,
					"url":   info.URL,
					"title": info.Title,
				},
			},
		})

	case "Runtime.evaluate":
		var params struct {
			TargetID     string `json:"targetId"`
			Expression   string `json:"expression"`
			World        string `json:"world"`
			AwaitPromise *bool  `json:"awaitPromise"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.resolveTarget(session, message, params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		if params.Expression == "" {
			return automationErrorResponse(message.ID, -32602, "missing expression", nil)
		}
		world := automationExecutionWorldPage
		if params.World == string(automationExecutionWorldAutomation) {
			world = automationExecutionWorldAutomation
		}
		awaitPromise := true
		if params.AwaitPromise != nil {
			awaitPromise = *params.AwaitPromise
		}
		result, err := target.evaluate(params.Expression, world, awaitPromise)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		return automationResultResponse(message.ID, map[string]any{
			"result": result,
		})

	case "Runtime.getExecutionContexts":
		var params struct {
			TargetID string `json:"targetId"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.resolveTarget(session, message, params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		caps := target.targetCapabilities()
		contexts := []map[string]any{
			{
				"id":   "page",
				"name": "page",
			},
		}
		if caps.AutomationRuntime {
			contexts = append(contexts, map[string]any{
				"id":   "automation",
				"name": "automation",
			})
		}
		return automationResultResponse(message.ID, map[string]any{
			"contexts": contexts,
		})

	case "Console.enable":
		target, err := s.resolveTarget(session, message, targetIDFromParams(message.Params))
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		session.enableConsole(target.targetID())
		return automationResultResponse(message.ID, map[string]any{
			"enabled": true,
		})

	case "Console.disable":
		target, err := s.resolveTarget(session, message, targetIDFromParams(message.Params))
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		session.disableConsole(target.targetID())
		return automationResultResponse(message.ID, map[string]any{
			"enabled": false,
		})

	case "Console.getBufferedMessages":
		target, err := s.resolveTarget(session, message, targetIDFromParams(message.Params))
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		return automationResultResponse(message.ID, map[string]any{
			"messages": target.bufferedConsoleMessages(),
		})

	case "Network.enable":
		target, err := s.resolveTarget(session, message, targetIDFromParams(message.Params))
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		session.enableNetwork(target.targetID())
		return automationResultResponse(message.ID, map[string]any{
			"enabled":     true,
			"captureMode": s.captureModeForTarget(target),
		})

	case "Network.disable":
		target, err := s.resolveTarget(session, message, targetIDFromParams(message.Params))
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		session.disableNetwork(target.targetID())
		return automationResultResponse(message.ID, map[string]any{
			"enabled": false,
		})

	case "Network.getCaptureMode":
		target, err := s.resolveTarget(session, message, targetIDFromParams(message.Params))
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		return automationResultResponse(message.ID, map[string]any{
			"captureMode": s.captureModeForTarget(target),
		})

	case "Inspection.enable", "Inspection.disable", "Inspection.getStatus":
		var params struct {
			TargetID string `json:"targetId"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.resolveTarget(session, message, params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		if message.Method == "Inspection.enable" {
			if err := target.setInspectable(true); err != nil {
				return automationErrorResponse(message.ID, -32000, err.Error(), nil)
			}
		}
		if message.Method == "Inspection.disable" {
			if err := target.setInspectable(false); err != nil {
				return automationErrorResponse(message.ID, -32000, err.Error(), nil)
			}
		}
		info := target.targetInfo()
		caps := target.targetCapabilities()
		return automationResultResponse(message.ID, map[string]any{
			"enabled":   info.InspectionEnabled,
			"supported": caps.Inspection,
		})

	case "App.enableInspection", "App.disableInspection":
		var params struct {
			TargetID string `json:"targetId"`
		}
		if err := decodeAutomationParams(message.Params, &params); err != nil {
			return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
		}
		target, err := s.resolveTarget(session, message, params.TargetID)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		if message.Method == "App.enableInspection" {
			if err := target.setInspectable(true); err != nil {
				return automationErrorResponse(message.ID, -32000, err.Error(), nil)
			}
		} else {
			if err := target.setInspectable(false); err != nil {
				return automationErrorResponse(message.ID, -32000, err.Error(), nil)
			}
		}
		return automationResultResponse(message.ID, map[string]any{
			"targetInfo": target.targetInfo(),
		})
	}

	if strings.HasPrefix(message.Method, "DOM.") || strings.HasPrefix(message.Method, "Storage.") || strings.HasPrefix(message.Method, "Accessibility.") {
		target, err := s.resolveTarget(session, message, targetIDFromParams(message.Params))
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}

		switch message.Method {
		case "Storage.getCookies":
			cookies, err := target.getCookies()
			if err != nil {
				return automationErrorResponse(message.ID, -32000, err.Error(), nil)
			}
			return automationResultResponse(message.ID, map[string]any{
				"cookies": cookies,
			})

		case "Storage.setCookie":
			var params struct {
				TargetID string           `json:"targetId"`
				Cookie   AutomationCookie `json:"cookie"`
			}
			if err := decodeAutomationParams(message.Params, &params); err != nil {
				return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
			}
			if params.Cookie.Name == "" {
				return automationErrorResponse(message.ID, -32602, "missing cookie.name", nil)
			}
			if err := target.setCookie(params.Cookie); err != nil {
				return automationErrorResponse(message.ID, -32000, err.Error(), nil)
			}
			return automationResultResponse(message.ID, map[string]any{
				"stored": true,
			})

		case "Storage.deleteCookie":
			var params struct {
				TargetID string `json:"targetId"`
				Name     string `json:"name"`
				Domain   string `json:"domain"`
				Path     string `json:"path"`
			}
			if err := decodeAutomationParams(message.Params, &params); err != nil {
				return automationErrorResponse(message.ID, -32602, "invalid params", err.Error())
			}
			if params.Name == "" {
				return automationErrorResponse(message.ID, -32602, "missing name", nil)
			}
			deleted, err := target.deleteCookie(params.Name, params.Domain, params.Path)
			if err != nil {
				return automationErrorResponse(message.ID, -32000, err.Error(), nil)
			}
			return automationResultResponse(message.ID, map[string]any{
				"deleted": deleted,
			})

		case "Storage.clearCookies":
			if err := target.clearCookies(); err != nil {
				return automationErrorResponse(message.ID, -32000, err.Error(), nil)
			}
			return automationResultResponse(message.ID, map[string]any{
				"cleared": true,
			})
		}

		result, err := target.invoke(message.Method, message.Params)
		if err != nil {
			return automationErrorResponse(message.ID, -32000, err.Error(), nil)
		}
		return automationResultResponse(message.ID, result)
	}

	return automationErrorResponse(message.ID, -32601, "method not found", message.Method)
}

func (s *AutomationServer) capabilities() AutomationCapabilities {
	targets := s.listTargets()
	nativeCaps := defaultAutomationCapabilities()
	for _, target := range s.listAutomationTargets() {
		caps := target.targetCapabilities()
		nativeCaps.PageRuntime = nativeCaps.PageRuntime || caps.PageRuntime
		nativeCaps.AutomationRuntime = nativeCaps.AutomationRuntime || caps.AutomationRuntime
		nativeCaps.AsyncRuntime = nativeCaps.AsyncRuntime || caps.AsyncRuntime
		nativeCaps.DOM = nativeCaps.DOM || caps.DOM
		nativeCaps.Storage = nativeCaps.Storage || caps.Storage
		nativeCaps.Cookies = nativeCaps.Cookies || caps.Cookies
		nativeCaps.NetworkBasic = nativeCaps.NetworkBasic || caps.NetworkBasic
		nativeCaps.NetworkProxy = nativeCaps.NetworkProxy || caps.NetworkProxy
		nativeCaps.Accessibility = nativeCaps.Accessibility || caps.Accessibility
		nativeCaps.Inspection = nativeCaps.Inspection || caps.Inspection
		nativeCaps.PDF = nativeCaps.PDF || caps.PDF
	}

	domains := map[string]AutomationDomainCapabilities{
		"App": {
			Commands: []string{"App.getVersion", "App.getCapabilities", "App.listTargets", "App.enableInspection", "App.disableInspection"},
			Events:   []string{"App.connected"},
		},
		"Target": {
			Commands: []string{"Target.getTargets", "Target.attachToTarget", "Target.detachFromTarget", "Target.setAutoAttach"},
			Events:   []string{"Target.targetCreated", "Target.targetInfoChanged", "Target.targetDestroyed", "Target.attachedToTarget"},
		},
		"Page": {
			Commands: []string{"Page.navigate", "Page.captureScreenshot", "Page.getNavigationHistory"},
			Events:   []string{"Page.frameStartedLoading", "Page.frameNavigated", "Page.domContentEventFired", "Page.loadEventFired", "Page.webContentProcessTerminated", "Page.windowOpenRequested"},
		},
		"Runtime": {
			Commands: []string{"Runtime.evaluate", "Runtime.getExecutionContexts"},
			Events:   []string{"Runtime.exceptionThrown"},
			Features: map[string]any{
				"pageWorld":       nativeCaps.PageRuntime,
				"automationWorld": nativeCaps.AutomationRuntime,
				"awaitPromise":    nativeCaps.AsyncRuntime,
			},
		},
		"Console": {
			Commands: []string{"Console.enable", "Console.disable", "Console.getBufferedMessages"},
			Events:   []string{"Console.messageAdded"},
		},
	}

	if nativeCaps.PDF {
		domain := domains["Page"]
		domain.Commands = append(domain.Commands, "Page.printToPDF")
		domains["Page"] = domain
	}
	if nativeCaps.DOM {
		domains["DOM"] = AutomationDomainCapabilities{
			Commands: []string{
				"DOM.getDocument",
				"DOM.querySelector",
				"DOM.querySelectorAll",
				"DOM.queryByRole",
				"DOM.queryByText",
				"DOM.queryByLabel",
				"DOM.getOuterHTML",
				"DOM.getInnerText",
				"DOM.getAttributes",
				"DOM.getBoundingClientRect",
				"DOM.scrollIntoView",
				"DOM.focus",
				"DOM.click",
				"DOM.fill",
				"DOM.selectOption",
				"DOM.waitForSelector",
				"DOM.waitForCondition",
			},
		}
	}
	if nativeCaps.Storage {
		domains["Storage"] = AutomationDomainCapabilities{
			Commands: []string{
				"Storage.getCookies",
				"Storage.setCookie",
				"Storage.deleteCookie",
				"Storage.clearCookies",
				"Storage.getLocalStorage",
				"Storage.setLocalStorageItem",
				"Storage.removeLocalStorageItem",
				"Storage.getSessionStorage",
				"Storage.setSessionStorageItem",
				"Storage.removeSessionStorageItem",
			},
		}
	}
	if nativeCaps.NetworkBasic || nativeCaps.NetworkProxy {
		domains["Network"] = AutomationDomainCapabilities{
			Commands: []string{"Network.enable", "Network.disable", "Network.getCaptureMode"},
			Events: []string{
				"Network.requestWillBeSent",
				"Network.responseReceived",
				"Network.loadingFinished",
				"Network.loadingFailed",
			},
			Features: map[string]any{
				"basic": nativeCaps.NetworkBasic,
				"proxy": nativeCaps.NetworkProxy,
			},
		}
	}
	if nativeCaps.Accessibility {
		domains["Accessibility"] = AutomationDomainCapabilities{
			Commands: []string{
				"Accessibility.getSnapshot",
				"Accessibility.queryByRole",
				"Accessibility.queryByLabel",
			},
		}
	}
	if nativeCaps.Inspection {
		domains["Inspection"] = AutomationDomainCapabilities{
			Commands: []string{"Inspection.enable", "Inspection.disable", "Inspection.getStatus"},
		}
	}

	supportedDomains := make([]string, 0, len(domains))
	for domain := range domains {
		supportedDomains = append(supportedDomains, domain)
	}
	sort.Strings(supportedDomains)

	return AutomationCapabilities{
		ProtocolVersion:    automationProtocolVersion,
		AppName:            s.options.AppName,
		AppVersion:         s.options.AppVersion,
		AppBuild:           s.options.AppBuild,
		Platform:           automationPlatform(),
		SupportedDomains:   supportedDomains,
		Domains:            domains,
		NetworkCaptureMode: captureModeFromCapabilities(nativeCaps),
		Targets:            targets,
	}
}

func (s *AutomationServer) listTargets() []TargetInfo {
	targets := s.listAutomationTargets()
	result := make([]TargetInfo, 0, len(targets))
	for _, target := range targets {
		result = append(result, target.targetInfo())
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TargetID < result[j].TargetID
	})
	return result
}

func (s *AutomationServer) listAutomationTargets() []automationTarget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]automationTarget, 0, len(s.targets))
	for _, registered := range s.targets {
		result = append(result, registered.target)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].targetID() < result[j].targetID()
	})
	return result
}

func (s *AutomationServer) findTarget(targetID string) (automationTarget, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if targetID != "" {
		target, ok := s.targets[targetID]
		if !ok {
			return nil, fmt.Errorf("unknown target: %s", targetID)
		}
		return target.target, nil
	}

	if len(s.targets) == 1 {
		for _, target := range s.targets {
			return target.target, nil
		}
	}
	if len(s.targets) == 0 {
		return nil, errors.New("no registered automation targets")
	}
	return nil, errors.New("multiple automation targets registered; specify targetId or attach to a target")
}

func (s *AutomationServer) resolveTarget(session *automationSession, message automationProtocolRequest, explicitTargetID string) (automationTarget, error) {
	if message.SessionID != "" {
		if targetID, ok := session.targetForSession(message.SessionID); ok {
			return s.findTarget(targetID)
		}
		return nil, fmt.Errorf("unknown session: %s", message.SessionID)
	}
	return s.findTarget(explicitTargetID)
}

func (s *AutomationServer) dispatchTargetEvent(event automationTargetEvent) {
	switch event.Scope {
	case automationEventScopeAll:
		s.broadcastAll(automationProtocolEnvelope{
			Method: event.Method,
			Params: event.Params,
		})

	case automationEventScopeAttached:
		s.broadcastToTarget(event.TargetID, false, automationProtocolEnvelope{
			Method: event.Method,
			Params: event.Params,
		})

	case automationEventScopeConsole:
		s.broadcastToTarget(event.TargetID, true, automationProtocolEnvelope{
			Method: event.Method,
			Params: event.Params,
		})

	case automationEventScopeNetwork:
		s.broadcastNetworkToTarget(event.TargetID, automationProtocolEnvelope{
			Method: event.Method,
			Params: event.Params,
		})
	}
}

func (s *AutomationServer) broadcastAll(message automationProtocolEnvelope) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for session := range s.sessions {
		session.trySend(message)
	}
}

func (s *AutomationServer) broadcastToTarget(targetID string, consoleOnly bool, message automationProtocolEnvelope) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for session := range s.sessions {
		sessionID, ok := session.sessionForTarget(targetID)
		if !ok {
			continue
		}
		if consoleOnly && !session.consoleEnabled(targetID) {
			continue
		}
		outbound := message
		outbound.SessionID = sessionID
		session.trySend(outbound)
	}
}

func (s *AutomationServer) broadcastNetworkToTarget(targetID string, message automationProtocolEnvelope) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for session := range s.sessions {
		sessionID, ok := session.sessionForTarget(targetID)
		if !ok || !session.networkEnabled(targetID) {
			continue
		}
		outbound := message
		outbound.SessionID = sessionID
		session.trySend(outbound)
	}
}

func (s *AutomationServer) authorized(req *http.Request) bool {
	authHeader := strings.TrimSpace(req.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:]) == s.options.Token
	}
	return req.URL.Query().Get("token") == s.options.Token
}

func (s *automationSession) writeLoop() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case message, ok := <-s.send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
			err := wsjson.Write(writeCtx, s.conn, message)
			cancel()
			if err != nil {
				s.close(websocket.StatusInternalError, "write failure")
				return
			}
		}
	}
}

func (s *automationSession) close(status websocket.StatusCode, reason string) {
	s.mu.Lock()
	timer := s.idleTimer
	s.idleTimer = nil
	s.mu.Unlock()

	if timer != nil {
		timer.Stop()
	}

	s.cancel()
	_ = s.conn.Close(status, reason)
}

func (s *automationSession) touch() {
	s.mu.RLock()
	timer := s.idleTimer
	s.mu.RUnlock()
	if timer != nil {
		timer.Reset(s.server.options.IdleTimeout)
	}
}

func (s *automationSession) sendBlocking(message any) error {
	select {
	case s.send <- message:
		return nil
	case <-s.ctx.Done():
		return s.ctx.Err()
	}
}

func (s *automationSession) trySend(message any) bool {
	select {
	case s.send <- message:
		return true
	default:
		return false
	}
}

func (s *automationSession) attachTarget(targetID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sessionID, ok := s.targetSessions[targetID]; ok {
		return sessionID
	}
	sessionID := "target-session-" + newAutomationToken()
	s.targetSessions[targetID] = sessionID
	s.sessionTargets[sessionID] = targetID
	return sessionID
}

func (s *automationSession) detachTarget(targetID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID, ok := s.targetSessions[targetID]
	if !ok {
		return
	}
	delete(s.targetSessions, targetID)
	delete(s.sessionTargets, sessionID)
	delete(s.consoleSubscriptions, targetID)
	delete(s.networkSubscriptions, targetID)
}

func (s *automationSession) detachSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	targetID, ok := s.sessionTargets[sessionID]
	if !ok {
		return
	}
	delete(s.sessionTargets, sessionID)
	delete(s.targetSessions, targetID)
	delete(s.consoleSubscriptions, targetID)
	delete(s.networkSubscriptions, targetID)
}

func (s *automationSession) targetForSession(sessionID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	targetID, ok := s.sessionTargets[sessionID]
	return targetID, ok
}

func (s *automationSession) sessionForTarget(targetID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sessionID, ok := s.targetSessions[targetID]
	return sessionID, ok
}

func (s *automationSession) setAutoAttach(value bool) {
	s.mu.Lock()
	s.autoAttach = value
	s.mu.Unlock()
}

func (s *automationSession) shouldAutoAttach() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.autoAttach
}

func (s *automationSession) enableConsole(targetID string) {
	s.mu.Lock()
	s.consoleSubscriptions[targetID] = true
	s.mu.Unlock()
}

func (s *automationSession) disableConsole(targetID string) {
	s.mu.Lock()
	delete(s.consoleSubscriptions, targetID)
	s.mu.Unlock()
}

func (s *automationSession) consoleEnabled(targetID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.consoleSubscriptions[targetID]
}

func (s *automationSession) enableNetwork(targetID string) {
	s.mu.Lock()
	s.networkSubscriptions[targetID] = true
	s.mu.Unlock()
}

func (s *automationSession) disableNetwork(targetID string) {
	s.mu.Lock()
	delete(s.networkSubscriptions, targetID)
	s.mu.Unlock()
}

func (s *automationSession) networkEnabled(targetID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.networkSubscriptions[targetID]
}

func automationResultResponse(id json.RawMessage, result any) *automationProtocolEnvelope {
	return &automationProtocolEnvelope{
		ID:     id,
		Result: result,
	}
}

func automationErrorResponse(id json.RawMessage, code int, message string, data any) *automationProtocolEnvelope {
	return &automationProtocolEnvelope{
		ID: id,
		Error: &automationProtocolError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

func decodeAutomationParams(raw json.RawMessage, target any) error {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	return json.Unmarshal(raw, target)
}

func newAutomationToken() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buffer)
}

func isLoopbackHost(host string) bool {
	host = strings.TrimSpace(host)
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func targetIDFromParams(raw json.RawMessage) string {
	var params struct {
		TargetID string `json:"targetId"`
	}
	if decodeAutomationParams(raw, &params) != nil {
		return ""
	}
	return params.TargetID
}

func (s *AutomationServer) captureModeForTarget(target automationTarget) string {
	return captureModeFromCapabilities(target.targetCapabilities())
}

func captureModeFromCapabilities(caps automationNativeCapabilities) string {
	if caps.NetworkProxy {
		return "proxy"
	}
	if caps.NetworkBasic {
		return "basic"
	}
	return "none"
}

func isLoopbackRemote(remoteAddr string) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return false
	}
	return isLoopbackHost(host)
}
