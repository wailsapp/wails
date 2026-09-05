package report

import (
	"strings"
	"testing"
	"time"
)

func TestEncodeDecodeRoundtrip(t *testing.T) {
	in := Event{Kind: KindStatus, Msg: "generating bindings"}
	line := Encode(in)
	if !strings.HasPrefix(line, sentinel) {
		t.Fatalf("encoded line missing sentinel: %q", line)
	}
	if strings.Contains(line, "\n") {
		t.Fatalf("encoded line must be single-line: %q", line)
	}
	out, ok := Decode(line)
	if !ok {
		t.Fatal("Decode failed on encoded line")
	}
	if out != in {
		t.Fatalf("roundtrip mismatch: got %+v want %+v", out, in)
	}
}

func TestDecodeIgnoresOrdinaryOutput(t *testing.T) {
	for _, line := range []string{
		"building frontend...",
		"# github.com/example/app",
		"",
		"\x1fwake\x1f not json",
	} {
		if _, ok := Decode(line); ok {
			t.Errorf("Decode wrongly accepted ordinary line %q", line)
		}
	}
}

func TestEmitWritesWireLine(t *testing.T) {
	var b strings.Builder
	Emit(&b, Event{Kind: KindInfo, Msg: "hi"})
	got := b.String()
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("Emit must terminate with newline: %q", got)
	}
	ev, ok := Decode(strings.TrimRight(got, "\n"))
	if !ok || ev.Msg != "hi" {
		t.Fatalf("emitted line did not decode: %q", got)
	}
}

type capture struct {
	Nop
	infos []string
}

func (c *capture) StepInfo(_ StepID, msg string) { c.infos = append(c.infos, msg) }

func TestRouteSendsEventToStepInfo(t *testing.T) {
	c := &capture{}
	Route(c, 1, Event{Kind: KindStatus, Msg: "phase 1"})
	Route(c, 1, Event{Kind: KindError, Msg: "oops"})
	if len(c.infos) != 2 || c.infos[0] != "phase 1" || c.infos[1] != "oops" {
		t.Fatalf("Route did not forward events: %v", c.infos)
	}
}

// Nop must satisfy Reporter and be safe to call.
func TestNopIsReporter(t *testing.T) {
	var r Reporter = Nop{}
	r.BuildStart("build", "t", 1)
	id := r.StepStart("t", "")
	r.StepInfo(id, "x")
	r.StepCommand(id, "c")
	r.StepOutput(id, "o")
	r.StepEnd(id, StatusOK, time.Second)
	r.StepFailed(id, Failure{})
	r.BuildEnd(time.Second, true)
	if r.Level() != Normal {
		t.Fatalf("Nop level = %v", r.Level())
	}
}
