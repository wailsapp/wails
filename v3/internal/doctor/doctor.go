package doctor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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

type DoctorReport struct {
	System       SystemReport       `json:"system"`
	Build        BuildReport        `json:"build"`
	Dependencies map[string]string  `json:"dependencies"`
	Signing      SigningStatus      `json:"signing"`
	Diagnostics  []DiagnosticResult `json:"diagnostics"`
	Ready        bool               `json:"ready"`
}

type SystemReport struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	ID           string            `json:"id"`
	Branding     string            `json:"branding,omitempty"`
	Platform     string            `json:"platform"`
	Architecture string            `json:"architecture"`
	CPU          string            `json:"cpu,omitempty"`
	GPU          string            `json:"gpu,omitempty"`
	Memory       string            `json:"memory,omitempty"`
	Extras       map[string]string `json:"extras,omitempty"`
}

type BuildReport struct {
	WailsVersion string            `json:"wailsVersion"`
	GoVersion    string            `json:"goVersion"`
	Settings     map[string]string `json:"settings,omitempty"`
}

func Run(jsonOutput bool) (err error) {
	report, err := collectReport(jsonOutput)
	if err != nil {
		return err
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}

	return renderReport(report)
}

func collectReport(quiet bool) (*DoctorReport, error) {
	report := &DoctorReport{
		Dependencies: make(map[string]string),
		Ready:        true,
	}

	_, err := buildinfo.Get()
	if err != nil {
		return nil, err
	}

	var spinner *pterm.SpinnerPrinter
	if !quiet {
		term.Header("Wails Doctor")
		spinner, _ = pterm.DefaultSpinner.WithRemoveWhenDone().Start("Scanning system - Please wait (this may take a long time)...")
	}

	BuildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		if spinner != nil {
			spinner.Fail()
		}
		return nil, fmt.Errorf("could not read build info from binary")
	}
	BuildSettings := lo.Associate(BuildInfo.Settings, func(setting debug.BuildSetting) (string, string) {
		return setting.Key, setting.Value
	})

	info, err := operatingsystem.Info()
	if err != nil {
		if spinner != nil {
			spinner.Fail()
		}
		return nil, err
	}

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

	platformExtras, platformOK := getInfo()
	if !platformOK {
		report.Ready = false
	}

	checkPlatformDependencies(report.Dependencies, &report.Ready)

	report.System = SystemReport{
		Name:         info.Name,
		Version:      info.Version,
		ID:           info.ID,
		Branding:     info.Branding,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		Extras:       platformExtras,
	}

	cpus, _ := ghw.CPU()
	if cpus != nil && len(cpus.Processors) > 0 {
		report.System.CPU = cpus.Processors[0].Model
	}

	gpu, _ := ghw.GPU(ghw.WithDisableWarnings())
	if gpu != nil && len(gpu.GraphicsCards) > 0 {
		card := gpu.GraphicsCards[0]
		if card.DeviceInfo != nil {
			report.System.GPU = fmt.Sprintf("%s (%s)", card.DeviceInfo.Product.Name, card.DeviceInfo.Vendor.Name)
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", "ioreg -l | grep gpu-core-count")
		output, err := cmd.Output()
		if err == nil {
			re := regexp.MustCompile(`= *(\d+)`)
			matches := re.FindAllStringSubmatch(string(output), -1)
			if len(matches) > 0 {
				report.System.GPU = matches[0][1] + " cores"
			}
		}
	}

	memory, _ := ghw.Memory()
	if memory != nil {
		report.System.Memory = strconv.Itoa(int(memory.TotalPhysicalBytes/1024/1024/1024)) + "GB"
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", "system_profiler SPHardwareDataType | grep 'Memory'")
		output, err := cmd.Output()
		if err == nil {
			output = bytes.Replace(output, []byte("Memory: "), []byte(""), 1)
			report.System.Memory = strings.TrimSpace(string(output))
		}
	}

	report.Build = BuildReport{
		WailsVersion: wailsVersion,
		GoVersion:    runtime.Version(),
		Settings:     BuildSettings,
	}

	report.Signing = CheckSigning()
	report.Diagnostics = RunDiagnostics()

	if len(report.Diagnostics) > 0 {
		report.Ready = false
	}

	if spinner != nil {
		spinner.Success()
	}

	return report, nil
}

