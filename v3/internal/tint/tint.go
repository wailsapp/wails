// Package tint provides a colorized slog.Handler for development logging.
package tint

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGray   = "\033[38;5;245m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[97m"
)

// Options configures the Handler.
type Options struct {
	Level      slog.Leveler
	TimeFormat string
	NoColor    bool
}

// NewHandler returns a new Handler that writes colored log records to w.
func NewHandler(w io.Writer, opts *Options) slog.Handler {
	if opts == nil {
		opts = &Options{}
	}
	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}
	tf := opts.TimeFormat
	if tf == "" {
		tf = time.StampMilli
	}
	return &handler{w: w, opts: *opts, timeFormat: tf}
}

type handler struct {
	mu         sync.Mutex
	w          io.Writer
	opts       Options
	timeFormat string
	preformatted []byte
	groups     []string
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *handler) Handle(_ context.Context, r slog.Record) error {
	var buf bytes.Buffer

	// Timestamp
	if !r.Time.IsZero() {
		ts := r.Time.Format(h.timeFormat)
		if h.opts.NoColor {
			buf.WriteString(ts)
		} else {
			buf.WriteString(colorGray + ts + colorReset)
		}
		buf.WriteByte(' ')
	}

	// Level
	levelStr, color := levelInfo(r.Level)
	if h.opts.NoColor {
		buf.WriteString(levelStr)
	} else {
		buf.WriteString(color + levelStr + colorReset)
	}
	buf.WriteByte(' ')

	// Message
	if h.opts.NoColor {
		buf.WriteString(r.Message)
	} else {
		buf.WriteString(colorWhite + r.Message + colorReset)
	}

	// Pre-formatted group attrs
	if len(h.preformatted) > 0 {
		buf.WriteByte(' ')
		buf.Write(h.preformatted)
	}

	// Record attrs
	r.Attrs(func(a slog.Attr) bool {
		buf.WriteByte(' ')
		h.appendAttr(&buf, a, h.groups)
		return true
	})

	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf.Bytes())
	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var buf bytes.Buffer
	for _, a := range attrs {
		buf.WriteByte(' ')
		h.appendAttr(&buf, a, h.groups)
	}
	extra := buf.Bytes()
	pre := make([]byte, len(h.preformatted)+len(extra))
	copy(pre, h.preformatted)
	copy(pre[len(h.preformatted):], extra)
	gs := make([]string, len(h.groups))
	copy(gs, h.groups)
	return &handler{
		w:            h.w,
		opts:         h.opts,
		timeFormat:   h.timeFormat,
		preformatted: pre,
		groups:       gs,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	gs := make([]string, len(h.groups)+1)
	copy(gs, h.groups)
	gs[len(h.groups)] = name
	pre := make([]byte, len(h.preformatted))
	copy(pre, h.preformatted)
	return &handler{
		w:            h.w,
		opts:         h.opts,
		timeFormat:   h.timeFormat,
		preformatted: pre,
		groups:       gs,
	}
}

func (h *handler) appendAttr(buf *bytes.Buffer, a slog.Attr, groups []string) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	key := a.Key
	for i := len(groups) - 1; i >= 0; i-- {
		key = groups[i] + "." + key
	}
	if h.opts.NoColor {
		fmt.Fprintf(buf, "%s=%v", key, a.Value)
	} else {
		fmt.Fprintf(buf, "%s%s%s=%v", colorCyan, key, colorReset, a.Value)
	}
}

func levelInfo(l slog.Level) (string, string) {
	switch {
	case l < slog.LevelInfo:
		return "DBG", colorCyan
	case l < slog.LevelWarn:
		return "INF", colorGreen
	case l < slog.LevelError:
		return "WRN", colorYellow
	default:
		return "ERR", colorRed
	}
}
