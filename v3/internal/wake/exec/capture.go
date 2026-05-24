package exec

import (
	"bytes"

	"github.com/wailsapp/wails/v3/internal/report"
)

// captureWriter is the sink for a command's stdout and stderr. It splits the
// stream into lines and, for each line:
//
//   - if the line is a wake wire event (emitted by a subprocess producer such
//     as the bindings generator), routes it to the live reporter and drops it
//     from the visible/captured output;
//   - otherwise records it in a buffer (shown only if the command fails) and,
//     when streaming (verbose), forwards it to the reporter live.
//
// The same captureWriter instance is used for both stdout and stderr; os/exec
// guarantees serialised Write calls when the two are the same writer, so no
// locking is needed here.
type captureWriter struct {
	rep     report.Reporter
	stream  bool
	buf     bytes.Buffer
	partial bytes.Buffer
}

func newCaptureWriter(rep report.Reporter, stream bool) *captureWriter {
	return &captureWriter{rep: rep, stream: stream}
}

func (w *captureWriter) Write(p []byte) (int, error) {
	w.partial.Write(p)
	data := w.partial.Bytes()
	start := 0
	for {
		i := bytes.IndexByte(data[start:], '\n')
		if i < 0 {
			break
		}
		w.handleLine(string(data[start : start+i]))
		start += i + 1
	}
	rest := append([]byte(nil), data[start:]...)
	w.partial.Reset()
	w.partial.Write(rest)
	return len(p), nil
}

func (w *captureWriter) handleLine(line string) {
	if ev, ok := report.Decode(line); ok {
		report.Route(w.rep, ev)
		return
	}
	w.buf.WriteString(line)
	w.buf.WriteByte('\n')
	if w.stream {
		w.rep.StepOutput(line)
	}
}

// flush handles any trailing line that lacked a newline.
func (w *captureWriter) flush() {
	if w.partial.Len() > 0 {
		w.handleLine(w.partial.String())
		w.partial.Reset()
	}
}

func (w *captureWriter) output() string { return w.buf.String() }
