package debug

import (
	"time"

	doctor "github.com/wailsapp/wails/v3/pkg/doctor-ng"
)

type SystemInfo = doctor.SystemInfo
type BuildInfo = doctor.BuildInfo
type Dependency = doctor.Dependency
type DiagnosticResult = doctor.DiagnosticResult
type HardwareInfo = doctor.HardwareInfo
type OSInfo = doctor.OSInfo
type CPUInfo = doctor.CPUInfo
type GPUInfo = doctor.GPUInfo

type CrashInfo struct {
	Timestamp      time.Time         `json:"timestamp"`
	OS             OSInfo            `json:"os"`
	Hardware       HardwareInfo     `json:"hardware"`
	Build          BuildInfo         `json:"build"`
	PanicMessage   string           `json:"panic_message,omitempty"`
	PanicStack     []StackFrame     `json:"panic_stack,omitempty"`
	Environment   map[string]string `json:"environment"`
	OpenFiles      []OpenFileInfo   `json:"open_files,omitempty"`
	LoadedModules  []ModuleInfo    `json:"loaded_modules,omitempty"`
	ProcessInfo    ProcessInfo     `json:"process_info,omitempty"`
	ThreadInfo     []ThreadInfo     `json:"threads,omitempty"`
	MemorySummary  MemorySummary   `json:"memory_summary,omitempty"`
	NetworkConns   []ConnInfo       `json:"network_connections,omitempty"`
}

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	PC       uintptr `json:"pc"`
}

type OpenFileInfo struct {
	Path    string `json:"path"`
	Mode    string `json:"mode,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Device  string `json:"device,omitempty"`
}

type ModuleInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Base   uint64 `json:"base"`
	Size   uint64 `json:"size"`
	Loaded bool   `json:"loaded"`
}

type ProcessInfo struct {
	PID            int    `json:"pid"`
	PPID           int    `json:"ppid"`
	Threads       int    `json:"threads"`
	Handles       int    `json:"handles"`
	Goroutines    int    `json:"goroutines,omitempty"`
	MemoryBytes   uint64 `json:"memory_bytes"`
	Status        string `json:"status,omitempty"`
	StartTime     time.Time `json:"start_time,omitempty"`
	CPUPercent    float64 `json:"cpu_percent,omitempty"`
	UserTime       float64 `json:"user_time_ms,omitempty"`
	SystemTime    float64 `json:"system_time_ms,omitempty"`
}

type ThreadInfo struct {
	ID      int         `json:"id"`
	State   string      `json:"state"`
	WaitChan chan struct{}

	Stack  StackInfo `json:"stack,omitempty"`
	UserTime float64 `json:"user_time_ms,omitempty"`
	SystemTime float64 `json:"system_time_ms,omitempty"`
}

type StackInfo struct {
	Start uintptr   `json:"start"`
	End   uintptr   `json:"end"`
	Frames []StackFrame `json:"frames,omitempty"`
}

type MemorySummary struct {
	TotalVirtual     uint64 `json:"total_virtual_bytes"`
	TotalWorkingSet  uint64 `json:"total_working_set_bytes"`
	PrivateBytes     uint64 `json:"private_bytes"`
	SharedBytes      uint64 `json:"shared_bytes"`
	GarbageCollector GCStats `json:"gc_stats,omitempty"`
}

type GCStats struct {
	NumGC         int    `json:"num_gc"`
	LastGC        uint32 `json:"last_gc"`
	TotalPause   uint64 `json:"total_pause_ns"`
	TotalAlloc    uint64 `json:"total_alloc_bytes"`
	HeapAlloc    uint64 `json:"heap_alloc_bytes"`
	HeapObjects  uint64 `json:"heap_objects"`
}

type ConnInfo struct {
	Protocol  string `json:"protocol"`
	LocalAddr string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
	State     string `json:"state"`
	UID       uint32 `json:"uid,omitempty"`
}

type Report struct {
	Timestamp    time.Time     `json:"timestamp"`
	System       SystemInfo    `json:"system"`
	Build        BuildInfo      `json:"build"`
	CrashInfo    *CrashInfo     `json:"crash_info,omitempty"`
	Diagnostics  []DiagnosticResult `json:"diagnostics,omitempty"`
}

func NewReport() *Report {
	return &Report{
		Timestamp: time.Now(),
	}
}