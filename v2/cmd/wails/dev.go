package main

import (
	"os"

	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/dev"
	"github.com/wailsapp/wails/v2/internal/tui"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

func devApplication(f *flags.Dev) error {
	if f.NoColour {
		tui.SetNoColour()
	}

	quiet := f.Verbosity == flags.Quiet

	// Create logger
	logger := clilogger.New(os.Stdout)
	logger.Mute(quiet)

	if !quiet {
		app.PrintBanner()
	}

	err := f.Process()
	if err != nil {
		return err
	}

	return dev.Application(f, logger)
}
