package commands

import (
	"github.com/wailsapp/wails/v3/internal/doctor"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func Doctor(options *flags.Doctor) error {
	return doctor.Run(options.JSON)
}
