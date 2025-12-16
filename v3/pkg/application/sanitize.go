package application

import (
	"encoding/json"
	"regexp"
	"strings"
)

// DefaultRedactFields contains field names that are redacted by default.
// Matching is case-insensitive and uses substring/contains matching.
// For example, "password" will match "userPassword", "password_hash", etc.
var DefaultRedactFields = []string{
	// Authentication
	"password", "passwd", "pwd", "pass",
	"token", "bearer", "jwt", "access_token", "refresh_token",
	"secret", "client_secret",
	"apikey", "api_key", "api-key",
	"auth", "authorization", "credential",

	// Cryptographic
	"private", "privatekey", "private_key",
	"signing", "encryption_key",

	// Session/Identity
	"session", "sessionid", "session_id",
	"cookie", "csrf", "xsrf",
}

// DefaultRedactPatterns contains regex patterns that match sensitive values.
// These are applied to string values regardless of field name.
var DefaultRedactPatterns = []*regexp.Regexp{
	// JWT tokens (header.payload.signature format)
	regexp.MustCompile(`eyJ[A-Za-z0-9_-]*\.eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`),
	// Bearer tokens in values
	regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9_-]+`),
	// Common API key formats (sk_xxx, pk_xxx with optional live/test prefix)
	regexp.MustCompile(`(?i)(sk|pk)_(live|test)_[A-Za-z0-9]{10,}`),
	// Generic API key formats (api_xxx, key_xxx with 20+ chars)
	regexp.MustCompile(`(?i)(api|key)[-_][A-Za-z0-9]{20,}`),
}

// DefaultReplacement is the default string used to replace redacted values.
const DefaultReplacement = "***"

// SanitizeOptions configures automatic redaction of sensitive data
// in logs and other output.
type SanitizeOptions struct {
	// RedactFields specifies additional field names to redact.
	// Matching is case-insensitive and uses substring matching.
	// These are merged with DefaultRedactFields unless DisableDefaults is true.
	// Example: []string{"ssn", "credit_card", "dob"}
	RedactFields []string

	// RedactPatterns specifies additional regex patterns to match against string values.
	// Matching values are replaced with the Replacement string.
	// These are merged with DefaultRedactPatterns unless DisableDefaults is true.
	// Example: []*regexp.Regexp{regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)} // SSN
	RedactPatterns []*regexp.Regexp

	// CustomSanitizeFunc provides full control over sanitization.
	// When set, this function is called for every key-value pair.
	// Return (newValue, true) to use newValue, or (_, false) to use default logic.
	//
	// Parameters:
	//   - key: The field name (e.g., "password")
	//   - value: The original value
	//   - path: Full dot-separated path (e.g., "args.user.password")
	//
	// Example:
	//
	//	CustomSanitizeFunc: func(key string, value any, path string) (any, bool) {
	//	    if strings.HasPrefix(path, "args.user.") {
	//	        return "***", true
	//	    }
	//	    return nil, false // Use default sanitization
	//	}
	CustomSanitizeFunc func(key string, value any, path string) (any, bool)

	// DisableDefaults disables the default redact fields and patterns.
	// When true, only explicitly specified RedactFields and RedactPatterns apply.
	DisableDefaults bool

	// Replacement is the string used to replace redacted values.
	// Default: "***"
	Replacement string

	// Disabled completely disables sanitization.
	// Use with caution - only for debugging sanitization issues.
	Disabled bool
}

// Sanitizer handles redaction of sensitive data in logs and other output.
type Sanitizer struct {
	fields      []string
	patterns    []*regexp.Regexp
	customFunc  func(key string, value any, path string) (any, bool)
	replacement string
	disabled    bool
}

// NewSanitizer creates a new Sanitizer with the given options.
// If opts is nil, default sanitization is applied.
func NewSanitizer(opts *SanitizeOptions) *Sanitizer {
	s := &Sanitizer{
		replacement: DefaultReplacement,
	}

	if opts == nil {
		// Default configuration
		s.fields = DefaultRedactFields
		s.patterns = DefaultRedactPatterns
		return s
	}

	if opts.Disabled {
		s.disabled = true
		return s
	}

	if opts.Replacement != "" {
		s.replacement = opts.Replacement
	}

	s.customFunc = opts.CustomSanitizeFunc

	// Build field list
	if opts.DisableDefaults {
		s.fields = opts.RedactFields
		s.patterns = opts.RedactPatterns
	} else {
		// Merge with defaults
		s.fields = make([]string, 0, len(DefaultRedactFields)+len(opts.RedactFields))
		s.fields = append(s.fields, DefaultRedactFields...)
		s.fields = append(s.fields, opts.RedactFields...)

		s.patterns = make([]*regexp.Regexp, 0, len(DefaultRedactPatterns)+len(opts.RedactPatterns))
		s.patterns = append(s.patterns, DefaultRedactPatterns...)
		s.patterns = append(s.patterns, opts.RedactPatterns...)
	}

	return s
}

// SanitizeValue sanitizes a single value based on its key and path.
// The path is the full dot-separated path to the value (e.g., "args.user.password").
func (s *Sanitizer) SanitizeValue(key string, value any, path string) any {
	if s.disabled {
		return value
	}

	// Check custom function first
	if s.customFunc != nil {
		if newValue, handled := s.customFunc(key, value, path); handled {
			return newValue
		}
	}

	// Check if key matches any redact field (case-insensitive contains)
	keyLower := strings.ToLower(key)
	for _, field := range s.fields {
		if strings.Contains(keyLower, strings.ToLower(field)) {
			return s.replacement
		}
	}

	// For string values, check patterns
	if str, ok := value.(string); ok {
		for _, pattern := range s.patterns {
			if pattern.MatchString(str) {
				return s.replacement
			}
		}
	}

	// Recursively handle nested structures
	switch v := value.(type) {
	case map[string]any:
		return s.sanitizeMapInternal(v, path)
	case []any:
		return s.sanitizeSlice(v, path)
	case json.RawMessage:
		return s.SanitizeJSON(v)
	}

	return value
}

// SanitizeMap sanitizes all values in a map, recursively processing nested structures.
func (s *Sanitizer) SanitizeMap(m map[string]any) map[string]any {
	if s.disabled || m == nil {
		return m
	}
	return s.sanitizeMapInternal(m, "")
}

func (s *Sanitizer) sanitizeMapInternal(m map[string]any, parentPath string) map[string]any {
	result := make(map[string]any, len(m))
	for k, v := range m {
		path := k
		if parentPath != "" {
			path = parentPath + "." + k
		}
		result[k] = s.SanitizeValue(k, v, path)
	}
	return result
}

func (s *Sanitizer) sanitizeSlice(slice []any, parentPath string) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		// For array elements, use index in path
		path := parentPath
		if path != "" {
			path = parentPath + "[]"
		}
		// Array elements don't have a "key" for field matching,
		// but we still process nested structures and check patterns
		result[i] = s.SanitizeValue("", v, path)
	}
	return result
}

// SanitizeJSON sanitizes JSON data, returning sanitized JSON.
// If the input is invalid JSON, it returns the original data unchanged.
func (s *Sanitizer) SanitizeJSON(data []byte) []byte {
	if s.disabled || len(data) == 0 {
		return data
	}

	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		// Not valid JSON, return as-is
		// But still check if the raw string matches any patterns
		str := string(data)
		for _, pattern := range s.patterns {
			if pattern.MatchString(str) {
				return []byte(`"` + s.replacement + `"`)
			}
		}
		return data
	}

	sanitized := s.SanitizeValue("", parsed, "")
	result, err := json.Marshal(sanitized)
	if err != nil {
		return data
	}
	return result
}

// SanitizeString checks if a string matches any redact patterns and returns
// the replacement string if so, otherwise returns the original string.
func (s *Sanitizer) SanitizeString(str string) string {
	if s.disabled {
		return str
	}

	for _, pattern := range s.patterns {
		if pattern.MatchString(str) {
			return s.replacement
		}
	}
	return str
}

// IsDisabled returns true if sanitization is disabled.
func (s *Sanitizer) IsDisabled() bool {
	return s.disabled
}

// Replacement returns the replacement string used for redacted values.
func (s *Sanitizer) Replacement() string {
	return s.replacement
}
