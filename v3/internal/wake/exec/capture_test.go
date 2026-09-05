package exec

import (
	"testing"

	"github.com/wailsapp/wails/v3/internal/report"
)

type recordingReporter struct {
	report.Nop
	infos   []string
	outputs []string
}

func (r *recordingReporter) StepInfo(_ report.StepID, msg string)    { r.infos = append(r.infos, msg) }
func (r *recordingReporter) StepOutput(_ report.StepID, line string) { r.outputs = append(r.outputs, line) }

func TestCaptureWriterStripsWireEventsAndCaptures(t *testing.T) {
	rep := &recordingReporter{}
	cw := newCaptureWriter(rep, 1, true) // stream=true so StepOutput fires

	ev := report.Encode(report.Event{Kind: report.KindStatus, Msg: "working"})
	// Two ordinary lines around a wire event, written in fragments to exercise
	// the partial-line buffering.
	_, _ = cw.Write([]byte("line one\n" + ev + "\nline "))
	_, _ = cw.Write([]byte("two\n"))
	cw.flush()

	if got := cw.output(); got != "line one\nline two\n" {
		t.Fatalf("captured output = %q, want wire event stripped", got)
	}
	if len(rep.infos) != 1 || rep.infos[0] != "working" {
		t.Fatalf("wire event not routed to StepInfo: %v", rep.infos)
	}
	want := []string{"line one", "line two"}
	if len(rep.outputs) != 2 || rep.outputs[0] != want[0] || rep.outputs[1] != want[1] {
		t.Fatalf("streamed output = %v, want %v", rep.outputs, want)
	}
}

func TestCaptureWriterFlushesTrailingLine(t *testing.T) {
	rep := &recordingReporter{}
	cw := newCaptureWriter(rep, 1, false)
	_, _ = cw.Write([]byte("no newline here"))
	if cw.output() != "" {
		t.Fatalf("partial line should not be captured before flush: %q", cw.output())
	}
	cw.flush()
	if cw.output() != "no newline here\n" {
		t.Fatalf("flush did not emit trailing line: %q", cw.output())
	}
}
