package monitor

import (
	"encoding/json"
	"time"
)

// Trace is a single IPC event, serialized as one NDJSON line on the wire.
type Trace struct {
	Seq        uint64          `json:"seq"`
	Time       time.Time       `json:"time"`
	Kind       string          `json:"kind"` // "call" | "result" | "error" | "event" | "cancel"
	Dir        string          `json:"dir"`  // "in" (JS->Go) | "out" (Go->JS)
	CallID     string          `json:"callId,omitempty"`
	Object     int             `json:"object"`
	ObjectName string          `json:"objectName,omitempty"`
	Method     string          `json:"method"`
	Window     string          `json:"window,omitempty"`
	ClientID   string          `json:"clientId,omitempty"`
	Args       json.RawMessage `json:"args,omitempty"`
	Result     json.RawMessage `json:"result,omitempty"`
	Error      *TraceError     `json:"error,omitempty"`
	DurationMS float64         `json:"durationMs,omitempty"`
}

// TraceError describes a failed call or errored event.
type TraceError struct {
	Message string `json:"message"`
	Kind    string `json:"kind,omitempty"`
}
