//go:build linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const samplerSupported = true

// WebKitGTK runs page content in a separate WebKitWebProcess, the analogue of
// macOS's com.apple.WebKit.WebContent.
const webContentProcName = "WebKitWebProcess"

// webContentPids enumerates every WebKitGTK web process on the machine.
// Callers diff against a pre-launch snapshot; summing all of them would fold
// in any other GTK/WebKit app running here.
func webContentPids() ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read /proc: %w", err)
	}

	var out []int
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue // not a pid directory
		}
		// /proc/<pid>/comm is capped at 15 chars (TASK_COMM_LEN-1), and
		// "WebKitWebProcess" is 16, so it always arrives truncated. Match on
		// argv[0] instead, which is not.
		raw, err := os.ReadFile(filepath.Join("/proc", e.Name(), "cmdline"))
		if err != nil || len(raw) == 0 {
			continue
		}
		argv0 := string(raw)
		if i := strings.IndexByte(argv0, 0); i >= 0 {
			argv0 = argv0[:i]
		}
		if filepath.Base(argv0) == webContentProcName {
			out = append(out, pid)
		}
	}
	return out, nil
}

// footprint reports Pss for pid — the closest Linux analogue to macOS
// phys_footprint, since it apportions shared pages rather than double-counting
// them. Falls back to VmRSS where smaps_rollup is unavailable.
//
// Caveat worth knowing when reading Linux results: memory that a process
// allocates, hands to another process via an fd, and then unmaps does not
// appear in Pss or RSS at all. macOS charges that to the owner and shows it as
// "owned unmapped"; Linux does not. So a memfd-based equivalent of the macOS
// retention would be invisible here — check open fd counts and Shmem in
// /proc/meminfo alongside these numbers.
func footprint(pid int) (uint64, bool) {
	if kb, ok := fieldKB(filepath.Join("/proc", strconv.Itoa(pid), "smaps_rollup"), "Pss:"); ok {
		return kb * 1024, true
	}
	if kb, ok := fieldKB(filepath.Join("/proc", strconv.Itoa(pid), "status"), "VmRSS:"); ok {
		return kb * 1024, true
	}
	return 0, false
}

func fieldKB(path, prefix string) (uint64, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}
	for _, line := range strings.Split(string(b), "\n") {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		f := strings.Fields(line)
		if len(f) < 2 {
			return 0, false
		}
		kb, err := strconv.ParseUint(f[1], 10, 64)
		if err != nil {
			return 0, false
		}
		return kb, true
	}
	return 0, false
}
