//go:build linux

package doctor

func getInfo() (map[string]string, bool) {
	result := make(map[string]string)
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

		pterm.Println() // Spacer for sponsor message
	*/
	return result, true
}
