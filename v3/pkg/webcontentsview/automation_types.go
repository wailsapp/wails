package webcontentsview

import (
	"runtime"
	"time"
)

const automationProtocolVersion = "wails.webcontentsview/0.1"

// AutomationServerOptions configures the loopback automation server used to expose WebContentsView targets.
type AutomationServerOptions struct {
	Address         string
	Path            string
	Token           string
	AppName         string
	AppVersion      string
	AppBuild        string
	IdleTimeout     time.Duration
	MaxPayloadBytes int64
}

// AutomationEndpoint describes how an automation client can connect to a started server.
type AutomationEndpoint struct {
	URL             string `json:"url"`
	Token           string `json:"token"`
	ProtocolVersion string `json:"protocolVersion"`
}

// AutomationCapabilities reports the protocol surface and registered targets available from the server.
type AutomationCapabilities struct {
	ProtocolVersion    string                                  `json:"protocolVersion"`
	AppName            string                                  `json:"appName,omitempty"`
	AppVersion         string                                  `json:"appVersion,omitempty"`
	AppBuild           string                                  `json:"appBuild,omitempty"`
	Platform           string                                  `json:"platform"`
	SupportedDomains   []string                                `json:"supportedDomains"`
	Domains            map[string]AutomationDomainCapabilities `json:"domains"`
	NetworkCaptureMode string                                  `json:"networkCaptureMode"`
	Targets            []TargetInfo                            `json:"targets"`
}

// AutomationDomainCapabilities describes the commands, events, and feature flags for a domain.
type AutomationDomainCapabilities struct {
	Commands []string       `json:"commands,omitempty"`
	Events   []string       `json:"events,omitempty"`
	Features map[string]any `json:"features,omitempty"`
}

// TargetInfo describes an automatable WebContentsView instance.
type TargetInfo struct {
	TargetID          string `json:"targetId"`
	Type              string `json:"type"`
	Name              string `json:"name,omitempty"`
	URL               string `json:"url,omitempty"`
	Title             string `json:"title,omitempty"`
	Loading           bool   `json:"loading"`
	Attached          bool   `json:"attached"`
	InspectionEnabled bool   `json:"inspectionEnabled"`
	Platform          string `json:"platform"`
}

// AutomationConsoleMessage represents a console event captured from a target.
type AutomationConsoleMessage struct {
	Level     string   `json:"level"`
	Text      string   `json:"text"`
	Args      []string `json:"args,omitempty"`
	Timestamp int64    `json:"timestamp"`
	Stack     string   `json:"stack,omitempty"`
	Source    string   `json:"source,omitempty"`
	Line      int64    `json:"line,omitempty"`
	Column    int64    `json:"column,omitempty"`
}

// AutomationException represents a page error or unhandled promise rejection.
type AutomationException struct {
	Message            string `json:"message"`
	Stack              string `json:"stack,omitempty"`
	Source             string `json:"source,omitempty"`
	Line               int64  `json:"line,omitempty"`
	Column             int64  `json:"column,omitempty"`
	Timestamp          int64  `json:"timestamp"`
	UnhandledRejection bool   `json:"unhandledRejection,omitempty"`
}

// AutomationCookie represents an HTTP cookie associated with a target data store.
type AutomationCookie struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Domain      string `json:"domain,omitempty"`
	Path        string `json:"path,omitempty"`
	Expires     int64  `json:"expires,omitempty"`
	Secure      bool   `json:"secure,omitempty"`
	HTTPOnly    bool   `json:"httpOnly,omitempty"`
	SessionOnly bool   `json:"sessionOnly,omitempty"`
	SameSite    string `json:"sameSite,omitempty"`
}

// AutomationNetworkEvent represents a best-effort network request lifecycle event.
type AutomationNetworkEvent struct {
	RequestID      string            `json:"requestId"`
	URL            string            `json:"url"`
	Method         string            `json:"method,omitempty"`
	Status         int64             `json:"status,omitempty"`
	StatusText     string            `json:"statusText,omitempty"`
	Type           string            `json:"type,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	RequestHeaders map[string]string `json:"requestHeaders,omitempty"`
	ErrorText      string            `json:"errorText,omitempty"`
	StartTime      int64             `json:"startTime,omitempty"`
	EndTime        int64             `json:"endTime,omitempty"`
	Duration       int64             `json:"duration,omitempty"`
}

type automationExecutionWorld string

const (
	automationExecutionWorldPage       automationExecutionWorld = "page"
	automationExecutionWorldAutomation automationExecutionWorld = "automation"
)

type automationRemoteObject struct {
	Type        string `json:"type"`
	Subtype     string `json:"subtype,omitempty"`
	Value       any    `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type automationNativeCapabilities struct {
	PageRuntime       bool `json:"pageRuntime"`
	AutomationRuntime bool `json:"automationRuntime"`
	AsyncRuntime      bool `json:"asyncRuntime"`
	DOM               bool `json:"dom"`
	Storage           bool `json:"storage"`
	Cookies           bool `json:"cookies"`
	NetworkBasic      bool `json:"networkBasic"`
	NetworkProxy      bool `json:"networkProxy"`
	Accessibility     bool `json:"accessibility"`
	Inspection        bool `json:"inspection"`
	PDF               bool `json:"pdf"`
}

type automationEventScope uint8

const (
	automationEventScopeAll automationEventScope = iota
	automationEventScopeAttached
	automationEventScopeConsole
	automationEventScopeNetwork
)

type automationTargetEvent struct {
	TargetID string
	Method   string
	Params   any
	Scope    automationEventScope
}

func defaultAutomationCapabilities() automationNativeCapabilities {
	return automationNativeCapabilities{}
}

func automationPlatform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}
