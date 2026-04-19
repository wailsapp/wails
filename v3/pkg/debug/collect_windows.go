//go:build windows

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
		TotalVirtual: m.TotalAlloc + m.Mallocs,
		HeapAlloc:   m.HeapAlloc,
		HeapObjects: m.HeapObjects,
	}
}

func collectThreadInfo(info *CrashInfo) {
}

func loadModules(info *CrashInfo) {
}

func writeCoreDump(path string, pid int) error {
	return fmt.Errorf("not implemented - requires golang.org/x/sys/windows")
}