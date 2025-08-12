package shell

import (
	"bytes"
	"os"
	"os/exec"
)

type Command struct {
	command    string
	args       []string
	env        []string
	dir        string
	stdo, stde bytes.Buffer
}

func NewCommand(command string) *Command {
	return &Command{
		command: command,
		env:     os.Environ(),
	}
}

func (c *Command) Dir(dir string) {
	c.dir = dir
}

func (c *Command) Env(name string, value string) {
	c.env = append(c.env, name+"="+value)
}

func (c *Command) Run() error {
	cmd := exec.Command(c.command, c.args...)
	if c.dir != "" {
		cmd.Dir = c.dir
	}
	cmd.Stdout = &c.stdo
	cmd.Stderr = &c.stde
	return cmd.Run()
}

func (c *Command) Stdout() string {
	return c.stdo.String()
}

func (c *Command) Stderr() string {
	return c.stde.String()
}

func (c *Command) AddArgs(args []string) {
	c.args = append(c.args, args...)
}

// CreateCommand returns a *Cmd struct that when run, will run the given command + args in the given directory
func CreateCommand(directory string, command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.Dir = directory
	return cmd
}

// RunCommand will run the given command + args in the given directory
// Will return stdout, stderr and error
func RunCommand(directory string, command string, args ...string) (string, string, error) {
	return RunCommandWithEnv(nil, directory, command, args...)
}

// RunCommandWithEnv will run the given command + args in the given directory and using the specified env.
//
// Env specifies the environment of the process. Each entry is of the form "key=value".
// If Env is nil, the new process uses the current process's environment.
//
// Will return stdout, stderr and error
func RunCommandWithEnv(env []string, directory string, command string, args ...string) (string, string, error) {
	cmd := CreateCommand(directory, command, args...)
	cmd.Env = env

	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	err := cmd.Run()
	return stdo.String(), stde.String(), err
}

// RunCommandVerbose will run the given command + args in the given directory
// Will return an error if one occurs
func RunCommandVerbose(directory string, command string, args ...string) error {
	cmd := CreateCommand(directory, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

// CommandExists returns true if the given command can be found on the shell
func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
