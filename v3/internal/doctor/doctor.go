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
	"github.com/pterm/pterm"
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

	spinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone().Start("Scanning system - Please wait (this may take a long time)...")

	defer func() {
		if err != nil {
			spinner.Fail()
		}
	}()

	/** Build **/

	// BuildSettings contains the build settings for the application
	var BuildSettings map[string]string

	// BuildInfo contains the build info for the application
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

	// Get system info
	info, err := operatingsystem.Info()
	if err != nil {
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
		// Get the latest commit hash
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

	spinner.Success()

	/** Output **/

	term.Section("System")

	systemTabledata := pterm.TableData{
		{pterm.Sprint("Name"), info.Name},
		{pterm.Sprint("Version"), info.Version},
		{pterm.Sprint("ID"), info.ID},
		{pterm.Sprint("Branding"), info.Branding},

		{pterm.Sprint("Platform"), runtime.GOOS},
		{pterm.Sprint("Architecture"), runtime.GOARCH},
	}

	mapKeys := lo.Keys(platformExtras)
	slices.Sort(mapKeys)
	for _, key := range mapKeys {
		systemTabledata = append(systemTabledata, []string{key, platformExtras[key]})
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
		systemTabledata = append(systemTabledata, []string{"CPU", "Unknown"})
	}

	// Probe GPU
	gpu, _ := ghw.GPU(ghw.WithDisableWarnings())
	if gpu != nil {
		for idx, card := range gpu.GraphicsCards {
			details := "Unknown"
			prefix := "GPU " + strconv.Itoa(idx+1)
			if card.DeviceInfo != nil {
				details = fmt.Sprintf("%s (%s) - Driver: %s ", card.DeviceInfo.Product.Name, card.DeviceInfo.Vendor.Name, card.DeviceInfo.Driver)
			}
			systemTabledata = append(systemTabledata, []string{prefix, details})
		}
	} else {
		if runtime.GOOS == "darwin" {
			var numCoresValue string
			cmd := exec.Command("sh", "-c", "ioreg -l | grep gpu-core-count")
			output, err := cmd.Output()
			if err == nil {
				// Look for an `=` sign, optional spaces and then an integer
				re := regexp.MustCompile(`= *(\d+)`)
				matches := re.FindAllStringSubmatch(string(output), -1)
				numCoresValue = "Unknown"
				if len(matches) > 0 {
					numCoresValue = matches[0][1]
				}

			}

			// Run `system_profiler SPDisplaysDataType | grep Metal`
			var metalSupport string
			cmd = exec.Command("sh", "-c", "system_profiler SPDisplaysDataType | grep Metal")
			output, err = cmd.Output()
			if err == nil {
				metalSupport = ", " + strings.TrimSpace(string(output))
			}
			systemTabledata = append(systemTabledata, []string{"GPU", numCoresValue + " cores" + metalSupport})

		} else {
			systemTabledata = append(systemTabledata, []string{"GPU", "Unknown"})
		}
	}

	memory, _ := ghw.Memory()
	var memoryText = "Unknown"
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
	systemTabledata = append(systemTabledata, []string{"Memory", memoryText})

	err = pterm.DefaultTable.WithBoxed().WithData(systemTabledata).Render()
	if err != nil {
		return err
	}

	// Build Environment

	term.Section("Build Environment")

	tableData := pterm.TableData{
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
			tableData = append(tableData, []string{name, buildSetting.Value})
		}
	}

	mapKeys = lo.Keys(BuildSettings)
	slices.Sort(mapKeys)
	for _, key := range mapKeys {
		tableData = append(tableData, []string{key, BuildSettings[key]})
	}

	err = pterm.DefaultTable.WithBoxed(true).WithData(tableData).Render()
	if err != nil {
		return err
	}

	// Dependencies
	term.Section("Dependencies")
	dependenciesBox := pterm.DefaultBox.WithTitleBottomCenter().WithTitle(pterm.Gray("*") + " - Optional Dependency")
	dependencyTableData := pterm.TableData{}
	if len(dependencies) == 0 {
		pterm.Info.Println("No dependencies found")
	} else {
		var optionals pterm.TableData
		mapKeys = lo.Keys(dependencies)
		for _, key := range mapKeys {
			if strings.HasPrefix(dependencies[key], "*") {
				optionals = append(optionals, []string{key, dependencies[key]})
			} else {
				dependencyTableData = append(dependencyTableData, []string{key, dependencies[key]})
			}
		}
		dependencyTableData = append(dependencyTableData, optionals...)
		dependenciesTableString, _ := pterm.DefaultTable.WithData(dependencyTableData).Srender()
		dependenciesBox.Println(dependenciesTableString)
	}

	// Run diagnostics after system info
	term.Section("Checking for issues")

	diagnosticResults := RunDiagnostics()
	if len(diagnosticResults) == 0 {
		pterm.Success.Println("No issues found")
	} else {
		pterm.Warning.Println("Found potential issues:")
		for _, result := range diagnosticResults {
			pterm.Printf("• %s: %s\n", result.TestName, result.ErrorMsg)
			url := result.HelpURL
			if strings.HasPrefix(url, "/") {
				url = "https://v3.wails.io" + url
			}
			pterm.Printf("  For more information: %s\n", term.Hyperlink(url, url))
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
