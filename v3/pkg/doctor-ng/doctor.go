package doctorng

import (
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/internal/version"
)

type Doctor struct {
	report *Report
}

func New() *Doctor {
	return &Doctor{
		report: NewReport(),
	}
}

func (d *Doctor) Run() (*Report, error) {
	if err := d.collectSystemInfo(); err != nil {
		return nil, err
	}

	if err := d.collectBuildInfo(); err != nil {
		return nil, err
	}

	if err := d.collectDependencies(); err != nil {
		return nil, err
	}

	d.runDiagnostics()
	d.generateSummary()

	return d.report, nil
}

func (d *Doctor) collectSystemInfo() error {
	info, err := operatingsystem.Info()
	if err != nil {
		return err
	}

	d.report.System.OS = OSInfo{
		Name:     info.Name,
		Version:  info.Version,
		ID:       info.ID,
		Branding: info.Branding,
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	d.report.System.Hardware = collectHardwareInfo()
	d.report.System.PlatformExtras = collectPlatformExtras()

	return nil
}

func (d *Doctor) collectBuildInfo() error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		d.report.Build = BuildInfo{
			WailsVersion: version.String(),
			GoVersion:    runtime.Version(),
		}
		return nil
	}

	settings := make(map[string]string)
	for _, s := range buildInfo.Settings {
		settings[s.Key] = s.Value
	}

	wailsVersion := strings.TrimSpace(version.String())
	wailsPackage, found := lo.Find(buildInfo.Deps, func(dep *debug.Module) bool {
		return dep.Path == "github.com/wailsapp/wails/v3"
	})

	if found && wailsPackage != nil && wailsPackage.Replace != nil {
		wailsVersion = "(local) => " + filepath.ToSlash(wailsPackage.Replace.Path)
		repo, err := git.PlainOpen(filepath.Join(wailsPackage.Replace.Path, ".."))
		if err == nil {
			head, err := repo.Head()
			if err == nil {
				wailsVersion += " (" + head.Hash().String()[:8] + ")"
			}
		}
	}

	d.report.Build = BuildInfo{
		WailsVersion: wailsVersion,
		GoVersion:    runtime.Version(),
		BuildMode:    settings["-buildmode"],
		Compiler:     settings["-compiler"],
		CGOEnabled:   settings["CGO_ENABLED"] == "1",
		Settings:     settings,
	}

	return nil
}

func (d *Doctor) generateSummary() {
	missing := d.report.Dependencies.RequiredMissing()
	errCount := 0
	warnCount := 0

	for _, diag := range d.report.Diagnostics {
		switch diag.Severity {
		case SeverityError:
			errCount++
		case SeverityWarning:
			warnCount++
		}
	}

	if len(missing) == 0 && errCount == 0 {
		d.report.Ready = true
		if warnCount > 0 {
			d.report.Summary = "System is ready for Wails development with some warnings"
		} else {
			d.report.Summary = "System is ready for Wails development!"
		}
	} else {
		d.report.Ready = false
		var parts []string
		if len(missing) > 0 {
			parts = append(parts, lo.Ternary(len(missing) == 1,
				"1 missing dependency",
				string(rune(len(missing)+'0'))+" missing dependencies"))
		}
		if errCount > 0 {
			parts = append(parts, lo.Ternary(errCount == 1,
				"1 error",
				string(rune(errCount+'0'))+" errors"))
		}
		d.report.Summary = "System has issues: " + strings.Join(parts, ", ")
	}
}
