package system

import (
	"fmt"
	"github.com/jaypipes/ghw"
	"github.com/wailsapp/wails/v2/internal/shell"
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
	"github.com/wailsapp/wails/v2/pkg/options"
	"runtime"
	"strings"
)

var GlobalOptions options.App

type Report struct {
	OS             string                        `json:"os"`
	Arch           string                        `json:"arch"`
	CPU            string                        `json:"cpu"`
	GPU            string                        `json:"gpu"`
	Memory         string                        `json:"memory"`
	OSInfo         *operatingsystem.OS           `json:"osInfo"`
	PackageManager string                        `json:"packageManager"`
	Dependencies   packagemanager.DependencyList `json:"dependencies"`
	Options        options.App                   `json:"options"`
}

func GenerateReport() (*Report, error) {
	info := &Report{}

	info.OS = runtime.GOOS
	info.Arch = runtime.GOARCH

	cpu, _ := ghw.CPU()
	if cpu != nil {
		info.CPU = cpu.Processors[0].Model
	}

	info.GPU = getGPUInfo()

	memory, _ := ghw.Memory()
	if memory != nil {
		info.Memory = fmt.Sprintf("%dGB", memory.TotalPhysicalBytes/1024/1024/1024)
	} else {
		info.Memory = "Unknown"
	}

	sysInfo, err := GetInfo()
	if err != nil {
		return nil, err
	}

	info.OSInfo = sysInfo.OS
	if sysInfo.PM == nil {
		info.PackageManager = "None"
	} else {
		info.PackageManager = sysInfo.PM.Name()
	}
	info.Dependencies = sysInfo.Dependencies

	// Add additional dependencies
	info.Dependencies = append(info.Dependencies,
		checkNodejs(),
		checkNPM(),
		checkUPX(),
		checkNSIS(),
		checkDocker(),
		checkLibrary("gtk+-3.0")(),
		checkLibrary("webkit2gtk-4.0")(),
	)

	info.Options = GlobalOptions
	info.Options.AssetServer.Middleware = nil

	return info, nil
}

func getGPUInfo() string {
	gpu, _ := ghw.GPU(ghw.WithDisableWarnings())
	if gpu != nil {
		var gpuDetails []string
		for idx, card := range gpu.GraphicsCards {
			prefix := "GPU"
			if len(gpu.GraphicsCards) > 1 {
				prefix = fmt.Sprintf("GPU %d", idx+1)
			}
			if card.DeviceInfo == nil {
				gpuDetails = append(gpuDetails, fmt.Sprintf("%s: Unknown", prefix))
				continue
			}
			details := fmt.Sprintf("%s: %s (%s) - Driver: %s", prefix, card.DeviceInfo.Product.Name, card.DeviceInfo.Vendor.Name, card.DeviceInfo.Driver)
			gpuDetails = append(gpuDetails, details)
		}
		return strings.Join(gpuDetails, "; ")
	}

	gpuInfo := "Unknown"
	if runtime.GOOS == "darwin" {
		if stdout, _, err := shell.RunCommand("", "system_profiler", "SPDisplaysDataType"); err == nil {
			var gpuInfoDetails []string
			startCapturing := false
			for _, line := range strings.Split(stdout, "\n") {
				if strings.Contains(line, "Chipset Model") {
					startCapturing = true
				}
				if startCapturing {
					gpuInfoDetails = append(gpuInfoDetails, strings.TrimSpace(line))
				}
				if strings.Contains(line, "Metal Support") {
					break
				}
			}
			if len(gpuInfoDetails) > 0 {
				gpuInfo = strings.Join(gpuInfoDetails, " ")
			}
		}
	}
	return gpuInfo
}
