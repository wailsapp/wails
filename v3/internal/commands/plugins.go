package commands

import (
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/plugins"
	"strings"

	"github.com/pterm/pterm"
)

func toCamelCasePlugin(s string) string {
	var camelCase string
	var capitalize = true

	for _, c := range s {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
			if capitalize {
				camelCase += strings.ToUpper(string(c))
				capitalize = false
			} else {
				camelCase += string(c)
			}
		} else if c >= 'A' && c <= 'Z' {
			camelCase += string(c)
			capitalize = false
		} else {
			capitalize = true
		}
	}

	return camelCase + "Plugin"
}

func PluginInit(options *flags.PluginInit) error {

	if options.Quiet {
		pterm.DisableOutput()
	}

	if options.PackageName == "" {
		options.PackageName = toCamelCasePlugin(options.Name)
	}

	return plugins.Install(options)
}
