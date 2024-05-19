package commands

import (
	"fmt"
	"github.com/wailsapp/wails/v3/internal/commands/services"
	"github.com/wailsapp/wails/v3/internal/flags"
	"os"
)

func ServiceCreate(options *flags.ServiceCreateOptions) error {
	if options.OutputDir == "./services" {
		// Check the existence of the "services" directory
		_, servicesErr := os.Stat("services")
		_, goModError := os.Stat("go.mod")
		if servicesErr != nil || goModError != nil {
			return fmt.Errorf("this command should be run from the project root")
		}
	}
	return services.Install(options)
}
