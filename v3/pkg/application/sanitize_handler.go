package application

import (
	"context"
	"log/slog"
)

// SanitizingHandler wraps a slog.Handler to sanitize sensitive data in log attributes.
// All log attributes are passed through the Sanitizer before being forwarded to the
// underlying handler.
type SanitizingHandler struct {
	handler   slog.Handler
	sanitizer *Sanitizer
	groups    []string // Track nested groups for path building
}

// NewSanitizingHandler creates a new SanitizingHandler that wraps the given handler.
// If sanitizer is nil, a default sanitizer is created.
func NewSanitizingHandler(handler slog.Handler, sanitizer *Sanitizer) *SanitizingHandler {
	if sanitizer == nil {
		sanitizer = NewSanitizer(nil)
	}
	return &SanitizingHandler{
		handler:   handler,
		sanitizer: sanitizer,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *SanitizingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle sanitizes the record's attributes and forwards to the underlying handler.
func (h *SanitizingHandler) Handle(ctx context.Context, record slog.Record) error {
	if h.sanitizer.IsDisabled() {
		return h.handler.Handle(ctx, record)
	}

	// Create a new record with sanitized attributes
	newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)

	// Sanitize each attribute
	record.Attrs(func(attr slog.Attr) bool {
		sanitized := h.sanitizeAttr(attr, h.buildPath(""))
		newRecord.AddAttrs(sanitized)
		return true
	})

	return h.handler.Handle(ctx, newRecord)
}

// WithAttrs returns a new handler with the given attributes added.
// The attributes are sanitized before being stored.
func (h *SanitizingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.sanitizer.IsDisabled() {
		return &SanitizingHandler{
			handler:   h.handler.WithAttrs(attrs),
			sanitizer: h.sanitizer,
			groups:    h.groups,
		}
	}

	// Sanitize attributes before passing to underlying handler
	sanitizedAttrs := make([]slog.Attr, len(attrs))
	for i, attr := range attrs {
		sanitizedAttrs[i] = h.sanitizeAttr(attr, h.buildPath(""))
	}

	return &SanitizingHandler{
		handler:   h.handler.WithAttrs(sanitizedAttrs),
		sanitizer: h.sanitizer,
		groups:    h.groups,
	}
}

// WithGroup returns a new handler with the given group name.
func (h *SanitizingHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups), len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups = append(newGroups, name)

	return &SanitizingHandler{
		handler:   h.handler.WithGroup(name),
		sanitizer: h.sanitizer,
		groups:    newGroups,
	}
}

// buildPath constructs the full path for an attribute key.
func (h *SanitizingHandler) buildPath(key string) string {
	if len(h.groups) == 0 {
		return key
	}

	path := ""
	for _, g := range h.groups {
		if path != "" {
			path += "."
		}
		path += g
	}
	if key != "" {
		if path != "" {
			path += "."
		}
		path += key
	}
	return path
}

// sanitizeAttr recursively sanitizes an slog.Attr.
func (h *SanitizingHandler) sanitizeAttr(attr slog.Attr, parentPath string) slog.Attr {
	key := attr.Key
	path := key
	if parentPath != "" {
		path = parentPath + "." + key
	}

	// Handle group attributes (nested)
	if attr.Value.Kind() == slog.KindGroup {
		groupAttrs := attr.Value.Group()
		sanitizedGroup := make([]slog.Attr, len(groupAttrs))
		for i, ga := range groupAttrs {
			sanitizedGroup[i] = h.sanitizeAttr(ga, path)
		}
		return slog.Attr{Key: key, Value: slog.GroupValue(sanitizedGroup...)}
	}

	// Sanitize the value
	sanitizedValue := h.sanitizeValue(key, attr.Value, path)
	return slog.Attr{Key: key, Value: sanitizedValue}
}

// sanitizeValue sanitizes an slog.Value using the Sanitizer.
func (h *SanitizingHandler) sanitizeValue(key string, value slog.Value, path string) slog.Value {
	// Resolve LogValuer interfaces
	value = value.Resolve()

	switch value.Kind() {
	case slog.KindString:
		sanitized := h.sanitizer.SanitizeValue(key, value.String(), path)
		if str, ok := sanitized.(string); ok {
			return slog.StringValue(str)
		}
		return slog.AnyValue(sanitized)

	case slog.KindAny:
		anyVal := value.Any()
		sanitized := h.sanitizer.SanitizeValue(key, anyVal, path)
		return slog.AnyValue(sanitized)

	case slog.KindGroup:
		// Already handled in sanitizeAttr
		return value

	default:
		// For other kinds (Int64, Uint64, Float64, Bool, Time, Duration),
		// we still need to check if the key matches a redact field
		// Convert to any and sanitize
		anyVal := value.Any()
		sanitized := h.sanitizer.SanitizeValue(key, anyVal, path)

		// If it was redacted (became a string), return as string
		if str, ok := sanitized.(string); ok && str == h.sanitizer.Replacement() {
			return slog.StringValue(str)
		}

		// Otherwise return original value
		return value
	}
}

// WrapLoggerWithSanitizer wraps an existing slog.Logger with sanitization.
// This is a convenience function for wrapping a logger.
func WrapLoggerWithSanitizer(logger *slog.Logger, sanitizer *Sanitizer) *slog.Logger {
	if logger == nil {
		return nil
	}
	handler := NewSanitizingHandler(logger.Handler(), sanitizer)
	return slog.New(handler)
}
