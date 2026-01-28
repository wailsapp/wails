package commands

import (
	"github.com/wailsapp/wails/v3/internal/setupwizard"
	"github.com/wailsapp/wails/v3/internal/term"
)

type SetupOptions struct{}

func Setup(_ *SetupOptions) error {
	DisableFooter = true

	term.Println("")
	term.Warning("The setup wizard is experimental and not fully tested on all platforms.")
	term.Println("If you encounter issues, please report them at:")
	term.Println("  " + term.Hyperlink("https://github.com/wailsapp/wails/issues/4904", "GitHub Issue #4904"))
	term.Println("")

	wizard := setupwizard.New()
	return wizard.Run()
}
