package dev

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
)

type testProcess struct{ done chan struct{} }

func (p testProcess) Done() <-chan struct{}            { return p.done }
func (p testProcess) Stop(time.Duration)               {}
func (p testProcess) Err() error                       { return nil }
func (p testProcess) NeedsStopBeforeReplacement() bool { return false }
func TestReadinessRequiresLaunchToken(t *testing.T) {
	r, e := NewReadiness()
	if e != nil {
		t.Fatal(e)
	}
	defer r.Close()
	p := testProcess{make(chan struct{})}
	req, _ := http.NewRequest("POST", r.address, nil)
	response, e := http.DefaultClient.Do(req)
	if e != nil {
		t.Fatal(e)
	}
	response.Body.Close()
	if response.StatusCode != 401 {
		t.Fatal(response.StatusCode)
	}
	select {
	case <-r.ready:
		t.Fatal("unauthenticated launch marked ready")
	default:
	}
	req, _ = http.NewRequest("POST", r.address, nil)
	req.Header.Set("Authorization", "Bearer "+r.token)
	response, e = http.DefaultClient.Do(req)
	if e != nil {
		t.Fatal(e)
	}
	response.Body.Close()
	if e = r.Wait(t.Context(), p, time.Second); e != nil {
		t.Fatal(e)
	}
}
func TestReadinessCancellationExitAndTimeout(t *testing.T) {
	for _, kind := range []string{"cancel", "exit", "timeout"} {
		t.Run(kind, func(t *testing.T) {
			r, e := NewReadiness()
			if e != nil {
				t.Fatal(e)
			}
			defer r.Close()
			p := testProcess{make(chan struct{})}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if kind == "cancel" {
				cancel()
			}
			if kind == "exit" {
				close(p.done)
			}
			e = r.Wait(ctx, p, time.Millisecond)
			if e == nil {
				t.Fatal("false readiness")
			}
			if kind == "cancel" && !errors.Is(e, context.Canceled) {
				t.Fatal(e)
			}
		})
	}
}
func TestOutputConcurrentStreamsPreserveLinesAndTail(t *testing.T) {
	var buffer bytes.Buffer
	o := NewOutput(&buffer, report.Normal)
	var wg sync.WaitGroup
	for _, name := range []string{"frontend", "backend"} {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			w := o.Writer(name, "stderr")
			defer w.Close()
			io.WriteString(w, "stack\n  frame\ntail")
		}(name)
	}
	wg.Wait()
	for _, name := range []string{"frontend", "backend"} {
		for _, line := range []string{"stack", "  frame", "tail"} {
			if !strings.Contains(buffer.String(), "["+name+"/stderr] "+line+"\n") {
				t.Fatal(buffer.String())
			}
		}
	}
}
func TestOutputBoundsNewlineFreeWrites(t *testing.T) {
	var b bytes.Buffer
	o := NewOutput(&b, report.Normal)
	w := o.Writer("backend", "stdout").(*lineWriter)
	w.Write(bytes.Repeat([]byte{'x'}, 200000))
	if len(w.pending) >= 65536 {
		t.Fatal(len(w.pending))
	}
	w.Close()
	if strings.Count(b.String(), "x") != 200000 {
		t.Fatal("lost output")
	}
}
func TestSessionOutputsAreIndependent(t *testing.T) {
	var a, b bytes.Buffer
	oa, ob := NewOutput(&a, report.Normal), NewOutput(&b, report.Normal)
	ca, cb := WithOutput(context.Background(), oa), WithOutput(context.Background(), ob)
	OutputFrom(ca).BuildReporter().BuildStart("build", "one", 1)
	OutputFrom(cb).BuildReporter().BuildStart("build", "two", 1)
	if strings.Contains(a.String(), "two") || strings.Contains(b.String(), "one") {
		t.Fatal("cross-session output")
	}
}

func TestFrontendReadinessWaitsThroughServerErrors(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(503)
		} else {
			w.WriteHeader(200)
		}
	}))
	defer server.Close()
	if err := WaitHTTP(t.Context(), testProcess{make(chan struct{})}, server.URL, time.Second); err != nil {
		t.Fatal(err)
	}
	if calls.Load() < 2 {
		t.Fatal("accepted an unavailable frontend")
	}
}
func TestFrontendReadinessRejectsUntrustedTLS(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
	defer server.Close()
	if err := WaitHTTP(t.Context(), testProcess{make(chan struct{})}, server.URL, 100*time.Millisecond); err == nil {
		t.Fatal("accepted untrusted TLS")
	}
}

func TestCanceledSessionDoesNotReportAnExitFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sessionProcessExit(ctx, "backend", testProcess{make(chan struct{})}); err != nil {
		t.Fatal(err)
	}
}
func TestReportedRebuildFailureIsNotRenderedTwice(t *testing.T) {
	var b bytes.Buffer
	o := NewOutput(&b, report.Normal)
	o.Failure(3, "build", "rebuild", report.MarkReported(errors.New("compiler diagnostic already rendered")), "Current app is still running.")
	if strings.Contains(b.String(), "compiler diagnostic") {
		t.Fatal(b.String())
	}
	if !strings.Contains(b.String(), "Current app is still running.") {
		t.Fatal(b.String())
	}
}
func TestTailKeepsTheLastBytes(t *testing.T) {
	var tail Tail
	tail.Write(bytes.Repeat([]byte{'a'}, 70000))
	tail.Write([]byte("last"))
	if len(tail.String()) != 65536 || !strings.HasSuffix(tail.String(), "last") {
		t.Fatal("incorrect diagnostic tail")
	}
}
