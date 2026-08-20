package application

import (
	"context"
	"encoding/json"
	"log/slog"
	"regexp"
	"strings"
	"testing"
)

func TestNewSanitizer_Defaults(t *testing.T) {
	s := NewSanitizer(nil)

	if s.IsDisabled() {
		t.Error("expected sanitizer to be enabled by default")
	}
	if s.Replacement() != DefaultReplacement {
		t.Errorf("expected replacement %q, got %q", DefaultReplacement, s.Replacement())
	}
	if len(s.fields) != len(DefaultRedactFields) {
		t.Errorf("expected %d default fields, got %d", len(DefaultRedactFields), len(s.fields))
	}
	if len(s.patterns) != len(DefaultRedactPatterns) {
		t.Errorf("expected %d default patterns, got %d", len(DefaultRedactPatterns), len(s.patterns))
	}
}

func TestNewSanitizer_Disabled(t *testing.T) {
	s := NewSanitizer(&SanitizeOptions{Disabled: true})

	if !s.IsDisabled() {
		t.Error("expected sanitizer to be disabled")
	}

	// Disabled sanitizer should pass through values unchanged
	result := s.SanitizeValue("password", "secret123", "password")
	if result != "secret123" {
		t.Errorf("expected unchanged value, got %v", result)
	}
}

func TestNewSanitizer_CustomReplacement(t *testing.T) {
	s := NewSanitizer(&SanitizeOptions{Replacement: "[REDACTED]"})

	if s.Replacement() != "[REDACTED]" {
		t.Errorf("expected replacement [REDACTED], got %q", s.Replacement())
	}

	result := s.SanitizeValue("password", "secret123", "password")
	if result != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", result)
	}
}

func TestNewSanitizer_DisableDefaults(t *testing.T) {
	s := NewSanitizer(&SanitizeOptions{
		DisableDefaults: true,
		RedactFields:    []string{"custom_field"},
	})

	// Default field should not be redacted
	result := s.SanitizeValue("password", "secret123", "password")
	if result == s.Replacement() {
		t.Error("expected password to NOT be redacted when defaults disabled")
	}

	// Custom field should be redacted
	result = s.SanitizeValue("custom_field", "value", "custom_field")
	if result != s.Replacement() {
		t.Error("expected custom_field to be redacted")
	}
}

func TestNewSanitizer_MergeFields(t *testing.T) {
	s := NewSanitizer(&SanitizeOptions{
		RedactFields: []string{"custom_field"},
	})

	// Default field should still be redacted
	result := s.SanitizeValue("password", "secret123", "password")
	if result != s.Replacement() {
		t.Error("expected password to be redacted")
	}

	// Custom field should also be redacted
	result = s.SanitizeValue("custom_field", "value", "custom_field")
	if result != s.Replacement() {
		t.Error("expected custom_field to be redacted")
	}
}

