package main

import "encoding/json"

// Helpers for reading "native:*" custom-event payloads emitted by the frontend.
// A single Emit("name", obj) delivers obj on event.Data; a multi-arg emit
// delivers an array, so accept both shapes.

func firstMap(data any) map[string]any {
	switch v := data.(type) {
	case map[string]any:
		return v
	case []any:
		if len(v) > 0 {
			if m, ok := v[0].(map[string]any); ok {
				return m
			}
		}
	}
	return nil
}

func eventBool(data any, key string, def bool) bool {
	if m := firstMap(data); m != nil {
		if b, ok := m[key].(bool); ok {
			return b
		}
	}
	return def
}

func eventString(data any, key string) string {
	if m := firstMap(data); m != nil {
		if s, ok := m[key].(string); ok {
			return s
		}
	}
	return ""
}

func eventFloat(data any, key string, def float64) float64 {
	if m := firstMap(data); m != nil {
		if f, ok := m[key].(float64); ok {
			return f
		}
	}
	return def
}

func payloadJSON(data any) string {
	if m := firstMap(data); m != nil {
		if b, err := json.Marshal(m); err == nil {
			return string(b)
		}
	}
	return "{}"
}
