package doctor

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/internal/term"

	"github.com/wailsapp/wails/v3/internal/buildinfo"

	"github.com/go-git/go-git/v5"
	"github.com/jaypipes/ghw"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/internal/version"
)

func Run() (err error) {

	get, err := buildinfo.Get()
	if err != nil {
		return err
	}
	_ = get

	term.Header("Wails Doctor")

	spinner := term.StartSpinner("Scanning system - Please wait (this may take a long time)…")

	defer func() {
		if err != nil {
			term.StopSpinner(spinner)
		}
	}()

	/** Build **/

	var BuildSettings map[string]string
	var BuildInfo *debug.BuildInfo

	var ok bool
	BuildInfo, ok = debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("could not read build info from binary")
	}
	BuildSettings = lo.Associate(BuildInfo.Settings, func(setting debug.BuildSetting) (string, string) {
		return setting.Key, setting.Value
	})

	/** Operating System **/

	info, err := operatingsystem.Info()
	if err != nil {
		term.StopSpinner(spinner)
		term.Error("Failed to get system information")
		return err
	}

	/** Wails **/

	wailsPackage, _ := lo.Find(BuildInfo.Deps, func(dep *debug.Module) bool {
		return dep.Path == "github.com/wailsapp/wails/v3"
	})

	wailsVersion := strings.TrimSpace(version.String())
	if wailsPackage != nil && wailsPackage.Replace != nil {
		wailsVersion = "(local) => " + filepath.ToSlash(wailsPackage.Replace.Path)
		repo, err := git.PlainOpen(filepath.Join(wailsPackage.Replace.Path, ".."))
		if err == nil {
			head, err := repo.Head()
			if err == nil {
				wailsVersion += " (" + head.Hash().String()[:8] + ")"
			}
		}
	}

	platformExtras, ok := getInfo()

	dependencies := make(map[string]string)
	checkPlatformDependencies(dependencies, &ok)

	term.StopSpinner(spinner)

	/** Output **/

	term.Section("System")

	systemRows := [][]string{
		{"Name", info.Name},
		{"Version", info.Version},
		{"ID", info.ID},
		{"Branding", info.Branding},
		{"Platform", runtime.GOOS},
		{"Architecture", runtime.GOARCH},
	}

	mapKeys := lo.Keys(platformExtras)
	slices.Sort(mapKeys)
	for _, key := range mapKeys {
		systemRows = append(systemRows, []string{key, platformExtras[key]})
	}

	cpus, _ := ghw.CPU()
	if cpus != nil {
		prefix := "CPU"
		for idx, cpu := range cpus.Processors {
			if len(cpus.Processors) > 1 {
				prefix = "CPU " + strconv.Itoa(idx+1)
			}
			systemRows = append(systemRows, []string{prefix, cpu.Model})
		}
	} else {
		systemRows = append(systemRows, []string{"CPU", "Unknown"})
	}

	gpu, _ := ghw.GPU(ghw.WithDisableWarnings())
	if gpu != nil {
		for idx, card := range gpu.GraphicsCards {
			details := "Unknown"
			prefix := "GPU " + strconv.Itoa(idx+1)
			if card.DeviceInfo != nil {
				details = fmt.Sprintf("%s (%s) - Driver: %s", card.DeviceInfo.Product.Name, card.DeviceInfo.Vendor.Name, card.DeviceInfo.Driver)
			}
			systemRows = append(systemRows, []string{prefix, details})
		}
	} else {
		if runtime.GOOS == "darwin" {
			var numCoresValue string
			cmd := exec.Command("sh", "-c", "ioreg -l | grep gpu-core-count")
			output, err := cmd.Output()
			if err == nil {
				re := regexp.MustCompile(`= *(\d+)`)
				matches := re.FindAllStringSubmatch(string(output), -1)
				numCoresValue = "Unknown"
				if len(matches) > 0 {
					numCoresValue = matches[0][1]
				}
			}
			var metalSupport string
			cmd = exec.Command("sh", "-c", "system_profiler SPDisplaysDataType | grep Metal")
			output, err = cmd.Output()
			if err == nil {
				metalSupport = ", " + strings.TrimSpace(string(output))
			}
			systemRows = append(systemRows, []string{"GPU", numCoresValue + " cores" + metalSupport})
		} else {
			systemRows = append(systemRows, []string{"GPU", "Unknown"})
		}
	}

	memory, _ := ghw.Memory()
	memoryText := "Unknown"
	if memory != nil {
		memoryText = strconv.Itoa(int(memory.TotalPhysicalBytes/1024/1024/1024)) + "GB"
	} else {
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("sh", "-c", "system_profiler SPHardwareDataType | grep 'Memory'")
			output, err := cmd.Output()
			if err == nil {
				output = bytes.Replace(output, []byte("Memory: "), []byte(""), 1)
				memoryText = strings.TrimSpace(string(output))
			}
		}
	}
	systemRows = append(systemRows, []string{"Memory", memoryText})
	term.Table(systemRows)

	term.Section("Build Environment")

	buildRows := [][]string{
		{"Wails CLI", wailsVersion},
		{"Go Version", runtime.Version()},
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
			buildRows = append(buildRows, []string{name, buildSetting.Value})
		}
	}

	mapKeys = lo.Keys(BuildSettings)
	slices.Sort(mapKeys)
	for _, key := range mapKeys {
		buildRows = append(buildRows, []string{key, BuildSettings[key]})
	}
	term.Table(buildRows)

	term.Section("Dependencies")
	if len(dependencies) == 0 {
		term.Info("No dependencies found")
	} else {
		var depRows [][]string
		var optionalRows [][]string
		mapKeys = lo.Keys(dependencies)
		slices.Sort(mapKeys)
		for _, key := range mapKeys {
			if strings.HasPrefix(dependencies[key], "*") {
				optionalRows = append(optionalRows, []string{key, dependencies[key]})
			} else {
				depRows = append(depRows, []string{key, dependencies[key]})
			}
		}
		depRows = append(depRows, optionalRows...)
		term.Table(depRows)
		term.Println(term.Dim("* - Optional Dependency"))
	}

	term.Section("Checking for issues")

	diagnosticResults := RunDiagnostics()
	if len(diagnosticResults) == 0 {
		term.Success("No issues found")
	} else {
		term.Warning("Found potential issues:")
		for _, result := range diagnosticResults {
			term.Printf("  • %s: %s\n", result.TestName, result.ErrorMsg)
			url := result.HelpURL
			if strings.HasPrefix(url, "/") {
				url = "https://v3.wails.io" + url
			}
			term.Printf("    For more information: %s\n", term.Hyperlink(url, url))
		}
	}

	term.Section("Diagnosis")
	if !ok {
		term.Warning("There are some items above that need addressing!")
	} else {
		term.Success("Your system is ready for Wails development!")
	}

	return nil
}
