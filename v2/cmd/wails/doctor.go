package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/internal/shell"
	"github.com/wailsapp/wails/v2/internal/system"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
	"github.com/wailsapp/wails/v2/internal/tui"
)

func diagnoseEnvironment(f *flags.Doctor) error {
	if f.NoColour {
		tui.SetNoColour()
	}

	fmt.Println()
	fmt.Println(tui.Bold("Wails Doctor"))
	fmt.Println()

	var info *system.Info
	err := tui.WithSpinner("Scanning system", func() error {
		var scanErr error
		info, scanErr = system.GetInfo()
		return scanErr
	})
	if err != nil {
		tui.Error("Failed to get system information")
		return err
	}

	tui.Section("Wails")

	wailsTableData := [][]string{
		{"Version", app.Version()},
	}

	if buildInfo, _ := debug.ReadBuildInfo(); buildInfo != nil {
		buildSettingToName := map[string]string{
			"vcs.revision": "Revision",
			"vcs.modified": "Modified",
		}
		for _, buildSetting := range buildInfo.Settings {
			name := buildSettingToName[buildSetting.Key]
			if name == "" {
				continue
			}
			wailsTableData = append(wailsTableData, []string{name, buildSetting.Value})
		}
	}

	if info.PM != nil {
		wailsTableData = append(wailsTableData, []string{"Package Manager", info.PM.Name()})
	}

	tui.Table(wailsTableData)

	tui.Section("System")

	systemTabledata := [][]string{
		{"OS", info.OS.Name},
		{"Version", info.OS.Version},
		{"ID", info.OS.ID},
		{"Branding", info.OS.Branding},
		{"Go Version", runtime.Version()},
		{"Platform", runtime.GOOS},
		{"Architecture", runtime.GOARCH},
	}

	// Probe CPU
	cpus, _ := ghw.CPU()
	if cpus != nil {
		prefix := "CPU"
		for idx, cpu := range cpus.Processors {
			if len(cpus.Processors) > 1 {
				prefix = "CPU " + strconv.Itoa(idx+1)
			}
			systemTabledata = append(systemTabledata, []string{prefix, cpu.Model})
		}
	} else {
		cpuInfo := "Unknown"
		if runtime.GOOS == "darwin" {
			if stdout, _, err := shell.RunCommand("", "sysctl", "-n", "machdep.cpu.brand_string"); err == nil {
				cpuInfo = strings.TrimSpace(stdout)
			}
		}
		systemTabledata = append(systemTabledata, []string{"CPU", cpuInfo})
	}

	// Probe GPU
	gpu, _ := ghw.GPU(ghw.WithDisableWarnings())
	if gpu != nil {
		prefix := "GPU"
		for idx, card := range gpu.GraphicsCards {
			if len(gpu.GraphicsCards) > 1 {
				prefix = "GPU " + strconv.Itoa(idx+1) + " "
			}
			if card.DeviceInfo == nil {
				systemTabledata = append(systemTabledata, []string{prefix, "Unknown"})
				continue
			}
			details := fmt.Sprintf("%s (%s) - Driver: %s", card.DeviceInfo.Product.Name, card.DeviceInfo.Vendor.Name, card.DeviceInfo.Driver)
			systemTabledata = append(systemTabledata, []string{prefix, details})
		}
	} else {
		gpuInfo := "Unknown"
		if runtime.GOOS == "darwin" {
			if stdout, _, err := shell.RunCommand("", "system_profiler", "SPDisplaysDataType"); err == nil {
				var (
					startCapturing bool
					gpuInfoDetails []string
				)
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
		systemTabledata = append(systemTabledata, []string{"GPU", gpuInfo})
	}

	memory, _ := ghw.Memory()
	if memory != nil {
		systemTabledata = append(systemTabledata, []string{"Memory", strconv.Itoa(int(memory.TotalPhysicalBytes/1024/1024/1024)) + "GB"})
	} else {
		memInfo := "Unknown"
		if runtime.GOOS == "darwin" {
			if stdout, _, err := shell.RunCommand("", "sysctl", "-n", "hw.memsize"); err == nil {
				if memSize, err := strconv.Atoi(strings.TrimSpace(stdout)); err == nil {
					memInfo = strconv.Itoa(memSize/1024/1024/1024) + "GB"
				}
			}
		}
		systemTabledata = append(systemTabledata, []string{"Memory", memInfo})
	}

	tui.Table(systemTabledata)

	tui.Section("Dependencies")

	var dependenciesMissing []string
	var externalPackages []*packagemanager.Dependency
	dependenciesAvailableRequired := 0
	dependenciesAvailableOptional := 0

	dependenciesTableData := [][]string{
		{"Dependency", "Package Name", "Status", "Version"},
	}

	hasOptionalDependencies := false
	for _, dependency := range info.Dependencies {
		name := dependency.Name

		if dependency.Optional {
			name = tui.Gray("*") + name
			hasOptionalDependencies = true
		}

		packageName := "Unknown"
		status := tui.Red("Not Found")

		if dependency.PackageName != "" {
			packageName = dependency.PackageName

			if dependency.Installed {
				status = tui.Green("Installed")
			} else {
				status = tui.Magenta("Available")

				if dependency.Optional {
					dependenciesAvailableOptional++
				} else {
					dependenciesAvailableRequired++
				}
			}
		} else {
			if !dependency.Optional {
				dependenciesMissing = append(dependenciesMissing, dependency.Name)
			}

			if dependency.External {
				externalPackages = append(externalPackages, dependency)
			}
		}

		dependenciesTableData = append(dependenciesTableData, []string{name, packageName, status, dependency.Version})
	}

	tui.HeaderTable(dependenciesTableData)

	if hasOptionalDependencies {
		fmt.Println("\n  " + tui.Gray("*") + " = Optional Dependency")
	}
	_ = externalPackages

	tui.Section("Diagnosis")

	if dependenciesAvailableRequired != 0 {
		fmt.Println("Required package(s) installation details:\n" + info.Dependencies.InstallAllRequiredCommand())
	}

	if dependenciesAvailableOptional != 0 {
		fmt.Println("Optional package(s) installation details:\n" + info.Dependencies.InstallAllOptionalCommand())
	}

	if len(dependenciesMissing) == 0 && dependenciesAvailableRequired == 0 {
		tui.Success("Your system is ready for Wails development!")
	} else {
		tui.Warning("Your system has missing dependencies!")
	}

	if len(dependenciesMissing) != 0 {
		fmt.Println("Fatal:")
		fmt.Println("Required dependencies missing: " + strings.Join(dependenciesMissing, " "))
	}

	fmt.Println()
	return nil
}
