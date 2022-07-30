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

	releaseNotes := showCommand.NewSubCommand("releasenotes", "Shows the release notes for the current version")
	releaseNotes.Action(func() error {
		app.PrintBanner()
		releaseNotes := github.GetReleaseNotes(internal.Version)
		_, _ = fmt.Fprintln(w, releaseNotes)
		return nil
	})
}
