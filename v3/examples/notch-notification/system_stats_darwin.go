//go:build darwin

package main

/*
#include <mach/mach.h>
#include <mach/mach_host.h>
#include <mach/processor_info.h>
#include <stdint.h>
#include <sys/mount.h>
#include <sys/sysctl.h>

typedef struct {
    uint64_t active_ticks;
    uint64_t total_ticks;
    uint64_t memory_used;
    uint64_t memory_total;
    uint64_t disk_used;
    uint64_t disk_total;
} WailsSystemSample;

static int wails_read_system_sample(WailsSystemSample* sample) {
    natural_t cpuCount = 0;
    processor_info_array_t cpuInfo = NULL;
    mach_msg_type_number_t cpuInfoCount = 0;
    kern_return_t result = host_processor_info(
        mach_host_self(),
        PROCESSOR_CPU_LOAD_INFO,
        &cpuCount,
        &cpuInfo,
        &cpuInfoCount);
    if (result != KERN_SUCCESS) {
        return 0;
    }

    processor_cpu_load_info_t loads = (processor_cpu_load_info_t)cpuInfo;
    for (natural_t index = 0; index < cpuCount; index++) {
        uint64_t user = loads[index].cpu_ticks[CPU_STATE_USER];
        uint64_t system = loads[index].cpu_ticks[CPU_STATE_SYSTEM];
        uint64_t nice = loads[index].cpu_ticks[CPU_STATE_NICE];
        uint64_t idle = loads[index].cpu_ticks[CPU_STATE_IDLE];
        sample->active_ticks += user + system + nice;
        sample->total_ticks += user + system + nice + idle;
    }
    vm_deallocate(
        mach_task_self(),
        (vm_address_t)cpuInfo,
        (vm_size_t)cpuInfoCount * sizeof(integer_t));

    uint64_t memoryTotal = 0;
    size_t memoryTotalSize = sizeof(memoryTotal);
    if (sysctlbyname("hw.memsize", &memoryTotal, &memoryTotalSize, NULL, 0) != 0) {
        return 0;
    }

    vm_statistics64_data_t vmStats;
    mach_msg_type_number_t vmStatsCount = HOST_VM_INFO64_COUNT;
    result = host_statistics64(
        mach_host_self(),
        HOST_VM_INFO64,
        (host_info64_t)&vmStats,
        &vmStatsCount);
    if (result != KERN_SUCCESS) {
        return 0;
    }

    vm_size_t pageSize = 0;
    if (host_page_size(mach_host_self(), &pageSize) != KERN_SUCCESS) {
        return 0;
    }
    uint64_t reclaimablePages = (uint64_t)vmStats.free_count +
        (uint64_t)vmStats.inactive_count +
        (uint64_t)vmStats.purgeable_count;
    uint64_t reclaimableBytes = reclaimablePages * (uint64_t)pageSize;
    sample->memory_total = memoryTotal;
    sample->memory_used = reclaimableBytes < memoryTotal ? memoryTotal - reclaimableBytes : 0;

    struct statfs diskStats;
    if (statfs("/", &diskStats) != 0) {
        return 0;
    }
    uint64_t blockSize = (uint64_t)diskStats.f_bsize;
    sample->disk_total = (uint64_t)diskStats.f_blocks * blockSize;
    uint64_t diskFree = (uint64_t)diskStats.f_bavail * blockSize;
    sample->disk_used = diskFree < sample->disk_total ? sample->disk_total - diskFree : 0;
    return 1;
}
*/
import "C"

import (
	"encoding/json"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

const bytesPerGiB = 1024 * 1024 * 1024

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
	var sample C.WailsSystemSample
	if C.wails_read_system_sample(&sample) == 0 {
		return SystemStats{}
	}

	active := uint64(sample.active_ticks)
	total := uint64(sample.total_ticks)
	memoryUsed := uint64(sample.memory_used)
	memoryTotal := uint64(sample.memory_total)
	diskUsed := uint64(sample.disk_used)
	diskTotal := uint64(sample.disk_total)

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

var (
	machineNameOnce   sync.Once
	cachedMachineName = "Apple silicon"
)

func machineName() string {
	machineNameOnce.Do(func() {
		output, err := exec.Command(
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
