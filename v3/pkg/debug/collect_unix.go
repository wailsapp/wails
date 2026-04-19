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
		MemoryBytes: 0,
		Goroutines: runtime.NumGoroutine(),
	}

	if pf, err := os.Open("/proc/self/status"); err == nil {
		defer pf.Close()
		scanner := bufio.NewScanner(pf)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "VmSize:") {
				if val := strings.TrimSpace(strings.TrimPrefix(line, "VmSize:")); strings.HasSuffix(val, "kB") {
					if v, err := strconv.ParseInt(strings.TrimSuffix(val, "kB"), 10, 64); err == nil {
						info.ProcessInfo.MemoryBytes = uint64(v * 1024)
					}
				}
			} else if strings.HasPrefix(line, "Threads:") {
				if val := strings.TrimSpace(strings.TrimPrefix(line, "Threads:")); v, err := strconv.Atoi(val); err == nil {
					info.ProcessInfo.Threads = v
				}
			}
		}
	}
}

func collectMemoryInfo(info *CrashInfo) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info.MemorySummary = MemorySummary{
		TotalVirtual: m.TotalAlloc + m.Mallocs,
		HeapAlloc:    m.HeapAlloc,
		HeapObjects:  m.HeapObjects,
		GarbageCollector: GCStats{
			NumGC:       int(m.NumGC),
			LastGC:      m.LastGC,
			TotalPause:  m.TotalPause,
			TotalAlloc:  m.TotalAlloc,
			HeapAlloc:   m.HeapAlloc,
			HeapObjects: m.HeapObjects,
		},
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

func writeCoreDump(path string, pid int) error {
	return fmt.Errorf("not implemented for %s", runtime.GOOS)
}