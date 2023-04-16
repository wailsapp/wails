package main

import (
	"os"

	"github.com/pterm/pterm"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v3/internal/commands"
)

func main() {
	app := clir.NewCli("wails", "The Wails CLI", "v3")
	app.NewSubCommandFunction("build", "Build the project", commands.Build)
	app.NewSubCommandFunction("init", "Initialise a new project", commands.Init)
	task := app.NewSubCommand("task", "Run and list tasks")
	var taskFlags commands.RunTaskOptions
	task.AddFlags(&taskFlags)
	task.Action(func() error {
		return commands.RunTask(&taskFlags, task.OtherArgs())
	})
	task.LongDescription("\nUsage: wails task [taskname] [flags]\n\nTasks are defined in the `Taskfile.yaml` file. See https://taskfile.dev for more information.")
	generate := app.NewSubCommand("generate", "Generation tools")
	generate.NewSubCommandFunction("defaults", "Generate default build assets", commands.Defaults)
	generate.NewSubCommandFunction("icons", "Generate icons", commands.GenerateIcons)
	generate.NewSubCommandFunction("syso", "Generate Windows .syso file", commands.GenerateSyso)
	generate.NewSubCommandFunction("bindings", "Generate bindings + models", commands.GenerateBindings)
	plugin := app.NewSubCommand("plugin", "Plugin tools")
	//plugin.NewSubCommandFunction("list", "List plugins", commands.PluginList)
	plugin.NewSubCommandFunction("init", "Initialise a new plugin", commands.PluginInit)
	//plugin.NewSubCommandFunction("add", "Add a plugin", commands.PluginAdd)

	err := app.Run()
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
