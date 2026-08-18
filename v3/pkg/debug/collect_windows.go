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
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	psapi    = windows.NewLazyDLL("psapi.dll")
	kernel32 = windows.NewLazyDLL("kernel32.dll")

	procGetProcessMemoryInfo  = psapi.NewProc("GetProcessMemoryInfo")
	procGetProcessHandleCount = kernel32.NewProc("GetProcessHandleCount")
)

type processMemoryCountersEx struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
	PrivateUsage               uintptr
}

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

	h := windows.CurrentProcess()

	var memCounters processMemoryCountersEx
	memCounters.CB = uint32(unsafe.Sizeof(memCounters))
	ret, _, _ := procGetProcessMemoryInfo.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&memCounters)),
		uintptr(memCounters.CB),
	)
	if ret != 0 {
		info.ProcessInfo.MemoryBytes = uint64(memCounters.WorkingSetSize)
	}

	var handleCount uint32
	ret, _, _ = procGetProcessHandleCount.Call(uintptr(h), uintptr(unsafe.Pointer(&handleCount)))
	if ret != 0 {
		info.ProcessInfo.Handles = int(handleCount)
	}

	var creationTime, exitTime, kernelTime, userTime windows.Filetime
	if err := windows.GetProcessTimes(h, &creationTime, &exitTime, &kernelTime, &userTime); err == nil {
		info.ProcessInfo.StartTime = time.Unix(0, creationTime.Nanoseconds())
		info.ProcessInfo.UserTime = float64(userTime.Nanoseconds()) / 1e6
		info.ProcessInfo.SystemTime = float64(kernelTime.Nanoseconds()) / 1e6
	}
}

func collectMemoryInfo(info *CrashInfo) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info.MemorySummary = MemorySummary{
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

	h := windows.CurrentProcess()
	var memCounters processMemoryCountersEx
	memCounters.CB = uint32(unsafe.Sizeof(memCounters))
	ret, _, _ := procGetProcessMemoryInfo.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&memCounters)),
		uintptr(memCounters.CB),
	)
	if ret != 0 {
		info.MemorySummary.TotalVirtual = uint64(memCounters.PagefileUsage)
		info.MemorySummary.TotalWorkingSet = uint64(memCounters.WorkingSetSize)
		info.MemorySummary.PrivateBytes = uint64(memCounters.PrivateUsage)
	}
}

func collectThreadInfo(info *CrashInfo) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return
	}
	defer windows.CloseHandle(snap)

	pid := uint32(os.Getpid())
	var te windows.ThreadEntry32
	te.Size = uint32(unsafe.Sizeof(te))

	if err := windows.Thread32First(snap, &te); err != nil {
		return
	}

	ownThreadCount := 0
	for {
		if te.OwnerProcessID == pid {
			ownThreadCount++
		}
		te.Size = uint32(unsafe.Sizeof(te))
		if err := windows.Thread32Next(snap, &te); err != nil {
			break
		}
	}
	info.ProcessInfo.Threads = ownThreadCount
}

func loadModules(info *CrashInfo) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE|windows.TH32CS_SNAPMODULE32, uint32(os.Getpid()))
	if err != nil {
		return
	}
	defer windows.CloseHandle(snap)

	var me windows.ModuleEntry32
	me.Size = uint32(unsafe.Sizeof(me))

	if err := windows.Module32First(snap, &me); err != nil {
		return
	}

	for {
		modPath := windows.UTF16ToString(me.ExePath[:])
		modName := windows.UTF16ToString(me.Module[:])
		info.LoadedModules = append(info.LoadedModules, ModuleInfo{
			Name:   modName,
			Path:   modPath,
			Base:   uint64(me.ModBaseAddr),
			Size:   uint64(me.ModBaseSize),
			Loaded: true,
		})
		me.Size = uint32(unsafe.Sizeof(me))
		if err := windows.Module32Next(snap, &me); err != nil {
			break
		}
	}
}

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

func writeCoreDump(path string, fullMemory bool) (string, error) {
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

	flags := uintptr(defaultDumpFlags)
	if fullMemory {
		flags |= miniDumpWithFullMemory
	}

	ret, _, callErr := syscall.SyscallN(
		proc,
		uintptr(windows.CurrentProcess()),
		uintptr(windows.GetCurrentProcessId()),
		uintptr(hFile),
		flags,
		0, 0, 0,
	)
	if ret == 0 {
		return "", fmt.Errorf("MiniDumpWriteDump failed: %w", callErr)
	}
	return abs, nil
}