func renderReport(report *DoctorReport) error {
	term.Section("System")

	systemTabledata := pterm.TableData{
		{"Name", report.System.Name},
		{"Version", report.System.Version},
		{"ID", report.System.ID},
		{"Branding", report.System.Branding},
		{"Platform", report.System.Platform},
		{"Architecture", report.System.Architecture},
	}

	mapKeys := lo.Keys(report.System.Extras)
	slices.Sort(mapKeys)
	for _, key := range mapKeys {
		systemTabledata = append(systemTabledata, []string{key, report.System.Extras[key]})
	}

	if report.System.CPU != "" {
		systemTabledata = append(systemTabledata, []string{"CPU", report.System.CPU})
	}
	if report.System.GPU != "" {
		systemTabledata = append(systemTabledata, []string{"GPU", report.System.GPU})
	}
	if report.System.Memory != "" {
		systemTabledata = append(systemTabledata, []string{"Memory", report.System.Memory})
	}

	if err := pterm.DefaultTable.WithBoxed().WithData(systemTabledata).Render(); err != nil {
		return err
	}

	term.Section("Build Environment")

	tableData := pterm.TableData{
		{"Wails CLI", report.Build.WailsVersion},
		{"Go Version", report.Build.GoVersion},
	}

	mapKeys = lo.Keys(report.Build.Settings)
	slices.Sort(mapKeys)
	for _, key := range mapKeys {
		tableData = append(tableData, []string{key, report.Build.Settings[key]})
	}

	if err := pterm.DefaultTable.WithBoxed(true).WithData(tableData).Render(); err != nil {
		return err
	}

	term.Section("Dependencies")
	dependenciesBox := pterm.DefaultBox.WithTitleBottomCenter().WithTitle(pterm.Gray("*") + " - Optional Dependency")
	if len(report.Dependencies) == 0 {
		pterm.Info.Println("No dependencies found")
	} else {
		var dependencyTableData pterm.TableData
		var optionals pterm.TableData
		mapKeys = lo.Keys(report.Dependencies)
		for _, key := range mapKeys {
			if strings.HasPrefix(report.Dependencies[key], "*") {
				optionals = append(optionals, []string{key, report.Dependencies[key]})
			} else {
				dependencyTableData = append(dependencyTableData, []string{key, report.Dependencies[key]})
			}
		}
		dependencyTableData = append(dependencyTableData, optionals...)
		dependenciesTableString, _ := pterm.DefaultTable.WithData(dependencyTableData).Srender()
		dependenciesBox.Println(dependenciesTableString)
	}

	term.Section("Signing")
	signingTableData := pterm.TableData{}
	signingMap := formatSigningStatus(report.Signing)
	for key, value := range signingMap {
		signingTableData = append(signingTableData, []string{key, value})
	}
	if err := pterm.DefaultTable.WithBoxed(true).WithData(signingTableData).Render(); err != nil {
		return err
	}

	term.Section("Checking for issues")
	if len(report.Diagnostics) == 0 {
		pterm.Success.Println("No issues found")
	} else {
		pterm.Warning.Println("Found potential issues:")
		for _, result := range report.Diagnostics {
			pterm.Printf("• %s: %s\n", result.TestName, result.ErrorMsg)
			url := result.HelpURL
			if strings.HasPrefix(url, "/") {
				url = "https://v3.wails.io" + url
			}
			pterm.Printf("  For more information: %s\n", term.Hyperlink(url, url))
		}
	}

	term.Section("Diagnosis")
	if !report.Ready {
		term.Warning("There are some items above that need addressing!")
	} else {
		term.Success("Your system is ready for Wails development!")
	}

	return nil
}
