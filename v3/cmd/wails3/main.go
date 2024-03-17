package main

import (
	"github.com/pterm/pterm"
	"github.com/samber/lo"
	"os"
	"runtime/debug"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v3/internal/commands"
)

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	commands.BuildSettings = lo.Associate(buildInfo.Settings, func(setting debug.BuildSetting) (string, string) {
		return setting.Key, setting.Value
	})
	// Iterate over the Deps and add them to the build settings using a prefix of "mod."
	for _, dep := range buildInfo.Deps {
		commands.BuildSettings["mod."+dep.Path] = dep.Version
	}
}

func main() {
	app := clir.NewCli("wails", "The Wails3 CLI", "v3")
	app.NewSubCommandFunction("init", "Initialise a new project", commands.Init)
	app.NewSubCommandFunction("build", "Build the project", commands.Build)
	app.NewSubCommandFunction("dev", "Run in Dev mode", commands.Dev)
	app.NewSubCommandFunction("package", "Package application", commands.Package)
	app.NewSubCommandFunction("doctor", "System status report", commands.Doctor)
	task := app.NewSubCommand("task", "Run and list tasks")
	var taskFlags commands.RunTaskOptions
	task.AddFlags(&taskFlags)
	task.Action(func() error {
		return commands.RunTask(&taskFlags, task.OtherArgs())
	})
	task.LongDescription("\nUsage: wails3 task [taskname] [flags]\n\nTasks are defined in the `Taskfile.yaml` file. See https://taskfile.dev for more information.")
	generate := app.NewSubCommand("generate", "Generation tools")
	generate.NewSubCommandFunction("build-assets", "Generate build assets", commands.GenerateBuildAssets)
	generate.NewSubCommandFunction("icons", "Generate icons", commands.GenerateIcons)
	generate.NewSubCommandFunction("syso", "Generate Windows .syso file", commands.GenerateSyso)
	generate.NewSubCommandFunction("bindings", "Generate bindings + models", commands.GenerateBindings)
	generate.NewSubCommandFunction("constants", "Generate JS constants from Go", commands.GenerateConstants)
	generate.NewSubCommandFunction(".desktop", "Generate .desktop file", commands.GenerateDotDesktop)
	generate.NewSubCommandFunction("appimage", "Generate Linux AppImage", commands.GenerateAppImage)

	plugin := app.NewSubCommand("plugin", "Plugin tools")
	//plugin.NewSubCommandFunction("list", "List plugins", commands.PluginList)
	plugin.NewSubCommandFunction("init", "Initialise a new plugin", commands.PluginInit)
	//plugin.NewSubCommandFunction("add", "Add a plugin", commands.PluginAdd)
	tool := app.NewSubCommand("tool", "Various tools")
	tool.NewSubCommandFunction("checkport", "Checks if a port is open. Useful for testing if vite is running.", commands.ToolCheckPort)
	tool.NewSubCommandFunction("watcher", "Watches files and runs a command when they change", commands.Watcher)
	tool.NewSubCommandFunction("cp", "Copy files", commands.Cp)

	app.NewSubCommandFunction("version", "Print the version", commands.Version)

	err := app.Run()
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
