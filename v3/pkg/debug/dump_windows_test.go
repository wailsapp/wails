//go:build windows

package debug

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const miniDumpMagic uint32 = 0x504D444D

func TestDump_WritesValidMinidump(t *testing.T) {
	dumpPath := filepath.Join(t.TempDir(), "test.dmp")
	path, err := Dump(WithPath(dumpPath))
	if err != nil {
		t.Fatalf("Dump: %v", err)
	}
	if path == "" {
		t.Fatal("Dump returned empty path")
	}
	if !strings.EqualFold(path, dumpPath) {
		t.Errorf("Dump path = %q, want %q", path, dumpPath)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat dump: %v", err)
	}
	if info.Size() < 1024 {
		t.Errorf("dump size = %d bytes, expected at least 1KB", info.Size())
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open dump: %v", err)
	}
	defer f.Close()

	var magic uint32
	if err := binary.Read(f, binary.LittleEndian, &magic); err != nil {
		t.Fatalf("read magic: %v", err)
	}
	if magic != miniDumpMagic {
		t.Errorf("magic = %#x, want %#x (MDMP)", magic, miniDumpMagic)
	}
}

func TestDump_DefaultPathInTempDir(t *testing.T) {
	path, err := Dump()
	if err != nil {
		t.Fatalf("Dump: %v", err)
	}
	defer os.Remove(path)

	if !strings.HasPrefix(path, os.TempDir()) {
		t.Errorf("Dump default path %q not under TempDir %q", path, os.TempDir())
	}
	if !strings.HasSuffix(path, ".dmp") {
		t.Errorf("Dump path %q missing .dmp suffix", path)
	}
}

func TestReport_WithDump_SetsDumpPath(t *testing.T) {
	dumpPath := filepath.Join(t.TempDir(), "report.dmp")
	r, err := Report(WithDumpPath(dumpPath))
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	if r.DumpPath == "" {
		t.Error("CrashReport.DumpPath empty after WithDumpPath")
	}
	if _, err := os.Stat(r.DumpPath); err != nil {
		t.Errorf("dump not at reported path %q: %v", r.DumpPath, err)
	}
}

func TestReport_Windows_ProcessInfo(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	pi := r.Crash.ProcessInfo

	if pi.Handles <= 0 {
		t.Errorf("Handles = %d, expected > 0", pi.Handles)
	}
	if pi.MemoryBytes == 0 {
		t.Error("MemoryBytes = 0, expected non-zero from GetProcessMemoryInfo")
	}
}

func TestReport_Windows_MemoryFromOS(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	ms := r.Crash.MemorySummary

	if ms.TotalVirtual == 0 {
		t.Error("TotalVirtual = 0, expected PagefileUsage from GetProcessMemoryInfo")
	}
	if ms.TotalWorkingSet == 0 {
		t.Error("TotalWorkingSet = 0, expected WorkingSetSize from GetProcessMemoryInfo")
	}
	if ms.PrivateBytes == 0 {
		t.Error("PrivateBytes = 0, expected PrivateUsage from GetProcessMemoryInfo")
	}
}

func TestReport_Windows_LoadedModules(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}

	if len(r.Crash.LoadedModules) == 0 {
		t.Fatal("LoadedModules empty, expected at least the executable module")
	}

	found := false
	for _, mod := range r.Crash.LoadedModules {
		if mod.Name == "" {
			t.Error("ModuleInfo.Name empty")
		}
		if mod.Path == "" {
			t.Error("ModuleInfo.Path empty")
		}
		if mod.Size == 0 {
			t.Errorf("ModuleInfo.Size = 0 for %s", mod.Name)
		}
		if !mod.Loaded {
			t.Errorf("ModuleInfo.Loaded = false for %s", mod.Name)
		}
		if strings.EqualFold(mod.Name, "ntdll.dll") {
			found = true
		}
	}
	if !found {
		t.Error("ntdll.dll not found in LoadedModules — module enumeration may be broken")
	}
}

func TestReport_Windows_ThreadCount(t *testing.T) {
	r, err := Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	if r.Crash.ProcessInfo.Threads <= 0 {
		t.Errorf("Threads = %d, expected > 0 from CreateToolhelp32Snapshot", r.Crash.ProcessInfo.Threads)
	}
}

func TestReport_Windows_JSONRoundTrip(t *testing.T) {
	r, err := Report(WithDump())
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	defer func() {
		if r.DumpPath != "" {
			os.Remove(r.DumpPath)
		}
	}()

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var r2 CrashReport
	if err := json.Unmarshal(data, &r2); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(r2.Crash.LoadedModules) != len(r.Crash.LoadedModules) {
		t.Errorf("LoadedModules count mismatch: got %d, want %d", len(r2.Crash.LoadedModules), len(r.Crash.LoadedModules))
	}
	if r2.Crash.ProcessInfo.Handles != r.Crash.ProcessInfo.Handles {
		t.Errorf("Handles mismatch: got %d, want %d", r2.Crash.ProcessInfo.Handles, r.Crash.ProcessInfo.Handles)
	}
}
