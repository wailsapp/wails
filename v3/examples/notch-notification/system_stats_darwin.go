//go:build darwin

package main

import (
	"context"
	"encoding/json"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/unix"
)

const (
	bytesPerGiB           = 1024 * 1024 * 1024
	machineProfileTimeout = 2 * time.Second
)

// SystemStats is the current machine telemetry returned to the webview.
type SystemStats struct {
	Available     bool    `json:"available"`
	CPUPercent    float64 `json:"cpuPercent"`
	CoreCount     int     `json:"coreCount"`
	MachineName   string  `json:"machineName"`
	MemoryUsedGB  float64 `json:"memoryUsedGB"`
	MemoryFreeGB  float64 `json:"memoryFreeGB"`
	MemoryPercent float64 `json:"memoryPercent"`
	DiskUsedGB    float64 `json:"diskUsedGB"`
	DiskFreeGB    float64 `json:"diskFreeGB"`
	DiskPercent   float64 `json:"diskPercent"`
}

func (controller *NotificationController) systemStats() SystemStats {
	sample, ok := readSystemSample()
	if !ok {
		return SystemStats{}
	}

	active := sample.activeTicks
	total := sample.totalTicks
	memoryUsed := sample.memoryUsed
	memoryTotal := sample.memoryTotal
	diskUsed := sample.diskUsed
	diskTotal := sample.diskTotal

	controller.mu.Lock()
	cpuPercent := controller.lastCPUPercent
	if controller.lastCPUTotal > 0 && total > controller.lastCPUTotal {
		activeDelta := active - controller.lastCPUActive
		totalDelta := total - controller.lastCPUTotal
		cpuPercent = float64(activeDelta) / float64(totalDelta) * 100
	}
	controller.lastCPUActive = active
	controller.lastCPUTotal = total
	controller.lastCPUPercent = cpuPercent
	controller.mu.Unlock()

	memoryFree := memoryTotal - memoryUsed
	diskFree := diskTotal - diskUsed
	return SystemStats{
		Available:     true,
		CPUPercent:    cpuPercent,
		CoreCount:     runtime.NumCPU(),
		MachineName:   machineName(),
		MemoryUsedGB:  float64(memoryUsed) / bytesPerGiB,
		MemoryFreeGB:  float64(memoryFree) / bytesPerGiB,
		MemoryPercent: percentage(memoryUsed, memoryTotal),
		DiskUsedGB:    float64(diskUsed) / bytesPerGiB,
		DiskFreeGB:    float64(diskFree) / bytesPerGiB,
		DiskPercent:   percentage(diskUsed, diskTotal),
	}
}

// systemSample is one reading of the machine counters. The mach calls go
// through purego so the example builds with and without cgo.
type systemSample struct {
	activeTicks uint64
	totalTicks  uint64
	memoryUsed  uint64
	memoryTotal uint64
	diskUsed    uint64
	diskTotal   uint64
}

const (
	processorCPULoadInfo = 2  // PROCESSOR_CPU_LOAD_INFO
	hostVMInfo64         = 4  // HOST_VM_INFO64
	hostVMInfo64Count    = 38 // HOST_VM_INFO64_COUNT: sizeof(vm_statistics64) / sizeof(integer_t)
	cpuStateUser         = 0
	cpuStateSystem       = 1
	cpuStateIdle         = 2
	cpuStateNice         = 3
	cpuStateMax          = 4
	kernSuccess          = 0
)

// vmStatistics64 mirrors vm_statistics64 from <mach/vm_statistics.h>.
type vmStatistics64 struct {
	freeCount, activeCount, inactiveCount, wireCount        uint32
	zeroFillCount, reactivations, pageins, pageouts, faults uint64
	cowFaults, lookups, hits, purges                        uint64
	purgeableCount, speculativeCount                        uint32
	decompressions, compressions, swapins, swapouts         uint64
	compressorPageCount, throttledCount                     uint32
	externalPageCount, internalPageCount                    uint32
	totalUncompressedPagesInCompressor                      uint64
}

