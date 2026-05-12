package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/github"
	"github.com/wailsapp/wails/v2/internal/tui"
)

func showReleaseNotes(f *flags.ShowReleaseNotes) error {
	if f.NoColour {
		tui.SetNoColour()
	}

	version := internal.Version
	if f.Version != "" {
		version = f.Version
	}

	app.PrintBanner()
	releaseNotes := github.GetReleaseNotes(version, f.NoColour)
	fmt.Println(releaseNotes)

	return nil
}
