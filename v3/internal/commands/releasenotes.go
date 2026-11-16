package commands

import (
	"github.com/wailsapp/wails/v3/internal/github"
	"github.com/wailsapp/wails/v3/internal/releasenotes"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/version"
)

type ReleaseNotesOptions struct {
	Version  string `name:"v" description:"The version to show release notes for"`
	NoColour bool   `name:"n" description:"Disable colour output"`
}

func ReleaseNotes(options *ReleaseNotesOptions) error {
	if options.NoColour {
		term.DisableColor()
	}

	term.Header("Release Notes")

	if version.IsDev() {
		term.Println("Release notes are not available for development builds")
		return nil
	}

	currentVersion := version.String()
	if options.Version != "" {
		currentVersion = options.Version
	}

	notes, err := github.GetReleaseNotes(currentVersion)
	if err != nil {
		term.Println(err.Error())
		return nil
	}
	term.Println(releasenotes.Render(notes, options.NoColour))
	return nil
}
