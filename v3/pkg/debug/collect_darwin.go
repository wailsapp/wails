//go:build darwin

package debug

import (
	"fmt"
	"os"
	"runtime"
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
}

func collectMemoryInfo(info *CrashInfo) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info.MemorySummary = MemorySummary{
		TotalVirtual:    m.HeapInuse + m.StackInuse,
		TotalWorkingSet: m.HeapAlloc,
		PrivateBytes:    m.HeapInuse,
		GarbageCollector: GCStats{
			NumGC:       int(m.NumGC),
			LastGC:      uint32(m.LastGC),
			TotalPause:  m.PauseTotalNs,
			TotalAlloc:  m.TotalAlloc,
			HeapAlloc:   m.HeapAlloc,
			HeapObjects: m.HeapObjects,
		},
	}
}

func collectThreadInfo(info *CrashInfo) {
}

func loadModules(info *CrashInfo) {
}

func writeCoreDump(path string, fullMemory bool) (string, error) {
	return "", fmt.Errorf("core dump not implemented for %s", runtime.GOOS)
}
