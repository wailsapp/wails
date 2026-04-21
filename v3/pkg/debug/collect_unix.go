//go:build !windows

package debug

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func collectEnvironmentVars() map[string]string {
	env := make(map[string]string)
	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}

func collectProcessInfo(info *CrashInfo) {
	info.ProcessInfo = ProcessInfo{
		PID:        os.Getpid(),
		Goroutines: runtime.NumGoroutine(),
	}

	if pf, err := os.Open("/proc/self/status"); err == nil {
		defer pf.Close()
		scanner := bufio.NewScanner(pf)
		for scanner.Scan() {
			line := scanner.Text()
			switch {
			case strings.HasPrefix(line, "VmSize:"):
				if v := parseProcKB(line); v > 0 {
					info.ProcessInfo.MemoryBytes = v
				}
			case strings.HasPrefix(line, "VmRSS:"):
				if v := parseProcKB(line); v > 0 {
					info.ProcessInfo.MemoryBytes = v
				}
			case strings.HasPrefix(line, "Threads:"):
				val := strings.TrimSpace(strings.TrimPrefix(line, "Threads:"))
				if v, err := strconv.Atoi(val); err == nil {
					info.ProcessInfo.Threads = v
				}
			}
		}
	}
}

func parseProcKB(line string) uint64 {
	val := strings.TrimSpace(strings.TrimPrefix(line, strings.SplitN(line, ":", 2)[0]+":"))
	val = strings.TrimSuffix(val, " kB")
	v, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return v * 1024
}

func collectMemoryInfo(info *CrashInfo) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info.MemorySummary = MemorySummary{
		PrivateBytes: m.HeapInuse,
		GarbageCollector: GCStats{
			NumGC:       int(m.NumGC),
			LastGC:      uint32(m.LastGC),
			TotalPause:  m.PauseTotalNs,
			TotalAlloc:  m.TotalAlloc,
			HeapAlloc:   m.HeapAlloc,
			HeapObjects: m.HeapObjects,
		},
	}

	if pf, err := os.Open("/proc/self/status"); err == nil {
		defer pf.Close()
		scanner := bufio.NewScanner(pf)
		for scanner.Scan() {
			line := scanner.Text()
			switch {
			case strings.HasPrefix(line, "VmSize:"):
				if v := parseProcKB(line); v > 0 {
					info.MemorySummary.TotalVirtual = v
				}
			case strings.HasPrefix(line, "VmRSS:"):
				if v := parseProcKB(line); v > 0 {
					info.MemorySummary.TotalWorkingSet = v
				}
			case strings.HasPrefix(line, "VmData:"):
				if v := parseProcKB(line); v > 0 {
					info.MemorySummary.PrivateBytes = v
				}
			}
		}
	}
}

func collectThreadInfo(info *CrashInfo) {
}

func loadModules(info *CrashInfo) {
	if mf, err := os.Open("/proc/self/maps"); err == nil {
		defer mf.Close()
		scanner := bufio.NewScanner(mf)
		seen := make(map[string]bool)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) >= 6 && !seen[parts[len(parts)-1]] {
				seen[parts[len(parts)-1]] = true
				info.LoadedModules = append(info.LoadedModules, ModuleInfo{
					Path: parts[len(parts)-1],
				})
			}
		}
	}
}

func collectNetworkConnections(info *CrashInfo) {
}

// writeCoreDump is Windows-only for now. On other platforms it returns an
// error so callers can fall back to stack-trace-only reporting.
//
//nolint:unused // fullMemory kept in signature for parity with Windows impl
func writeCoreDump(path string, fullMemory bool) (string, error) {
	return "", fmt.Errorf("core dump not implemented for %s", runtime.GOOS)
}
