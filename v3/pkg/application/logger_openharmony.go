//go:build openharmony && cgo && !server

package application

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

// openharmonyConsoleHandler forwards slog records to hilog via the shell's
// console bridge (mirrors androidConsoleHandler → logcat). Go's stdout/stderr
// are not routed anywhere on OHOS, so this is how app.Logger output becomes
// visible in `hilog`.
type openharmonyConsoleHandler struct {
	level slog.Leveler
	attrs []slog.Attr
	group string
}

func (h *openharmonyConsoleHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	if h.level == nil {
		return true
	}
	return lvl >= h.level.Level()
}

func (h *openharmonyConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	var b bytes.Buffer
	b.WriteString(r.Time.Format(time.Kitchen))
	b.WriteString(" ")
	if h.group != "" {
		b.WriteString("[")
		b.WriteString(h.group)
		b.WriteString("] ")
	}
	b.WriteString(r.Message)

	writeAttr := func(a slog.Attr) {
		a.Value = a.Value.Resolve()
		b.WriteString(" ")
		b.WriteString(a.Key)
		b.WriteString("=")
		switch a.Value.Kind() {
		case slog.KindString:
			b.WriteString(a.Value.String())
		default:
			js, _ := json.Marshal(a.Value.Any())
			b.Write(js)
		}
	}
	for _, a := range h.attrs {
		writeAttr(a)
	}
	r.Attrs(func(a slog.Attr) bool { writeAttr(a); return true })

	ohBridgeVoid("log", string(mustJSON(map[string]string{
		"level": r.Level.String(),
		"msg":   b.String(),
	})))
	return nil
}

func (h *openharmonyConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &openharmonyConsoleHandler{level: h.level, attrs: append(h.attrs, attrs...), group: h.group}
}

func (h *openharmonyConsoleHandler) WithGroup(name string) slog.Handler {
	return &openharmonyConsoleHandler{level: h.level, attrs: h.attrs, group: name}
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
