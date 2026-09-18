//go:build android

package application

import (
	"bytes"
	"context"
	"log/slog"
	"time"

	"encoding/json"
)

// androidConsoleHandler implements slog.Handler and forwards records to logcat
// via androidLogf. Go's stdout/stderr are not routed anywhere on Android, so
// this is how app.Logger output becomes visible in `adb logcat` (tag "Wails").
type androidConsoleHandler struct {
	level slog.Leveler
	attrs []slog.Attr
	group string
}

func (h *androidConsoleHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	if h.level == nil {
		return true
	}
	return lvl >= h.level.Level()
}

func (h *androidConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	// Build a compact "time message key=value …" line.
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

	androidLogf(androidLevelString(r.Level), "%s", b.String())
	return nil
}

func (h *androidConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := *h
	clone.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &clone
}

func (h *androidConsoleHandler) WithGroup(name string) slog.Handler {
	clone := *h
	if clone.group == "" {
		clone.group = name
	} else if name != "" {
		clone.group = clone.group + "." + name
	}
	return &clone
}

func androidLevelString(l slog.Level) string {
	switch {
	case l <= slog.LevelDebug:
		return "debug"
	case l <= slog.LevelInfo:
		return "info"
	case l <= slog.LevelWarn:
		return "warn"
	default:
		return "error"
	}
}

// DefaultLogger for Android forwards all logs to logcat.
func DefaultLogger(level slog.Leveler) *slog.Logger {
	if level == nil {
		var lv slog.LevelVar
		lv.Set(slog.LevelInfo)
		level = &lv
	}
	return slog.New(&androidConsoleHandler{level: level})
}
