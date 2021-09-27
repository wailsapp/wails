//go:build dev
// +build dev

package main

func init() {

	commandDescription := `This command provides access to developer tooling.`
	devCommand := app.Command("dev", "A selection of developer tools").
		LongDescription(commandDescription)

	// Add subcommands
	newTemplate(devCommand)

	devCommand.Action(func() error {
		devCommand.PrintHelp()
		return nil
	})
}
