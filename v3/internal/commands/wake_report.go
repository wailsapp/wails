package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v3/internal/report"
)

// underWake reports whether this process was spawned by a wake build. When true,
// long-running producers (the bindings generator) should emit structured wire
// events for the build UI instead of drawing their own spinner/headers.
func underWake() bool {
	return os.Getenv("WAKE_REPORT") == "1"
}

// wakeLogger is a generator/config.Logger that forwards messages to the parent
// wake build as wire events on stdout, where the build's output capture decodes
// them into live sub-lines under the current step.
type wakeLogger struct {
	w io.Writer
}

func newWakeLogger() wakeLogger { return wakeLogger{w: os.Stdout} }

func (l wakeLogger) emit(kind report.EventKind, format string, a ...any) {
	report.Emit(l.w, report.Event{Kind: kind, Msg: fmt.Sprintf(format, a...)})
}

func (l wakeLogger) Errorf(format string, a ...any)   { l.emit(report.KindError, format, a...) }
func (l wakeLogger) Warningf(format string, a ...any) { l.emit(report.KindWarn, format, a...) }
func (l wakeLogger) Infof(format string, a ...any)    { l.emit(report.KindInfo, format, a...) }
func (l wakeLogger) Statusf(format string, a ...any)  { l.emit(report.KindStatus, format, a...) }

// Debugf is dropped: debug chatter would flood the build UI. It still reaches a
// human running the generator directly (non-wake path).
func (l wakeLogger) Debugf(format string, a ...any) {}
