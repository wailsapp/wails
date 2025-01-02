package commands

import (
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/github"
	"github.com/wailsapp/wails/v3/internal/version"
)

type ReleaseNotesOptions struct {
	Version  string `name:"v" description:"The version to show release notes for"`
	NoColour bool   `name:"n" description:"Disable colour output"`
}

func ReleaseNotes(options *ReleaseNotesOptions) error {
	if options.NoColour {
		pterm.DisableColor()
	}

	currentVersion := version.VersionString
	if options.Version != "" {
		currentVersion = options.Version
	}

	releaseNotes := github.GetReleaseNotes(currentVersion, options.NoColour)
	pterm.Println(releaseNotes)
	return nil
}
