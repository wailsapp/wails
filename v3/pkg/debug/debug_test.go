package debug

import (
	"runtime"
	"strings"
	"testing"
)

func TestReport_PopulatesBasicFields(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	if r == nil {
		t.Fatal("Report returned nil")
	}
	if r.Timestamp.IsZero() {
		t.Error("Timestamp not set")
	}
	if r.Build.GoVersion == "" {
		t.Error("Build.GoVersion empty")
	}
	if !strings.HasPrefix(r.Build.GoVersion, "go") {
		t.Errorf("Build.GoVersion = %q, expected to start with 'go'", r.Build.GoVersion)
	}
	if r.Crash == nil {
		t.Fatal("Crash missing")
	}
	if r.Crash.ProcessInfo.Goroutines <= 0 {
		t.Errorf("ProcessInfo.Goroutines = %d, expected > 0", r.Crash.ProcessInfo.Goroutines)
	}
	if r.Crash.ProcessInfo.PID <= 0 {
		t.Errorf("ProcessInfo.PID = %d, expected > 0", r.Crash.ProcessInfo.PID)
	}
	if r.Crash.Environment == nil || len(r.Crash.Environment) == 0 {
		t.Error("Environment empty")
	}
}

func TestReport_DumpPathEmptyByDefault(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	if r.DumpPath != "" {
		t.Errorf("DumpPath = %q, expected empty when WithDump not set", r.DumpPath)
	}
}

func TestDump_UnsupportedOnNonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows has a real implementation; see dump_windows_test.go")
	}
	path, err := Dump()
	if err == nil {
		t.Errorf("Dump on %s returned nil error; expected unimplemented error (got path=%q)", runtime.GOOS, path)
	}
	if path != "" {
		t.Errorf("Dump on %s returned path=%q; expected empty on failure", runtime.GOOS, path)
	}
}

func TestOptions_Composition(t *testing.T) {
	cfg := &dumpConfig{}
	WithPath("/tmp/foo.dmp")(cfg)
	WithFullMemory()(cfg)
	if cfg.path != "/tmp/foo.dmp" {
		t.Errorf("WithPath: got %q, want /tmp/foo.dmp", cfg.path)
	}
	if !cfg.fullMemory {
		t.Error("WithFullMemory: flag not set")
	}

	rcfg := &reportConfig{}
	WithDumpPath("/tmp/bar.dmp")(rcfg)
	if !rcfg.withDump {
		t.Error("WithDumpPath should imply WithDump")
	}
	if rcfg.dumpPath != "/tmp/bar.dmp" {
		t.Errorf("WithDumpPath: got %q, want /tmp/bar.dmp", rcfg.dumpPath)
	}

	rcfg2 := &reportConfig{}
	WithDumpFullMemory()(rcfg2)
	if !rcfg2.withDump {
		t.Error("WithDumpFullMemory should imply WithDump")
	}
	if !rcfg2.fullMemory {
		t.Error("WithDumpFullMemory: flag not set")
	}
}
