package commands

import (
	"fmt"

	"github.com/konoui/lipo/pkg/lipo"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func ToolLipo(options *flags.Lipo) error {
	DisableFooter = true

	if len(options.Inputs) < 2 {
		return fmt.Errorf("lipo requires at least 2 input files")
	}

	if options.Output == "" {
		return fmt.Errorf("output file is required (-output)")
	}

	l := lipo.New(
		lipo.WithInputs(options.Inputs...),
		lipo.WithOutput(options.Output),
	)

	if err := l.Create(); err != nil {
		return fmt.Errorf("failed to create universal binary: %w", err)
	}

	return nil
}
