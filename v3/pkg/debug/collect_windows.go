//go:build windows

package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
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

// Minidump type flags (subset of MINIDUMP_TYPE from dbghelp.h).
// We default to a "rich but not huge" mix: thread info, handle data, and
// unloaded modules — enough for WinDbg/Visual Studio postmortem without
// writing the entire process address space to disk.
const (
	miniDumpNormal              = 0x00000000
	miniDumpWithFullMemory      = 0x00000002
	miniDumpWithHandleData      = 0x00000004
	miniDumpWithUnloadedModules = 0x00000020
	miniDumpWithThreadInfo      = 0x00001000

	defaultDumpFlags = miniDumpNormal |
		miniDumpWithThreadInfo |
		miniDumpWithHandleData |
		miniDumpWithUnloadedModules
)

// writeCoreDump writes a Windows minidump of the current process to path.
// If path is empty, <TempDir>/wails-crash-<pid>-<unix>.dmp is used.
// Returns the absolute path of the dump on success.
//
// Runs under the calling user's token — no elevation, no SeDebugPrivilege,
// no OpenProcess: dumping the current process uses the
// GetCurrentProcess() pseudo-handle which bypasses ACL checks.
func writeCoreDump(path string) (string, error) {
	if path == "" {
		name := fmt.Sprintf("wails-crash-%d-%d.dmp", os.Getpid(), time.Now().Unix())
		path = filepath.Join(os.TempDir(), name)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve dump path: %w", err)
	}

	pathUTF16, err := windows.UTF16PtrFromString(abs)
	if err != nil {
		return "", fmt.Errorf("encode dump path: %w", err)
	}
	hFile, err := windows.CreateFile(
		pathUTF16,
		windows.GENERIC_WRITE,
		0,
		nil,
		windows.CREATE_ALWAYS,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return "", fmt.Errorf("create dump file %s: %w", abs, err)
	}
	defer windows.CloseHandle(hFile)

	dbghelp, err := windows.LoadLibrary("Dbghelp.dll")
	if err != nil {
		return "", fmt.Errorf("load Dbghelp.dll: %w", err)
	}
	defer windows.FreeLibrary(dbghelp)

	proc, err := windows.GetProcAddress(dbghelp, "MiniDumpWriteDump")
	if err != nil {
		return "", fmt.Errorf("resolve MiniDumpWriteDump: %w", err)
	}

	// BOOL MiniDumpWriteDump(
	//   HANDLE hProcess,           // arg0: pseudo-handle to self
	//   DWORD  ProcessId,          // arg1: current pid
	//   HANDLE hFile,              // arg2: destination file
	//   MINIDUMP_TYPE DumpType,    // arg3: flags
	//   PMINIDUMP_EXCEPTION_INFORMATION ExceptionParam,      // arg4: nil
	//   PMINIDUMP_USER_STREAM_INFORMATION UserStreamParam,   // arg5: nil
	//   PMINIDUMP_CALLBACK_INFORMATION CallbackParam)        // arg6: nil
	ret, _, callErr := syscall.SyscallN(
		proc,
		uintptr(windows.CurrentProcess()),
		uintptr(windows.GetCurrentProcessId()),
		uintptr(hFile),
		uintptr(defaultDumpFlags),
		0, 0, 0,
	)
	if ret == 0 {
		return "", fmt.Errorf("MiniDumpWriteDump failed: %w", callErr)
	}
	return abs, nil
}