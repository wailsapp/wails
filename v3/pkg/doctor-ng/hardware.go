package doctorng

import (
	"bytes"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/internal/hostinfo"
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
	cpus, err := hostinfo.Processors()
	if err != nil || cpus == nil {
		return []CPUInfo{{Model: "Unknown"}}
	}

	result := make([]CPUInfo, 0, len(cpus))
	for _, cpu := range cpus {
		result = append(result, CPUInfo{
			Model: cpu.Model,
			Cores: cpu.Cores,
		})
	}
	return result
}

func collectGPUs() []GPUInfo {
	cards, err := hostinfo.Graphics()
	if err == nil && len(cards) > 0 {
		result := make([]GPUInfo, 0, len(cards))
		for _, card := range cards {
			result = append(result, GPUInfo{Name: card.Name, Vendor: card.Vendor, Driver: card.Driver})
		}
		return result
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
	memory, err := hostinfo.Memory()
	if err == nil {
		return strconv.Itoa(int(memory/1024/1024/1024)) + "GB"
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