func TestSanitizeValue_FieldMatching(t *testing.T) {
	s := NewSanitizer(nil)

	tests := []struct {
		name     string
		key      string
		value    any
		expected any
	}{
		// Exact matches
		{"exact password", "password", "secret", "***"},
		{"exact token", "token", "abc123", "***"},
		{"exact secret", "secret", "mysecret", "***"},
		{"exact apikey", "apikey", "key123", "***"},

		// Case insensitive
		{"uppercase PASSWORD", "PASSWORD", "secret", "***"},
		{"mixed case Token", "Token", "abc123", "***"},
		{"mixed case API_KEY", "API_KEY", "key123", "***"},

		// Substring/contains matching
		{"userPassword contains password", "userPassword", "secret", "***"},
		{"password_hash contains password", "password_hash", "hash", "***"},
		{"auth_token contains token", "auth_token", "abc", "***"},
		{"jwt_token contains token", "jwt_token", "xyz", "***"},
		{"MySecretKey contains secret", "MySecretKey", "value", "***"},

		// Non-sensitive fields should pass through
		{"username not sensitive", "username", "john", "john"},
		{"email not sensitive", "email", "john@example.com", "john@example.com"},
		{"count not sensitive", "count", 42, 42},
		{"enabled not sensitive", "enabled", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.SanitizeValue(tt.key, tt.value, tt.key)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSanitizeValue_PatternMatching(t *testing.T) {
	s := NewSanitizer(nil)

	tests := []struct {
		name        string
		key         string
		value       string
		shouldRedact bool
	}{
		// JWT tokens
		{"JWT token", "data", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U", true},

		// Bearer tokens
		{"Bearer token", "header", "Bearer abc123xyz789", true},
		{"bearer lowercase", "header", "bearer mytoken123", true},

		// API keys (Stripe-style sk_live_xxx / pk_test_xxx)
		{"sk_live API key", "key", "sk_live_abcdefghij12", true},
		{"pk_test API key", "key", "pk_test_abcdefghij12", true},
		// Generic API keys (api_xxx with 20+ chars)
		{"api_ key", "key", "api_abcdefghij12345678901234", true},

		// Non-matching
		{"normal string", "message", "Hello, World!", false},
		{"short sk key", "key", "sk_live_short", false}, // Too short (needs 10+ after prefix)
		{"email", "data", "user@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.SanitizeValue(tt.key, tt.value, tt.key)
			if tt.shouldRedact {
				if result != s.Replacement() {
					t.Errorf("expected value to be redacted, got %v", result)
				}
			} else {
				if result == s.Replacement() {
					t.Errorf("expected value to NOT be redacted, got %v", result)
				}
			}
		})
	}
}

func TestSanitizeValue_NestedMap(t *testing.T) {
	s := NewSanitizer(nil)

	input := map[string]any{
		"user": map[string]any{
			"name":     "John",
			"password": "secret123",
			"settings": map[string]any{
				"theme":     "dark",
				"api_token": "mytoken",
			},
		},
		"count": 42,
	}

	result := s.SanitizeValue("data", input, "data")
	resultMap, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %T", result)
	}

	userMap := resultMap["user"].(map[string]any)
	if userMap["name"] != "John" {
		t.Errorf("expected name to be unchanged, got %v", userMap["name"])
	}
	if userMap["password"] != "***" {
		t.Errorf("expected password to be redacted, got %v", userMap["password"])
	}

	settingsMap := userMap["settings"].(map[string]any)
	if settingsMap["theme"] != "dark" {
		t.Errorf("expected theme to be unchanged, got %v", settingsMap["theme"])
	}
	if settingsMap["api_token"] != "***" {
		t.Errorf("expected api_token to be redacted, got %v", settingsMap["api_token"])
	}

	if resultMap["count"] != 42 {
		t.Errorf("expected count to be unchanged, got %v", resultMap["count"])
	}
}

func TestSanitizeValue_Slice(t *testing.T) {
	s := NewSanitizer(nil)

	input := []any{
		"normal string",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.signature123",
		map[string]any{
			"password": "secret",
			"name":     "test",
		},
	}

	result := s.SanitizeValue("items", input, "items")
	resultSlice, ok := result.([]any)
	if !ok {
		t.Fatalf("expected slice result, got %T", result)
	}

	if resultSlice[0] != "normal string" {
		t.Errorf("expected first element unchanged, got %v", resultSlice[0])
	}
	if resultSlice[1] != "***" {
		t.Errorf("expected JWT to be redacted, got %v", resultSlice[1])
	}

	mapElement := resultSlice[2].(map[string]any)
	if mapElement["password"] != "***" {
		t.Errorf("expected password in slice element to be redacted, got %v", mapElement["password"])
	}
	if mapElement["name"] != "test" {
		t.Errorf("expected name in slice element unchanged, got %v", mapElement["name"])
	}
}

func TestSanitizeMap(t *testing.T) {
	s := NewSanitizer(nil)

	input := map[string]any{
		"username":    "john",
		"password":    "secret123",
		"accessToken": "token123",
		"data": map[string]any{
			"secret_key": "mysecret",
			"public_id":  "pub123",
		},
	}

	result := s.SanitizeMap(input)

	if result["username"] != "john" {
		t.Errorf("expected username unchanged, got %v", result["username"])
	}
	if result["password"] != "***" {
		t.Errorf("expected password redacted, got %v", result["password"])
	}
	if result["accessToken"] != "***" {
		t.Errorf("expected accessToken redacted, got %v", result["accessToken"])
	}

	dataMap := result["data"].(map[string]any)
	if dataMap["secret_key"] != "***" {
		t.Errorf("expected secret_key redacted, got %v", dataMap["secret_key"])
	}
	if dataMap["public_id"] != "pub123" {
		t.Errorf("expected public_id unchanged, got %v", dataMap["public_id"])
	}
}

func TestSanitizeMap_Nil(t *testing.T) {
	s := NewSanitizer(nil)
	result := s.SanitizeMap(nil)
	if result != nil {
		t.Errorf("expected nil result for nil input, got %v", result)
	}
}

func TestSanitizeJSON(t *testing.T) {
	s := NewSanitizer(nil)

	input := `{"username":"john","password":"secret123","nested":{"token":"abc"}}`
	result := s.SanitizeJSON([]byte(input))

	var parsed map[string]any
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("failed to parse result JSON: %v", err)
	}

	if parsed["username"] != "john" {
		t.Errorf("expected username unchanged, got %v", parsed["username"])
	}
	if parsed["password"] != "***" {
		t.Errorf("expected password redacted, got %v", parsed["password"])
	}

	nested := parsed["nested"].(map[string]any)
	if nested["token"] != "***" {
		t.Errorf("expected nested token redacted, got %v", nested["token"])
	}
}

func TestSanitizeJSON_InvalidJSON(t *testing.T) {
	s := NewSanitizer(nil)

	// Invalid JSON should be returned as-is (unless it matches a pattern)
	input := []byte("not valid json")
	result := s.SanitizeJSON(input)
	if string(result) != string(input) {
		t.Errorf("expected invalid JSON to be returned unchanged, got %s", result)
	}

	// But if it contains a pattern match, it should be redacted
	inputWithJWT := []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.sig")
	result = s.SanitizeJSON(inputWithJWT)
	if string(result) != `"***"` {
		t.Errorf("expected JWT in invalid JSON to be redacted, got %s", result)
	}
}

func TestSanitizeJSON_Empty(t *testing.T) {
	s := NewSanitizer(nil)
	result := s.SanitizeJSON([]byte{})
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %s", result)
	}
}

