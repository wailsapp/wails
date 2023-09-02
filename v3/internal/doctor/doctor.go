package doctor

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/internal/version"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
)

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

	platformExtras := getInfo()

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

	pterm.DefaultSection.Println("Operating System")

	systemTabledata := pterm.TableData{
		{pterm.Sprint("Name"), info.Name},
		{pterm.Sprint("Version"), info.Version},
		{pterm.Sprint("ID"), info.ID},
		{pterm.Sprint("Branding"), info.Branding},

		{pterm.Sprint("Platform"), runtime.GOOS},
		{pterm.Sprint("Architecture"), runtime.GOARCH},
	}

	mapKeys = lo.Keys(platformExtras)
	slices.Sort(mapKeys)
	for _, key := range mapKeys {
		systemTabledata = append(systemTabledata, []string{key, platformExtras[key]})
	}

	err = pterm.DefaultTable.WithData(systemTabledata).Render()
	if err != nil {
		return err
	}
	/*
		pterm.DefaultSection.Println("Dependencies")

		// Output Dependencies Status
		var dependenciesMissing []string
		var externalPackages []*packagemanager.Dependency
		var dependenciesAvailableRequired = 0
		var dependenciesAvailableOptional = 0

		dependenciesTableData := pterm.TableData{
			{"Dependency", "Package Name", "Status", "Version"},
		}

		hasOptionalDependencies := false
		// Loop over dependencies
		for _, dependency := range info.Dependencies {
			name := dependency.Name

			if dependency.Optional {
				name = pterm.Gray("*") + name
				hasOptionalDependencies = true
			}

			packageName := "Unknown"
			status := pterm.LightRed("Not Found")

			// If we found the package
			if dependency.PackageName != "" {
				packageName = dependency.PackageName

				// If it's installed, update the status
				if dependency.Installed {
					status = pterm.LightGreen("Installed")
				} else {
					// Generate meaningful status text
					status = pterm.LightMagenta("Available")

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

		dependenciesTableString, _ := pterm.DefaultTable.WithHasHeader(true).WithData(dependenciesTableData).Srender()
		dependenciesBox := pterm.DefaultBox.WithTitleBottomCenter()

		if hasOptionalDependencies {
			dependenciesBox = dependenciesBox.WithTitle(pterm.Gray("*") + " - Optional Dependency")
		}

		dependenciesBox.Println(dependenciesTableString)

		pterm.DefaultSection.Println("Diagnosis")

		// Generate an appropriate diagnosis

		if dependenciesAvailableRequired != 0 {
			pterm.Println("Required package(s) installation details: \n" + info.Dependencies.InstallAllRequiredCommand())
		}

		if dependenciesAvailableOptional != 0 {
			pterm.Println("Optional package(s) installation details: \n" + info.Dependencies.InstallAllOptionalCommand())
		}

		if len(dependenciesMissing) == 0 && dependenciesAvailableRequired == 0 {
			pterm.Success.Println("Your system is ready for Wails development!")
		} else {
			pterm.Warning.Println("Your system has missing dependencies!")
		}

		if len(dependenciesMissing) != 0 {
			pterm.Println("Fatal:")
			pterm.Println("Required dependencies missing: " + strings.Join(dependenciesMissing, " "))
			pterm.Println("Please read this article on how to resolve this: https://wails.io/guides/resolving-missing-packages")
		}

		pterm.Println() // Spacer for sponsor message
	*/
	return nil
}
