package debug

import (
	"encoding/json"
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

func TestReport_ProcessInfo(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	pi := r.Crash.ProcessInfo

	if pi.PID <= 0 {
		t.Errorf("PID = %d, want > 0", pi.PID)
	}
	if pi.Goroutines <= 0 {
		t.Errorf("Goroutines = %d, want > 0", pi.Goroutines)
	}
	if pi.MemoryBytes == 0 {
		t.Error("MemoryBytes = 0, expected non-zero")
	}
}

func TestReport_MemorySummary(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	ms := r.Crash.MemorySummary

	if ms.TotalVirtual == 0 {
		t.Error("TotalVirtual = 0, expected non-zero")
	}
	if ms.TotalWorkingSet == 0 {
		t.Error("TotalWorkingSet = 0, expected non-zero")
	}
	if ms.PrivateBytes == 0 {
		t.Error("PrivateBytes = 0, expected non-zero")
	}
	if ms.GarbageCollector.NumGC < 0 {
		t.Errorf("GCStats.NumGC = %d, want >= 0", ms.GarbageCollector.NumGC)
	}
}

func TestReport_LoadedModules(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("see dump_windows_test.go for Windows module tests")
	}

	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}

	if len(r.Crash.LoadedModules) == 0 {
		t.Error("LoadedModules empty, expected at least one module")
	}

	for _, mod := range r.Crash.LoadedModules {
		if mod.Path == "" {
			t.Error("ModuleInfo.Path empty")
		}
	}
}

func TestReport_Threads(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}

	if r.Crash.ProcessInfo.Threads <= 0 {
		t.Errorf("Threads = %d, expected > 0", r.Crash.ProcessInfo.Threads)
	}
}

func TestReport_DiagnosticsMayBeEmptyOnHealthySystem(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	// Diagnostics come from doctor-ng. On a healthy dev system where Go,
	// WebView2, and a package manager are all installed, the list will be
	// empty — that is correct, not a bug.
	// We only verify the field is non-nil so callers can safely range it.
	if r.Diagnostics == nil {
		t.Error("Diagnostics is nil, expected empty slice on healthy system")
	}
}

func TestReport_JSONRoundTrip(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("Marshalled JSON is empty")
	}

	var r2 CrashReport
	if err := json.Unmarshal(data, &r2); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if r2.Timestamp.IsZero() {
		t.Error("Timestamp lost in round-trip")
	}
	if r2.Crash == nil {
		t.Fatal("Crash lost in round-trip")
	}
	if r2.Crash.ProcessInfo.PID != r.Crash.ProcessInfo.PID {
		t.Errorf("PID mismatch: got %d, want %d", r2.Crash.ProcessInfo.PID, r.Crash.ProcessInfo.PID)
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

func TestThreadInfo_NoChannelField(t *testing.T) {
	ti := ThreadInfo{
		ID:         1,
		State:      "running",
		WaitReason: "",
	}
	data, err := json.Marshal(ti)
	if err != nil {
		t.Fatalf("Marshal ThreadInfo: %v", err)
	}
	if strings.Contains(string(data), "WaitChan") {
		t.Error("ThreadInfo JSON contains WaitChan — channel field should not be present")
	}
	if !strings.Contains(string(data), "wait_reason") {
		t.Error("ThreadInfo JSON missing wait_reason field")
	}
}
