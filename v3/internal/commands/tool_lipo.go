package commands

import (
	"fmt"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/universalbinary"
)

func ToolLipo(options *flags.Lipo) error {
	DisableFooter = true

	if len(options.Inputs) < 2 {
		return fmt.Errorf("lipo requires at least 2 input files")
	}

	if options.Output == "" {
		return fmt.Errorf("output file is required (-output)")
	}

	if err := universalbinary.Create(options.Output, options.Inputs); err != nil {
		return fmt.Errorf("failed to create universal binary: %w", err)
	}

	return nil
}
