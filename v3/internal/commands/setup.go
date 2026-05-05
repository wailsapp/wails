package commands

import (
	"github.com/wailsapp/wails/v3/internal/setupwizard"
)

type SetupOptions struct{}

func Setup(_ *SetupOptions) error {
	DisableFooter = true
	wizard := setupwizard.New()
	return wizard.Run()
}
