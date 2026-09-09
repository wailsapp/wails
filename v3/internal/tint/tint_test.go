package tint

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// logRecord builds a slog.Record directly so we can control the timestamp.
func logRecord(t time.Time, level slog.Level, msg string) slog.Record {
	return slog.NewRecord(t, level, msg, 0)
}

// ---- NewHandler ----

func TestNewHandler_NilOpts(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestNewHandler_EmptyOpts(t *testing.T) {
	var buf bytes.Buffer
	// &Options{} has nil Level and empty TimeFormat — exercises both defaults
	h := NewHandler(&buf, &Options{})
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestNewHandler_ExplicitOpts(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05",
		NoColor:    false,
	})
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

// ---- Enabled ----

func TestEnabled_RespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelWarn})
	ctx := context.Background()

	if h.Enabled(ctx, slog.LevelDebug) {
		t.Error("Debug should be disabled when min level is Warn")
	}
	if h.Enabled(ctx, slog.LevelInfo) {
		t.Error("Info should be disabled when min level is Warn")
	}
	if !h.Enabled(ctx, slog.LevelWarn) {
		t.Error("Warn should be enabled when min level is Warn")
	}
	if !h.Enabled(ctx, slog.LevelError) {
		t.Error("Error should be enabled when min level is Warn")
	}
}

// ---- Handle: all four levels ----

func TestHandle_LevelAbbreviations(t *testing.T) {
	cases := []struct {
		level slog.Level
		abbr  string
	}{
		{slog.LevelDebug, "DBG"},
		{slog.LevelInfo, "INF"},
		{slog.LevelWarn, "WRN"},
		{slog.LevelError, "ERR"},
	}
	ctx := context.Background()
	for _, tc := range cases {
		t.Run(tc.abbr, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHandler(&buf, &Options{Level: slog.LevelDebug})
			r := logRecord(time.Now(), tc.level, "hello")
			if err := h.Handle(ctx, r); err != nil {
				t.Fatalf("Handle: %v", err)
			}
			out := buf.String()
			if !strings.Contains(out, tc.abbr) {
				t.Errorf("expected %q in output %q", tc.abbr, out)
			}
			if !strings.Contains(out, "hello") {
				t.Errorf("expected message 'hello' in output %q", out)
			}
		})
	}
}

// ---- Handle: no time (zero time) ----

func TestHandle_ZeroTime(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug})
	r := logRecord(time.Time{}, slog.LevelInfo, "no-time")
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "INF") {
		t.Errorf("expected INF in output %q", out)
	}
}

// ---- Handle: TimeFormat option ----

func TestHandle_TimeFormat(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{
		Level:      slog.LevelDebug,
		TimeFormat: "2006-01-02",
	})
	fixed := time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC)
	r := logRecord(fixed, slog.LevelInfo, "formatted")
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "2024-03-14") {
		t.Errorf("expected date string in output %q", out)
	}
}

// ---- NoColor ----

func TestHandle_NoColor_ProducesNoEscapeCodes(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true})
	r := logRecord(time.Now(), slog.LevelInfo, "plain")
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI escape codes in NoColor output, got %q", out)
	}
}

func TestHandle_Color_ContainsEscapeCodes(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: false})
	r := logRecord(time.Now(), slog.LevelInfo, "colorful")
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in colored output, got %q", out)
	}
}

// ---- WithAttrs ----

func TestWithAttrs_AddsAttrsToSubsequentRecords(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true})
	h2 := h.WithAttrs([]slog.Attr{slog.String("key", "val")})

	r := logRecord(time.Time{}, slog.LevelInfo, "msg")
	if err := h2.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "key=val") {
		t.Errorf("expected key=val in output %q", out)
	}
}

func TestWithAttrs_ParentUnchanged(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true}).(*handler)
	origLen := len(h.preformatted)

	child := h.WithAttrs([]slog.Attr{slog.String("x", "1")}).(*handler)
	// child adds attrs
	if len(child.preformatted) <= origLen {
		t.Error("child should have more preformatted bytes than parent")
	}
	// parent slice not extended
	if len(h.preformatted) != origLen {
		t.Error("parent preformatted slice was modified")
	}
}

// ---- WithGroup ----

func TestWithGroup_PrefixesKeys(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true})
	h2 := h.WithGroup("grp")

	r := logRecord(time.Time{}, slog.LevelInfo, "msg")
	r.AddAttrs(slog.String("field", "v"))
	if err := h2.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "grp.field=v") {
		t.Errorf("expected 'grp.field=v' in output %q", out)
	}
}

