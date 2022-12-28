package main

import (
	"log"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/exp/internal/commands"
)

func main() {

	app := clir.NewCli("wails", "The Wails CLI", "v3")
	app.NewSubCommandFunction("init", "Initialise a new project", commands.Init)
	app.NewSubCommandFunction("build", "Build the project", commands.Build)
	tool := app.NewSubCommand("tool", "Various build tools")
	tool.NewSubCommandFunction("icon", "Generate icons", commands.Icon)

	err := app.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
