package main

import (
	"fmt"

	"github.com/wailsapp/wails/cmd"
)

func init() {

	commandDescription := `Outputs the current version of the wails cli tool.`

	versionCommand := app.Command("version", "Print Wails cli version").
		LongDescription(commandDescription)

	versionCommand.Action(func() error {
		fmt.Println(cmd.Version)
		return nil
	})
}