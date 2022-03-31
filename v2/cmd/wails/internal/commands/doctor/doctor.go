package doctor

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/internal/system"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

// AddSubcommand adds the `doctor` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("doctor", "Diagnose your environment")

	command.Action(func() error {

		logger := clilogger.New(w)

		app.PrintBanner()

		logger.Print("Scanning system - Please wait (this may take a long time)...")

		// Get system info
		info, err := system.GetInfo()
		if err != nil {
			logger.Println("Failed.")
			return err
		}
		logger.Println("Done.")

		logger.Println("")

		// Start a new tabwriter
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)

		// Write out the system information
		fmt.Fprintf(w, "System\n")
		fmt.Fprintf(w, "------\n")
		fmt.Fprintf(w, "%s\t%s\n", "OS:", info.OS.Name)
		fmt.Fprintf(w, "%s\t%s\n", "Version: ", info.OS.Version)
		fmt.Fprintf(w, "%s\t%s\n", "ID:", info.OS.ID)

		// Output Go Information
		fmt.Fprintf(w, "%s\t%s\n", "Go Version:", runtime.Version())
		fmt.Fprintf(w, "%s\t%s\n", "Platform:", runtime.GOOS)
		fmt.Fprintf(w, "%s\t%s\n", "Architecture:", runtime.GOARCH)

		// Write out the wails information
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "Wails\n")
		fmt.Fprintf(w, "------\n")
		fmt.Fprintf(w, "%s\t%s\n", "Version: ", app.Version())

		printBuildSettings(w)

		// Exit early if PM not found
		if info.PM != nil {
			fmt.Fprintf(w, "%s\t%s\n", "Package Manager: ", info.PM.Name())
		}

		// Output Dependencies Status
		var dependenciesMissing = []string{}
		var externalPackages = []*packagemanager.Dependancy{}
		var dependenciesAvailableRequired = 0
		var dependenciesAvailableOptional = 0
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "Dependency\tPackage Name\tStatus\tVersion\n")
		fmt.Fprintf(w, "----------\t------------\t------\t-------\n")

		hasOptionalDependencies := false
		// Loop over dependencies
		for _, dependency := range info.Dependencies {

			name := dependency.Name
			if dependency.Optional {
				name = "*" + name
				hasOptionalDependencies = true
			}
			packageName := "Unknown"
			status := "Not Found"

			// If we found the package
			if dependency.PackageName != "" {

				packageName = dependency.PackageName

				// If it's installed, update the status
				if dependency.Installed {
					status = "Installed"
				} else {
					// Generate meaningful status text
					status = "Available"

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

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, packageName, status, dependency.Version)
		}
		if hasOptionalDependencies {
			fmt.Fprintf(w, "\n")
			fmt.Fprintf(w, "* - Optional Dependency\n")
		}
		w.Flush()
		logger.Println("")
		logger.Println("Diagnosis")
		logger.Println("---------")

		// Generate an appropriate diagnosis

		if len(dependenciesMissing) == 0 && dependenciesAvailableRequired == 0 {
			logger.Println("Your system is ready for Wails development!")
		} else {
			logger.Println("Your system has missing dependencies!\n")
		}

		if dependenciesAvailableRequired != 0 {
			logger.Println("Required package(s) installation details: \n" + info.Dependencies.InstallAllRequiredCommand())
		}

		if dependenciesAvailableOptional != 0 {
			logger.Println("Optional package(s) installation details: \n" + info.Dependencies.InstallAllOptionalCommand())
		}
		//
		//if len(externalPackages) > 0 {
		//	for _, p := range externalPackages {
		//		if p.Optional {
		//			print("[Optional] ")
		//		}
		//		logger.Println("Install " + p.Name + ": " + p.InstallCommand)
		//	}
		//}

		if len(dependenciesMissing) != 0 {
			// TODO: Check if apps are available locally and if so, adjust the diagnosis
			logger.Println("Fatal:")
			logger.Println("Required dependencies missing: " + strings.Join(dependenciesMissing, " "))
			logger.Println("Please read this article on how to resolve this: https://wails.app/guides/resolving-missing-packages")
		}

		logger.Println("")
		return nil
	})

	return nil
}
