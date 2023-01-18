package main

import (
	"os"

	"github.com/pterm/pterm"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v3/internal/commands"
)

func main() {
	app := clir.NewCli("wails", "The Wails CLI", "v3")
	app.NewSubCommandFunction("init", "Initialise a new project", commands.Init)
	task := app.NewSubCommand("task", "Run and list tasks")
	task.NewSubCommandFunction("run", "Run a task", commands.RunTask)
	task.NewSubCommandFunction("list", "List tasks", commands.ListTasks)
	generate := app.NewSubCommand("generate", "Generation tools")
	generate.NewSubCommandFunction("defaults", "Generate default build assets", commands.Defaults)
	generate.NewSubCommandFunction("icons", "Generate icons", commands.GenerateIcons)
	generate.NewSubCommandFunction("syso", "Generate Windows .syso file", commands.GenerateSyso)
	generate.NewSubCommandFunction("bindings", "Generate bindings + models", commands.GenerateBindings)
	err := app.Run()
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
