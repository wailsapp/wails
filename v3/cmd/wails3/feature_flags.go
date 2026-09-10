package main

import (
	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v3/internal/commands"
	"github.com/wailsapp/wails/v3/internal/features"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func addBuildFlags(command *clir.Command, options *flags.Build) {
	if features.WakeEnabled() {
		command.AddFlags(options)
		return
	}
	command.AddFlags(&options.Common)
	command.StringFlag("tags", "Additional build tags to pass to the Go compiler (comma-separated)", &options.Tags)
	command.BoolFlag("obfuscated", "Build with garble and stable obfuscated binding IDs", &options.Obfuscated)
	command.StringFlag("garbleargs", "Additional arguments to pass to garble before the build command", &options.GarbleArgs)
}

func addDevFlags(command *clir.Command, options *commands.DevOptions) {
	if features.WakeEnabled() {
		command.AddFlags(options)
		return
	}
	command.StringsFlag("appargs", "Application arguments (requires an active HCL project)", &options.AppArgs)
	options.Config = "./build/config.yml"
	command.AddFlags(&options.Common)
	command.StringFlag("config", "The config file including path", &options.Config)
	command.IntFlag("port", "Specify the vite dev server port", &options.VitePort)
	command.BoolFlag("s", "Enable HTTPS", &options.Secure)
	command.StringFlag("tags", "Additional development build tags to pass to the Go compiler (comma-separated)", &options.Tags)
}