func TestWithGroup_ParentGroupsUnchanged(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug}).(*handler)
	origGroups := len(h.groups)

	child := h.WithGroup("g1").(*handler)
	if len(child.groups) != origGroups+1 {
		t.Errorf("child should have one more group, got %d", len(child.groups))
	}
	if len(h.groups) != origGroups {
		t.Error("parent groups slice was modified")
	}
}

// ---- WithAttrs + WithGroup don't alias ----

func TestWithAttrsAndGroup_NoAlias(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true}).(*handler)

	// Add group then attrs
	hg := h.WithGroup("g").(*handler)
	ha := hg.WithAttrs([]slog.Attr{slog.String("a", "1")}).(*handler)

	// Modify child — parent should be unaffected
	_ = ha.WithAttrs([]slog.Attr{slog.String("b", "2")})

	// hg.preformatted should not contain "a=1" (ha's attr)
	if strings.Contains(string(hg.preformatted), "a=1") {
		t.Error("hg.preformatted was contaminated by ha's WithAttrs")
	}
}

// ---- appendAttr: empty attr is skipped ----

func TestHandle_EmptyAttrSkipped(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true})
	r := logRecord(time.Time{}, slog.LevelInfo, "msg")
	// Add a zero Attr — should be silently dropped
	r.AddAttrs(slog.Attr{})
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	// output should just be "INF msg\n"
	if strings.Contains(out, "=") {
		t.Errorf("empty attr should not produce key=val output, got %q", out)
	}
}

// ---- appendAttr with NoColor on an attr ----

func TestHandle_AttrNoColor(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true})
	r := logRecord(time.Time{}, slog.LevelInfo, "msg")
	r.AddAttrs(slog.Int("n", 42))
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "=42") {
		t.Errorf("expected =42 in output %q", out)
	}
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI codes in NoColor output, got %q", out)
	}
}

func TestHandle_AttrColor(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: false})
	r := logRecord(time.Time{}, slog.LevelInfo, "msg")
	r.AddAttrs(slog.Int("n", 42))
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "=42") {
		t.Errorf("expected =42 in output %q", out)
	}
	// attr key should be colored with colorCyan
	if !strings.Contains(out, colorCyan) {
		t.Errorf("expected colorCyan escape in colored attr output, got %q", out)
	}
}

// ---- levelInfo covers all four branches ----

func TestLevelInfo_AllBranches(t *testing.T) {
	cases := []struct {
		level  slog.Level
		abbr   string
		color  string
	}{
		{slog.LevelDebug, "DBG", colorCyan},
		{slog.LevelInfo, "INF", colorGreen},
		{slog.LevelWarn, "WRN", colorYellow},
		{slog.LevelError, "ERR", colorRed},
	}
	for _, tc := range cases {
		abbr, color := levelInfo(tc.level)
		if abbr != tc.abbr {
			t.Errorf("levelInfo(%v): got abbr %q, want %q", tc.level, abbr, tc.abbr)
		}
		if color != tc.color {
			t.Errorf("levelInfo(%v): got color %q, want %q", tc.level, color, tc.color)
		}
	}
}

// ---- WithAttrs with preformatted in color mode ----

func TestWithAttrs_ColorMode(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: false})
	h2 := h.WithAttrs([]slog.Attr{slog.String("env", "prod")})

	r := logRecord(time.Time{}, slog.LevelInfo, "msg")
	if err := h2.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	// The key "env" is wrapped in color codes, but "=prod" should appear verbatim
	if !strings.Contains(out, "=prod") {
		t.Errorf("expected =prod in output %q", out)
	}
}

// ---- Timestamp in color mode ----

func TestHandle_TimestampColor(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: false})
	fixed := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	r := logRecord(fixed, slog.LevelInfo, "ts-color")
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	// Both timestamp and ANSI codes should be present
	if !strings.Contains(out, "Jan") {
		t.Errorf("expected timestamp in output %q", out)
	}
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI codes in colored timestamp output %q", out)
	}
}

// ---- Timestamp in NoColor mode ----

func TestHandle_TimestampNoColor(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true})
	fixed := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	r := logRecord(fixed, slog.LevelInfo, "ts-nocolor")
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Jan") {
		t.Errorf("expected timestamp in output %q", out)
	}
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI codes in NoColor timestamp output %q", out)
	}
}

// ---- Message in NoColor mode ----

func TestHandle_MessageNoColor(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, &Options{Level: slog.LevelDebug, NoColor: true})
	r := logRecord(time.Time{}, slog.LevelInfo, "plain-message")
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "plain-message") {
		t.Errorf("expected message in output %q", out)
	}
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI codes in NoColor message output %q", out)
	}
}
