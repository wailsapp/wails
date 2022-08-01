package show

import (
	"fmt"
	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/github"
	"io"
)

// AddSubcommand adds the `show` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) {
	showCommand := app.NewSubCommand("show", "Shows various information")

	version := internal.Version
	releaseNotes := showCommand.NewSubCommand("releasenotes", "Shows the release notes for the current version")
	releaseNotes.StringFlag("version", "The version to show the release notes for", &version)
	releaseNotes.Action(func() error {
		app.PrintBanner()
		releaseNotes := github.GetReleaseNotes(version)
		_, _ = fmt.Fprintln(w, releaseNotes)
		return nil
	})
}
