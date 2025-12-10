//go:build ios

package application

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework WebKit
#include "webview_window_ios.h"
*/
import "C"

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"
	"unsafe"
)

// iosConsoleHandler implements slog.Handler and forwards records to the WKWebView console.
type iosConsoleHandler struct {
	level   slog.Leveler
	attrs   []slog.Attr
	group   string
}

func (h *iosConsoleHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	if h.level == nil {
		return true
	}
	return lvl >= h.level.Level()
}

func (h *iosConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	// Build a compact log line
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
		// Resolve attr values
		a.Value = a.Value.Resolve()
		b.WriteString(" ")
		b.WriteString(a.Key)
		b.WriteString("=")
		// Use JSON for complex values
		switch a.Value.Kind() {
		case slog.KindString:
			b.WriteString(strconvQuote(a.Value.String()))
		default:
			js, _ := json.Marshal(a.Value.Any())
			b.Write(js)
		}
	}

	for _, a := range h.attrs {
		writeAttr(a)
	}
	r.Attrs(func(a slog.Attr) bool { writeAttr(a); return true })

	lvl := levelToConsole(r.Level)

	msg := b.String()
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	clvl := C.CString(lvl)
	defer C.free(unsafe.Pointer(clvl))
	C.ios_console_log(clvl, cmsg)
	return nil
}

func (h *iosConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := *h
	clone.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &clone
}

func (h *iosConsoleHandler) WithGroup(name string) slog.Handler {
	clone := *h
	if clone.group == "" {
		clone.group = name
	} else if name != "" {
		clone.group = clone.group + "." + name
	}
	return &clone
}

func levelToConsole(l slog.Level) string {
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

// strconvQuote quotes a string for compact key=value output (no surrounding quotes if no spaces).
func strconvQuote(s string) string {
	if strings.IndexFunc(s, func(r rune) bool { return r == ' ' || r == '\n' || r == '\t' }) >= 0 {
		bs, _ := json.Marshal(s)
		return string(bs)
	}
	return s
}

// DefaultLogger for iOS forwards all logs to the browser console.
func DefaultLogger(level slog.Leveler) *slog.Logger {
	// Ensure there's always a leveler
	if level == nil {
		var lv slog.LevelVar
		lv.Set(slog.LevelInfo)
		level = &lv
	}
	return slog.New(&iosConsoleHandler{level: level})
}