var (
	machOnce          sync.Once
	machErr           error
	machHostSelf      func() uint32
	machTaskSelf      uint32
	hostProcessorInfo func(host uint32, flavor int32, count *uint32, info *unsafe.Pointer, infoCount *uint32) int32
	hostStatistics64  func(host uint32, flavor int32, info unsafe.Pointer, count *uint32) int32
	hostPageSize      func(host uint32, size *uintptr) int32
	vmDeallocate      func(task uint32, address unsafe.Pointer, size uintptr) int32
)

func loadMach() error {
	machOnce.Do(func() {
		lib, err := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			machErr = err
			return
		}
		task, err := purego.Dlsym(lib, "mach_task_self_")
		if err != nil {
			machErr = err
			return
		}
		// task is the address of the mach_task_self_ variable.
		machTaskSelf = **(**uint32)(unsafe.Pointer(&task))
		purego.RegisterLibFunc(&machHostSelf, lib, "mach_host_self")
		purego.RegisterLibFunc(&hostProcessorInfo, lib, "host_processor_info")
		purego.RegisterLibFunc(&hostStatistics64, lib, "host_statistics64")
		purego.RegisterLibFunc(&hostPageSize, lib, "host_page_size")
		purego.RegisterLibFunc(&vmDeallocate, lib, "vm_deallocate")
	})
	return machErr
}

func readSystemSample() (systemSample, bool) {
	var sample systemSample
	if loadMach() != nil {
		return sample, false
	}
	host := machHostSelf()

	var cpuCount, infoCount uint32
	var info unsafe.Pointer
	if hostProcessorInfo(host, processorCPULoadInfo, &cpuCount, &info, &infoCount) != kernSuccess {
		return sample, false
	}
	ticks := unsafe.Slice((*uint32)(info), int(cpuCount)*cpuStateMax)
	for cpu := 0; cpu < int(cpuCount); cpu++ {
		load := ticks[cpu*cpuStateMax : (cpu+1)*cpuStateMax]
		busy := uint64(load[cpuStateUser]) + uint64(load[cpuStateSystem]) + uint64(load[cpuStateNice])
		sample.activeTicks += busy
		sample.totalTicks += busy + uint64(load[cpuStateIdle])
	}
	vmDeallocate(machTaskSelf, info, uintptr(infoCount)*4)

	memoryTotal, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return sample, false
	}

	var vmStats vmStatistics64
	vmStatsCount := uint32(hostVMInfo64Count)
	if hostStatistics64(host, hostVMInfo64, unsafe.Pointer(&vmStats), &vmStatsCount) != kernSuccess {
		return sample, false
	}
	var pageSize uintptr
	if hostPageSize(host, &pageSize) != kernSuccess {
		return sample, false
	}
	reclaimable := (uint64(vmStats.freeCount) + uint64(vmStats.inactiveCount) + uint64(vmStats.purgeableCount)) * uint64(pageSize)
	sample.memoryTotal = memoryTotal
	if reclaimable < memoryTotal {
		sample.memoryUsed = memoryTotal - reclaimable
	}

	var disk unix.Statfs_t
	if unix.Statfs("/", &disk) != nil {
		return sample, false
	}
	blockSize := uint64(disk.Bsize)
	sample.diskTotal = disk.Blocks * blockSize
	if free := disk.Bavail * blockSize; free < sample.diskTotal {
		sample.diskUsed = sample.diskTotal - free
	}
	return sample, true
}

var (
	machineNameOnce   sync.Once
	cachedMachineName = "Apple silicon"
)

func machineName() string {
	machineNameOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), machineProfileTimeout)
		defer cancel()
		output, err := exec.CommandContext(
			ctx,
			"/usr/sbin/system_profiler",
			"SPHardwareDataType",
			"-json",
			"-detailLevel",
			"mini",
		).Output()
		if err != nil {
			return
		}

		var profile struct {
			Hardware []struct {
				ChipType    string `json:"chip_type"`
				MachineName string `json:"machine_name"`
			} `json:"SPHardwareDataType"`
		}
		if json.Unmarshal(output, &profile) != nil || len(profile.Hardware) == 0 {
			return
		}

		chip := strings.TrimSpace(strings.TrimPrefix(profile.Hardware[0].ChipType, "Apple "))
		model := strings.TrimSpace(profile.Hardware[0].MachineName)
		if chip != "" && model != "" {
			cachedMachineName = chip + " " + model
		}
	})
	return cachedMachineName
}

func percentage(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) / float64(total) * 100
}
