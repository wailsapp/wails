//go:build windows

package debug

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// minidump files start with the 4-byte signature "MDMP" (0x504D444D LE).
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