func TestSanitizeString(t *testing.T) {
	s := NewSanitizer(nil)

	// Normal string - unchanged
	result := s.SanitizeString("Hello, World!")
	if result != "Hello, World!" {
		t.Errorf("expected unchanged string, got %s", result)
	}

	// JWT - redacted
	result = s.SanitizeString("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.sig")
	if result != "***" {
		t.Errorf("expected JWT redacted, got %s", result)
	}
}

func TestCustomSanitizeFunc(t *testing.T) {
	s := NewSanitizer(&SanitizeOptions{
		CustomSanitizeFunc: func(key string, value any, path string) (any, bool) {
			// Custom handling for user paths
			if strings.HasPrefix(path, "user.") {
				return "[USER_DATA]", true
			}
			// Let defaults handle password
			if key == "password" {
				return nil, false
			}
			return nil, false
		},
	})

	// Custom function should handle user paths
	result := s.SanitizeValue("email", "john@example.com", "user.email")
	if result != "[USER_DATA]" {
		t.Errorf("expected custom handling for user path, got %v", result)
	}

	// Default should still handle password
	result = s.SanitizeValue("password", "secret", "password")
	if result != "***" {
		t.Errorf("expected default handling for password, got %v", result)
	}

	// Non-sensitive field should pass through
	result = s.SanitizeValue("count", 42, "count")
	if result != 42 {
		t.Errorf("expected count unchanged, got %v", result)
	}
}

