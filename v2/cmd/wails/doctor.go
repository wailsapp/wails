package main

import (
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/pterm/pterm"

	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/internal/colour"
	"github.com/wailsapp/wails/v2/internal/system"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

func diagnoseEnvironment(f *flags.Doctor) error {
	if f.NoColour {
		pterm.DisableColor()
		colour.ColourEnabled = false
	}

	pterm.DefaultSection = *pterm.DefaultSection.
		WithBottomPadding(0).
		WithStyle(pterm.NewStyle(pterm.FgBlue, pterm.Bold))

	pterm.Println() // Spacer
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println("Wails Doctor")
	pterm.Println() // Spacer

	spinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone().Start("Scanning system - Please wait (this may take a long time)...")

	// Get system info
	info, err := system.GetInfo()
	if err != nil {
		spinner.Fail()
		pterm.Error.Println("Failed to get system information")
		return err
	}
	spinner.Success()

	pterm.DefaultSection.Println("Wails")

	wailsTableData := pterm.TableData{
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

	// Exit early if PM not found
	if info.PM != nil {
		wailsTableData = append(wailsTableData, []string{"Package Manager", info.PM.Name()})
	}

	err = pterm.DefaultTable.WithData(wailsTableData).Render()
	if err != nil {
		return err
	}

	pterm.DefaultSection.Println("System")

	systemTabledata := pterm.TableData{
		{pterm.Bold.Sprint("OS"), info.OS.Name},
		{pterm.Bold.Sprint("Version"), info.OS.Version},
		{pterm.Bold.Sprint("ID"), info.OS.ID},
		{pterm.Bold.Sprint("Go Version"), runtime.Version()},
		{pterm.Bold.Sprint("Platform"), runtime.GOOS},
		{pterm.Bold.Sprint("Architecture"), runtime.GOARCH},
	}

	err = pterm.DefaultTable.WithBoxed().WithData(systemTabledata).Render()
	if err != nil {
		return err
	}

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
	return nil
}
