//go:build windows && !server

package application

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// forwardingInjections returns the worker-forwarding evaluations recorded for a
// session, in the order they were sent.
func forwardingInjections(f *cancellationTestCDP, session string) []int {
	var found []int
	for i, call := range f.calls {
		if call.session != session || call.method != "Runtime.evaluate" {
			continue
		}
		var params struct{ Expression string }
		if err := json.Unmarshal([]byte(call.params), &params); err != nil {
			continue
		}
		if params.Expression == workerForwardingScript {
			found = append(found, i)
		}
	}
	return found
}

// The forwarding shim must only be installed in worker sessions; the document
// session keeps its native Worker constructor.
func TestWorkerForwardingScriptInjectedIntoWorkerSessionsOnly(t *testing.T) {
	c, f := cancellationFixture(t)
	c.enableSession("", func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})
	if got := forwardingInjections(f, ""); len(got) != 0 {
		t.Fatalf("document session must not install the forwarding shim, got %d", len(got))
	}

	f.calls = nil
	c.enableSession("worker", func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})
	injected := forwardingInjections(f, "worker")
	if len(injected) != 1 {
		t.Fatalf("worker session must install the forwarding shim once, got %d", len(injected))
	}
	// The shim has to be installed before the startup pause is released,
	// otherwise the worker's first script can create its workers unforwarded.
	for _, call := range f.calls[injected[0]+1:] {
		if call.method == "Runtime.runIfWaitingForDebugger" {
			t.Fatal("forwarding shim installed after the worker was resumed")
		}
	}
	if last := f.calls[len(f.calls)-1]; last.method != "Target.setAutoAttach" {
		t.Fatalf("worker session setup ended with %q", last.method)
	}
}

// Installing the shim must not be able to fail the whole session; the commands
// after it still run.
func TestWorkerForwardingScriptIsBestEffort(t *testing.T) {
	c, f := cancellationFixture(t)
	f.calls = nil
	f.fail = "Runtime.evaluate"
	var done error
	c.configure("worker", []cancellationCommand{
		{method: "Network.enable", params: struct{}{}},
		{method: "Runtime.evaluate", params: map[string]string{"expression": workerForwardingScript}, optional: true},
		{method: "Target.setAutoAttach", params: struct{}{}},
	}, func(err error) { done = err })
	if done != nil {
		t.Fatalf("optional command failure must not fail the session: %v", done)
	}
	if last := f.calls[len(f.calls)-1]; last.method != "Target.setAutoAttach" {
		t.Fatalf("optional failure skipped the remaining commands, ended with %q", last.method)
	}
}

// A non-optional command failure must still fail the session, unchanged.
func TestWorkerForwardingRequiredCommandStillFails(t *testing.T) {
	c, f := cancellationFixture(t)
	f.calls = nil
	f.fail = "Network.enable"
	var done error
	c.configure("worker", []cancellationCommand{
		{method: "Network.enable", params: struct{}{}},
		{method: "Target.setAutoAttach", params: struct{}{}},
	}, func(err error) { done = err })
	if done == nil {
		t.Fatal("required command failure must fail the session")
	}
	if last := f.calls[len(f.calls)-1]; last.method != "Network.enable" {
		t.Fatalf("required failure kept configuring, ended with %q", last.method)
	}
}

// Guards the shim against accidental edits when node is unavailable.
func TestWorkerForwardingScriptContent(t *testing.T) {
	for _, fragment := range []string{
		"globalThis.__wailsWorkerForwarding",
		"globalThis.__wailsWorkerBase",
		"importScripts(",
		"NativeImportScripts",
		"await import(",
		"URL.createObjectURL",
		"URL.revokeObjectURL",
		"globalThis.fetch =",
		"globalThis.Request =",
	} {
		if !strings.Contains(workerForwardingScript, fragment) {
			t.Errorf("worker forwarding shim is missing %q", fragment)
		}
	}
}

// Exercises the forwarding decisions in node when it is available.
func TestWorkerForwardingScriptBehaviour(t *testing.T) {
	if _, err := os.Stat("worker_forwarding_windows_test.mjs"); err != nil {
		t.Skip("worker forwarding JavaScript test is unavailable")
	}
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for executable JavaScript shim tests")
	}
	out, err := exec.Command(node, "worker_forwarding_windows_test.mjs").CombinedOutput()
	if err != nil {
		t.Fatalf("worker forwarding shim behaviour: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "ok") {
		t.Fatalf("worker forwarding shim behaviour produced no result:\n%s", out)
	}
}
