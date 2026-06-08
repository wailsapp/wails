//go:build linux

package monitor

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// clockTicks is the kernel USER_HZ. It is effectively always 100 on Linux; we
// hard-code it to avoid a cgo sysconf call. CPU time in /proc is in these ticks.
const clockTicks = 100

// readProc reads /proc/self for RSS, CPU time, thread and fd counts.
func readProc() (procProbe, bool) {
	var p procProbe
	pageSize := uint64(os.Getpagesize())

	// /proc/self/statm: fields are pages; field 2 (index 1) is resident.
	if data, err := os.ReadFile("/proc/self/statm"); err == nil {
		f := strings.Fields(string(data))
		if len(f) >= 2 {
			if res, err := strconv.ParseUint(f[1], 10, 64); err == nil {
				p.RSS = res * pageSize
			}
		}
	}

	// /proc/self/stat: utime (14), stime (15), num_threads (20). 1-indexed per
	// proc(5). The comm field (2) can contain spaces/parens, so split after the
	// trailing ')'.
	if data, err := os.ReadFile("/proc/self/stat"); err == nil {
		s := string(data)
		if i := strings.LastIndexByte(s, ')'); i >= 0 && i+2 < len(s) {
			f := strings.Fields(s[i+2:]) // f[0] == state, i.e. field 3 onward
			// utime is field 14 -> index 11 here; stime field 15 -> index 12;
			// num_threads field 20 -> index 17.
			if len(f) >= 18 {
				utime, _ := strconv.ParseUint(f[11], 10, 64)
				stime, _ := strconv.ParseUint(f[12], 10, 64)
				ticks := utime + stime
				p.CPUTime = time.Duration(ticks) * time.Second / clockTicks
				if n, err := strconv.Atoi(f[17]); err == nil {
					p.Threads = n
				}
			}
		}
	}

	// Open file descriptors.
	if entries, err := os.ReadDir("/proc/self/fd"); err == nil {
		p.FDs = len(entries)
	}

	return p, p.RSS > 0 || p.CPUTime > 0
}
