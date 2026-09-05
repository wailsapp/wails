package doctorng

import (
	"bytes"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw"
)

func collectHardwareInfo() HardwareInfo {
	hw := HardwareInfo{
		CPUs:   make([]CPUInfo, 0),
		GPUs:   make([]GPUInfo, 0),
		Memory: "Unknown",
	}

	hw.CPUs = collectCPUs()
	hw.GPUs = collectGPUs()
	hw.Memory = collectMemory()

	return hw
}

func collectCPUs() []CPUInfo {
	cpus, err := ghw.CPU()
	if err != nil || cpus == nil {
		return []CPUInfo{{Model: "Unknown"}}
	}

	result := make([]CPUInfo, 0, len(cpus.Processors))
	for _, cpu := range cpus.Processors {
		result = append(result, CPUInfo{
			Model: cpu.Model,
			Cores: int(cpu.NumCores),
		})
	}
	return result
}

func collectGPUs() []GPUInfo {
	gpu, err := ghw.GPU(ghw.WithDisableWarnings())
	if err == nil && gpu != nil {
		result := make([]GPUInfo, 0, len(gpu.GraphicsCards))
		for _, card := range gpu.GraphicsCards {
			info := GPUInfo{Name: "Unknown"}
			if card.DeviceInfo != nil {
				if card.DeviceInfo.Product != nil {
					info.Name = card.DeviceInfo.Product.Name
				}
				if card.DeviceInfo.Vendor != nil {
					info.Vendor = card.DeviceInfo.Vendor.Name
				}
				info.Driver = card.DeviceInfo.Driver
			}
			result = append(result, info)
		}
		if len(result) > 0 {
			return result
		}
	}

	if runtime.GOOS == "darwin" {
		return collectMacGPU()
	}

	return []GPUInfo{{Name: "Unknown"}}
}

func collectMacGPU() []GPUInfo {
	var numCores string
	cmd := exec.Command("sh", "-c", "ioreg -l | grep gpu-core-count")
	output, err := cmd.Output()
	if err == nil {
		re := regexp.MustCompile(`= *(\d+)`)
		matches := re.FindAllStringSubmatch(string(output), -1)
		if len(matches) > 0 {
			numCores = matches[0][1]
		}
	}

	var metalSupport string
	cmd = exec.Command("sh", "-c", "system_profiler SPDisplaysDataType | grep Metal")
	output, err = cmd.Output()
	if err == nil {
		metalSupport = strings.TrimSpace(string(output))
	}

	name := "Apple GPU"
	if numCores != "" {
		name = numCores + " cores"
	}
	if metalSupport != "" {
		name += ", " + metalSupport
	}

	return []GPUInfo{{Name: name}}
}

func collectMemory() string {
	memory, err := ghw.Memory()
	if err == nil && memory != nil {
		return strconv.Itoa(int(memory.TotalPhysicalBytes/1024/1024/1024)) + "GB"
	}

	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", "system_profiler SPHardwareDataType | grep 'Memory'")
		output, err := cmd.Output()
		if err == nil {
			output = bytes.Replace(output, []byte("Memory: "), []byte(""), 1)
			return strings.TrimSpace(string(output))
		}
	}

	return "Unknown"
}
