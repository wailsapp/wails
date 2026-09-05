package commands

import (
	"github.com/wailsapp/wails/v3/pkg/doctor-ng/tui"
)

type DoctorNgOptions struct {
	NonInteractive bool `name:"n" description:"Run in non-interactive mode (no TUI)"`
}

func DoctorNg(options *DoctorNgOptions) error {
	DisableFooter = true

	if options.NonInteractive {
		return tui.RunNonInteractive()
	}

	return tui.RunSimple()
}
