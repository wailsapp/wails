package application

import (
	"encoding/json"
	"fmt"
)

// TransportCodec defines the interface for encoding/decoding transport data.
// This allows developers to plug in custom serialization strategies
// (e.g., MessagePack, Protobuf, custom binary formats) instead of the default base64/JSON.
type TransportCodec interface {
	// EncodeResponse encodes response data for transport.
	// The data parameter contains the raw response bytes from the message processor.
	// The contentType indicates the type of data (e.g., "application/json", "text/plain").
	// Returns the encoded data (which may be []byte, string, or any JSON-marshalable type).
	EncodeResponse(data []byte, contentType string) (interface{}, error)

	// DecodeRequest decodes incoming request arguments.
	// The data parameter contains the encoded arguments from the transport.
	// Returns the decoded arguments as a JSON string for compatibility with MessageProcessor.
	DecodeRequest(data interface{}) (string, error)
}

// DefaultCodec uses Go's standard JSON marshaling behavior.
// When []byte is marshaled to JSON, it automatically becomes a base64 string.
// This matches the frontend Base64JSONCodec for decoding.
type DefaultCodec struct{}

// NewDefaultCodec creates a new DefaultCodec instance
func NewDefaultCodec() *DefaultCodec {
	return &DefaultCodec{}
}

// EncodeResponse returns the raw bytes as-is.
// When this []byte is marshaled to JSON by the transport, it will automatically
// be encoded as a base64 string by Go's json.Marshal.
func (c *DefaultCodec) EncodeResponse(data []byte, contentType string) (interface{}, error) {
	return data, nil
}

// DecodeRequest expects a JSON string and returns it as-is
func (c *DefaultCodec) DecodeRequest(data interface{}) (string, error) {
	if str, ok := data.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("expected string, got %T", data)
}

// RawJSONCodec sends JSON data directly as a JSON object without base64 encoding.
// This is more efficient when the transport can handle JSON natively.
type RawJSONCodec struct{}

// NewRawJSONCodec creates a new RawJSONCodec instance
func NewRawJSONCodec() *RawJSONCodec {
	return &RawJSONCodec{}
}

// EncodeResponse converts the bytes to json.RawMessage so it's embedded directly in the response.
// This avoids the base64 encoding step.
func (c *RawJSONCodec) EncodeResponse(data []byte, contentType string) (interface{}, error) {
	// If the content is JSON, parse it so it's embedded as an object
	if contentType == "application/json" {
		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return jsonData, nil
	}
	// For non-JSON content, return as string
	return string(data), nil
}

// DecodeRequest expects a string or attempts to marshal any object to JSON string
func (c *RawJSONCodec) DecodeRequest(data interface{}) (string, error) {
	if str, ok := data.(string); ok {
		return str, nil
	}
	// Try to marshal to JSON string
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data to JSON: %w", err)
	}
	return string(bytes), nil
}

// RawStringCodec handles plain string data without any encoding.
// Useful for text-only protocols or when you want to avoid encoding overhead.
type RawStringCodec struct{}

// NewRawStringCodec creates a new RawStringCodec instance
func NewRawStringCodec() *RawStringCodec {
	return &RawStringCodec{}
}

// EncodeResponse converts bytes directly to a string
func (c *RawStringCodec) EncodeResponse(data []byte, contentType string) (interface{}, error) {
	return string(data), nil
}

// DecodeRequest expects a string and returns it as-is
func (c *RawStringCodec) DecodeRequest(data interface{}) (string, error) {
	if str, ok := data.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("expected string, got %T", data)
}
