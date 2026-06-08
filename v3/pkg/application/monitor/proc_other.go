//go:build !linux

package monitor

// readProc is a no-op on platforms without a /proc probe. The sampler still
// emits Go runtime stats (heap, goroutines); OS RSS/CPU are reported as zero.
// TODO(darwin/windows): task_info / GetProcessMemoryInfo probes.
func readProc() (procProbe, bool) {
	return procProbe{}, false
}
