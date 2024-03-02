package doctor

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/jaypipes/ghw"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/internal/version"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"slices"
	"strconv"
)

func getTaskVersion() (string, bool) {

	var taskVersion string = ""

	// Execute task
	output, err := exec.Command("task", "--version").Output()
	if err != nil {
		taskVersion = "Not Installed"
		return taskVersion, false
	}

	// Extract the version using regular expression
	match := regexp.MustCompile(`v(\d+\.\d+\.\d+)`).FindStringSubmatch(string(output))

	// Check if the version is found
	if match != nil {
		taskVersion = match[1]
	} else {
		taskVersion = "Unknown"
	}

	return taskVersion, true
}

func Run() (err error) {

	pterm.DefaultSection = *pterm.DefaultSection.
		WithBottomPadding(0).
		WithStyle(pterm.NewStyle(pterm.FgBlue, pterm.Bold))

	pterm.Println() // Spacer
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println("Wails Doctor")
	pterm.Println() // Spacer

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
		pterm.Error.Println("Failed to get system information")
		return err
	}

	/** Wails **/
	wailsPackage, _ := lo.Find(BuildInfo.Deps, func(dep *debug.Module) bool {
		return dep.Path == "github.com/wailsapp/wails/v3"
	})

	wailsVersion := version.VersionString
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

	spinner.Success()

	/** Output **/

	pterm.DefaultSection.Println("Build Environment")

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

	mapKeys := lo.Keys(BuildSettings)
	slices.Sort(mapKeys)
	for _, key := range mapKeys {
		tableData = append(tableData, []string{key, BuildSettings[key]})
	}

	//// Exit early if PM not found
	//if info.PM != nil {
	//	wailsTableData = append(wailsTableData, []string{"Package Manager", info.PM.Name()})
	//}

	err = pterm.DefaultTable.WithData(tableData).Render()
	if err != nil {
		return err
	}

	pterm.DefaultSection.Println("System")

	systemTabledata := pterm.TableData{
		{pterm.Sprint("Name"), info.Name},
		{pterm.Sprint("Version"), info.Version},
		{pterm.Sprint("ID"), info.ID},
		{pterm.Sprint("Branding"), info.Branding},

		{pterm.Sprint("Platform"), runtime.GOOS},
		{pterm.Sprint("Architecture"), runtime.GOARCH},
	}

	/* Task (Taskfile) */
	taskVersion, ok := getTaskVersion()
	systemTabledata = append(systemTabledata, []string{"Task", taskVersion})

	mapKeys = lo.Keys(platformExtras)
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
		systemTabledata = append(systemTabledata, []string{"GPU", "Unknown"})
	}

	memory, _ := ghw.Memory()
	if memory != nil {
		systemTabledata = append(systemTabledata, []string{"Memory", strconv.Itoa(int(memory.TotalPhysicalBytes/1024/1024/1024)) + "GB"})
	} else {
		systemTabledata = append(systemTabledata, []string{"Memory", "Unknown"})
	}

	//systemTabledata = append(systemTabledata, []string{"CPU", cpu.Processors[0].Model})

	err = pterm.DefaultTable.WithData(systemTabledata).Render()
	if err != nil {
		return err
	}

	pterm.DefaultSection.Println("Diagnosis")
	if !ok {
		pterm.Warning.Println("There are some items above that need addressing!")
	} else {
		pterm.Success.Println("Your system is ready for Wails development!")
	}

	return nil
}
