package main

import (
	"log"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/exp/internal/commands"
)

func main() {
	app := clir.NewCli("wails", "The Wails CLI", "v3")
	app.NewSubCommandFunction("init", "Initialise a new project", commands.Init)
	app.NewSubCommandFunction("run", "Run a task", commands.Run)
	generate := app.NewSubCommand("generate", "Generation tools")
	generate.NewSubCommandFunction("defaults", "Generate default build assets", commands.Defaults)
	generate.NewSubCommandFunction("icons", "Generate icons", commands.GenerateIcons)
	generate.NewSubCommandFunction("syso", "Generate Windows .syso file", commands.GenerateSyso)
	err := app.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