func TestCustomSanitizeFunc_FullOverride(t *testing.T) {
	s := NewSanitizer(&SanitizeOptions{
		CustomSanitizeFunc: func(key string, value any, path string) (any, bool) {
			// Handle everything - even "password" should use custom logic
			if key == "password" {
				return "[CUSTOM_REDACTED]", true
			}
			return value, true // Return original for everything else
		},
	})

	// Custom function overrides default password handling
	result := s.SanitizeValue("password", "secret", "password")
	if result != "[CUSTOM_REDACTED]" {
		t.Errorf("expected custom redaction, got %v", result)
	}

	// Custom function returns original for other fields
	result = s.SanitizeValue("token", "abc123", "token")
	if result != "abc123" {
		t.Errorf("expected token unchanged by custom func, got %v", result)
	}
}

func TestCustomPatterns(t *testing.T) {
	// Custom pattern for SSN
	ssnPattern := regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)

	s := NewSanitizer(&SanitizeOptions{
		RedactPatterns: []*regexp.Regexp{ssnPattern},
	})

	// SSN should be redacted
	result := s.SanitizeValue("data", "SSN: 123-45-6789", "data")
	if result != "***" {
		t.Errorf("expected SSN redacted, got %v", result)
	}

	// Default patterns should still work
	result = s.SanitizeValue("data", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.sig", "data")
	if result != "***" {
		t.Errorf("expected JWT still redacted, got %v", result)
	}
}

func TestSanitizeValue_JSONRawMessage(t *testing.T) {
	s := NewSanitizer(nil)

	input := json.RawMessage(`{"password":"secret","name":"test"}`)
	result := s.SanitizeValue("data", input, "data")

	resultBytes, ok := result.([]byte)
	if !ok {
		t.Fatalf("expected []byte result, got %T", result)
	}

	var parsed map[string]any
	if err := json.Unmarshal(resultBytes, &parsed); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	if parsed["password"] != "***" {
		t.Errorf("expected password redacted in RawMessage, got %v", parsed["password"])
	}
	if parsed["name"] != "test" {
		t.Errorf("expected name unchanged, got %v", parsed["name"])
	}
}

func TestPathTracking(t *testing.T) {
	var capturedPaths []string

	s := NewSanitizer(&SanitizeOptions{
		CustomSanitizeFunc: func(key string, value any, path string) (any, bool) {
			capturedPaths = append(capturedPaths, path)
			return nil, false // Let defaults handle it
		},
	})

	input := map[string]any{
		"user": map[string]any{
			"password": "secret",
			"profile": map[string]any{
				"name": "John",
			},
		},
	}

	s.SanitizeMap(input)

	expectedPaths := []string{
		"user",
		"user.password",
		"user.profile",
		"user.profile.name",
	}

	for _, expected := range expectedPaths {
		found := false
		for _, captured := range capturedPaths {
			if captured == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %q to be captured, got paths: %v", expected, capturedPaths)
		}
	}
}

// Benchmark tests
func BenchmarkSanitizeMap_Small(b *testing.B) {
	s := NewSanitizer(nil)
	input := map[string]any{
		"username": "john",
		"password": "secret123",
		"email":    "john@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SanitizeMap(input)
	}
}

