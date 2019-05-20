package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// NewCli - Creates a new Cli application object
func NewCli(name, description string) *Cli {
	result := &Cli{}
	result.rootCommand = NewCommand(name, description, result, "")
	result.log = NewLogger()
	return result
}

// Cli - The main application object
type Cli struct {
	rootCommand    *Command
	defaultCommand *Command
	preRunCommand  func(*Cli) error
	log            *Logger
}

// Version - Set the Application version string
func (c *Cli) Version(version string) {
	c.rootCommand.AppVersion = version
}

// PrintHelp - Prints the application's help
func (c *Cli) PrintHelp() {
	c.rootCommand.PrintHelp()
}

// Run - Runs the application with the given arguments
func (c *Cli) Run(args ...string) error {
	if c.preRunCommand != nil {
		err := c.preRunCommand(c)
		if err != nil {
			return err
		}
	}
	if len(args) == 0 {
		args = os.Args[1:]
	}
	return c.rootCommand.Run(args)
}

// DefaultCommand - Sets the given command as the command to run when
// no other commands given
func (c *Cli) DefaultCommand(defaultCommand *Command) *Cli {
	c.defaultCommand = defaultCommand
	return c
}

// Command - Adds a command to the application
func (c *Cli) Command(name, description string) *Command {
	return c.rootCommand.Command(name, description)
}

// PreRun - Calls the given function before running the specific command
func (c *Cli) PreRun(callback func(*Cli) error) {
	c.preRunCommand = callback
}

// BoolFlag - Adds a boolean flag to the root command
func (c *Cli) BoolFlag(name, description string, variable *bool) *Command {
	c.rootCommand.BoolFlag(name, description, variable)
	return c.rootCommand
}

// StringFlag - Adds a string flag to the root command
func (c *Cli) StringFlag(name, description string, variable *string) *Command {
	c.rootCommand.StringFlag(name, description, variable)
	return c.rootCommand
}

// Action represents a function that gets calls when the command is called by
// the user
type Action func() error

// Command represents a command that may be run by the user
type Command struct {
	Name              string
	CommandPath       string
	Shortdescription  string
	Longdescription   string
	AppVersion        string
	SubCommands       []*Command
	SubCommandsMap    map[string]*Command
	longestSubcommand int
	ActionCallback    Action
	App               *Cli
	Flags             *flag.FlagSet
	flagCount         int
	log               *Logger
	helpFlag          bool
	hidden            bool
}

// NewCommand creates a new Command
func NewCommand(name string, description string, app *Cli, parentCommandPath string) *Command {
	result := &Command{
		Name:             name,
		Shortdescription: description,
		SubCommandsMap:   make(map[string]*Command),
		App:              app,
		log:              NewLogger(),
		hidden:           false,
	}

	// Set up command path
	if parentCommandPath != "" {
		result.CommandPath += parentCommandPath + " "
	}
	result.CommandPath += name

	// Set up flag set
	result.Flags = flag.NewFlagSet(result.CommandPath, flag.ContinueOnError)
	result.BoolFlag("help", "Get help on the '"+result.CommandPath+"' command.", &result.helpFlag)

	// result.Flags.Usage = result.PrintHelp

	return result
}

// parseFlags parses the given flags
func (c *Command) parseFlags(args []string) error {
	// Parse flags
	tmp := os.Stderr
	os.Stderr = nil
	err := c.Flags.Parse(args)
	os.Stderr = tmp
	if err != nil {
		fmt.Printf("Error: %s\n\n", err.Error())
		c.PrintHelp()
	}
	return err
}

// Run - Runs the Command with the given arguments
func (c *Command) Run(args []string) error {

	// If we have arguments, process them
	if len(args) > 0 {
		// Check for subcommand
		subcommand := c.SubCommandsMap[args[0]]
		if subcommand != nil {
			return subcommand.Run(args[1:])
		}

		// Parse flags
		err := c.parseFlags(args)
		if err != nil {
			fmt.Printf("Error: %s\n\n", err.Error())
			c.PrintHelp()
			return err
		}

		// Help takes precedence
		if c.helpFlag {
			c.PrintHelp()
			return nil
		}
	}

	// Do we have an action?
	if c.ActionCallback != nil {
		return c.ActionCallback()
	}

	// If we haven't specified a subcommand
	// check for an app level default command
	if c.App.defaultCommand != nil {
		// Prevent recursion!
		if c.App.defaultCommand != c {
			// only run default command if no args passed
			if len(args) == 0 {
				return c.App.defaultCommand.Run(args)
			}
		}
	}

	// Nothing left we can do
	c.PrintHelp()

	return nil
}

// Action - Define an action from this command
func (c *Command) Action(callback Action) *Command {
	c.ActionCallback = callback
	return c
}

// PrintHelp - Output the help text for this command
func (c *Command) PrintHelp() {
	c.log.PrintBanner()

	commandTitle := c.CommandPath
	if c.Shortdescription != "" {
		commandTitle += " - " + c.Shortdescription
	}
	// Ignore root command
	if c.CommandPath != c.Name {
		c.log.Yellow(commandTitle)
	}
	if c.Longdescription != "" {
		fmt.Println()
		fmt.Println(c.Longdescription + "\n")
	}
	if len(c.SubCommands) > 0 {
		c.log.White("Available commands:")
		fmt.Println("")
		for _, subcommand := range c.SubCommands {
			if subcommand.isHidden() {
				continue
			}
			spacer := strings.Repeat(" ", 3+c.longestSubcommand-len(subcommand.Name))
			isDefault := ""
			if subcommand.isDefaultCommand() {
				isDefault = "[default]"
			}
			fmt.Printf("   %s%s%s %s\n", subcommand.Name, spacer, subcommand.Shortdescription, isDefault)
		}
		fmt.Println("")
	}
	if c.flagCount > 0 {
		c.log.White("Flags:")
		fmt.Println()
		c.Flags.SetOutput(os.Stdout)
		c.Flags.PrintDefaults()
		c.Flags.SetOutput(os.Stderr)

	}
	fmt.Println()
}

// isDefaultCommand returns true if called on the default command
func (c *Command) isDefaultCommand() bool {
	return c.App.defaultCommand == c
}

// isHidden returns true if the command is a hidden command
func (c *Command) isHidden() bool {
	return c.hidden
}

// Hidden hides the command from the Help system
func (c *Command) Hidden() {
	c.hidden = true
}

// Command - Defines a subcommand
func (c *Command) Command(name, description string) *Command {
	result := NewCommand(name, description, c.App, c.CommandPath)
	result.log = c.log
	c.SubCommands = append(c.SubCommands, result)
	c.SubCommandsMap[name] = result
	if len(name) > c.longestSubcommand {
		c.longestSubcommand = len(name)
	}
	return result
}

// BoolFlag - Adds a boolean flag to the command
func (c *Command) BoolFlag(name, description string, variable *bool) *Command {
	c.Flags.BoolVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// StringFlag - Adds a string flag to the command
func (c *Command) StringFlag(name, description string, variable *string) *Command {
	c.Flags.StringVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// LongDescription - Sets the long description for the command
func (c *Command) LongDescription(Longdescription string) *Command {
	c.Longdescription = Longdescription
	return c
}