func BenchmarkSanitizeMap_Nested(b *testing.B) {
	s := NewSanitizer(nil)
	input := map[string]any{
		"user": map[string]any{
			"name":     "John",
			"password": "secret123",
			"settings": map[string]any{
				"theme":     "dark",
				"api_token": "mytoken",
			},
		},
		"metadata": map[string]any{
			"created": "2024-01-01",
			"secret":  "hidden",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SanitizeMap(input)
	}
}

func BenchmarkSanitizeJSON(b *testing.B) {
	s := NewSanitizer(nil)
	input := []byte(`{"username":"john","password":"secret123","nested":{"token":"abc","data":"value"}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SanitizeJSON(input)
	}
}

func BenchmarkSanitizeValue_NoRedaction(b *testing.B) {
	s := NewSanitizer(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SanitizeValue("username", "john", "username")
	}
}

func BenchmarkSanitizeValue_WithRedaction(b *testing.B) {
	s := NewSanitizer(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SanitizeValue("password", "secret123", "password")
	}
}

func BenchmarkSanitizer_Disabled(b *testing.B) {
	s := NewSanitizer(&SanitizeOptions{Disabled: true})
	input := map[string]any{
		"username": "john",
		"password": "secret123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SanitizeMap(input)
	}
}

// SanitizingHandler tests

type testLogEntry struct {
	Level   slog.Level
	Message string
	Attrs   map[string]any
}

type testHandler struct {
	entries []testLogEntry
}

func (h *testHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *testHandler) Handle(_ context.Context, record slog.Record) error {
	entry := testLogEntry{
		Level:   record.Level,
		Message: record.Message,
		Attrs:   make(map[string]any),
	}
	record.Attrs(func(attr slog.Attr) bool {
		h.collectAttr(entry.Attrs, "", attr)
		return true
	})
	h.entries = append(h.entries, entry)
	return nil
}

func (h *testHandler) collectAttr(m map[string]any, prefix string, attr slog.Attr) {
	key := attr.Key
	if prefix != "" {
		key = prefix + "." + key
	}

	if attr.Value.Kind() == slog.KindGroup {
		for _, ga := range attr.Value.Group() {
			h.collectAttr(m, key, ga)
		}
		return
	}

	m[key] = attr.Value.Any()
}

func (h *testHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *testHandler) WithGroup(name string) slog.Handler {
	return h
}

func TestSanitizingHandler_Basic(t *testing.T) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(nil)
	handler := NewSanitizingHandler(underlying, sanitizer)
	logger := slog.New(handler)

	logger.Info("test message",
		"username", "john",
		"password", "secret123",
		"count", 42,
	)

	if len(underlying.entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(underlying.entries))
	}

	entry := underlying.entries[0]
	if entry.Message != "test message" {
		t.Errorf("expected message 'test message', got %q", entry.Message)
	}
	if entry.Attrs["username"] != "john" {
		t.Errorf("expected username 'john', got %v", entry.Attrs["username"])
	}
	if entry.Attrs["password"] != "***" {
		t.Errorf("expected password '***', got %v", entry.Attrs["password"])
	}
	if entry.Attrs["count"] != int64(42) {
		t.Errorf("expected count 42, got %v", entry.Attrs["count"])
	}
}

func TestSanitizingHandler_NestedGroup(t *testing.T) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(nil)
	handler := NewSanitizingHandler(underlying, sanitizer)
	logger := slog.New(handler)

	logger.Info("test message",
		slog.Group("user",
			slog.String("name", "john"),
			slog.String("password", "secret123"),
		),
	)

	if len(underlying.entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(underlying.entries))
	}

	entry := underlying.entries[0]
	if entry.Attrs["user.name"] != "john" {
		t.Errorf("expected user.name 'john', got %v", entry.Attrs["user.name"])
	}
	if entry.Attrs["user.password"] != "***" {
		t.Errorf("expected user.password '***', got %v", entry.Attrs["user.password"])
	}
}

func TestSanitizingHandler_PatternMatching(t *testing.T) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(nil)
	handler := NewSanitizingHandler(underlying, sanitizer)
	logger := slog.New(handler)

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.signature"
	logger.Info("auth",
		"data", jwt,
		"normal", "hello",
	)

	if len(underlying.entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(underlying.entries))
	}

	entry := underlying.entries[0]
	if entry.Attrs["data"] != "***" {
		t.Errorf("expected JWT redacted, got %v", entry.Attrs["data"])
	}
	if entry.Attrs["normal"] != "hello" {
		t.Errorf("expected normal 'hello', got %v", entry.Attrs["normal"])
	}
}

func TestSanitizingHandler_Disabled(t *testing.T) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(&SanitizeOptions{Disabled: true})
	handler := NewSanitizingHandler(underlying, sanitizer)
	logger := slog.New(handler)

	logger.Info("test", "password", "secret123")

	if len(underlying.entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(underlying.entries))
	}

	entry := underlying.entries[0]
	if entry.Attrs["password"] != "secret123" {
		t.Errorf("expected password unchanged when disabled, got %v", entry.Attrs["password"])
	}
}

func TestSanitizingHandler_WithAttrs(t *testing.T) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(nil)
	handler := NewSanitizingHandler(underlying, sanitizer)

	// Create handler with pre-set attrs
	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("api_key", "secret_key_123"),
	})

	logger := slog.New(handlerWithAttrs)
	logger.Info("test", "username", "john")

	// The WithAttrs should have sanitized the api_key
	// Note: Our test handler doesn't properly handle WithAttrs,
	// but we're testing that WithAttrs returns a SanitizingHandler
	if _, ok := handlerWithAttrs.(*SanitizingHandler); !ok {
		t.Error("expected WithAttrs to return a SanitizingHandler")
	}
}

func TestSanitizingHandler_WithGroup(t *testing.T) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(nil)
	handler := NewSanitizingHandler(underlying, sanitizer)

	handlerWithGroup := handler.WithGroup("request")

	if h, ok := handlerWithGroup.(*SanitizingHandler); !ok {
		t.Error("expected WithGroup to return a SanitizingHandler")
	} else {
		if len(h.groups) != 1 || h.groups[0] != "request" {
			t.Errorf("expected groups ['request'], got %v", h.groups)
		}
	}
}

func TestSanitizingHandler_AnyValue(t *testing.T) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(nil)
	handler := NewSanitizingHandler(underlying, sanitizer)
	logger := slog.New(handler)

	// Log a map as Any value
	data := map[string]any{
		"username": "john",
		"password": "secret",
	}
	logger.Info("test", "data", data)

	if len(underlying.entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(underlying.entries))
	}

	entry := underlying.entries[0]
	dataResult, ok := entry.Attrs["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data to be map, got %T", entry.Attrs["data"])
	}

	if dataResult["username"] != "john" {
		t.Errorf("expected username 'john', got %v", dataResult["username"])
	}
	if dataResult["password"] != "***" {
		t.Errorf("expected password '***', got %v", dataResult["password"])
	}
}

func TestWrapLoggerWithSanitizer(t *testing.T) {
	underlying := &testHandler{}
	originalLogger := slog.New(underlying)
	sanitizer := NewSanitizer(nil)

	wrappedLogger := WrapLoggerWithSanitizer(originalLogger, sanitizer)

	wrappedLogger.Info("test", "password", "secret123")

	if len(underlying.entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(underlying.entries))
	}

	entry := underlying.entries[0]
	if entry.Attrs["password"] != "***" {
		t.Errorf("expected password '***', got %v", entry.Attrs["password"])
	}
}

func TestWrapLoggerWithSanitizer_Nil(t *testing.T) {
	result := WrapLoggerWithSanitizer(nil, nil)
	if result != nil {
		t.Error("expected nil result for nil logger")
	}
}

func BenchmarkSanitizingHandler(b *testing.B) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(nil)
	handler := NewSanitizingHandler(underlying, sanitizer)
	logger := slog.New(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		underlying.entries = nil // Reset
		logger.Info("test",
			"username", "john",
			"password", "secret123",
			"count", 42,
		)
	}
}

func BenchmarkSanitizingHandler_Disabled(b *testing.B) {
	underlying := &testHandler{}
	sanitizer := NewSanitizer(&SanitizeOptions{Disabled: true})
	handler := NewSanitizingHandler(underlying, sanitizer)
	logger := slog.New(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		underlying.entries = nil
		logger.Info("test",
			"username", "john",
			"password", "secret123",
			"count", 42,
		)
	}
}
