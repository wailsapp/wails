package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// ProgramHelper - Utility functions around installed applications
type ProgramHelper struct {
	shell *ShellHelper
}

// NewProgramHelper - Creates a new ProgramHelper
func NewProgramHelper() *ProgramHelper {
	return &ProgramHelper{
		shell: NewShellHelper(),
	}
}

// IsInstalled tries to determine if the given binary name is installed
func (p *ProgramHelper) IsInstalled(programName string) bool {
	_, err := exec.LookPath(programName)
	return err == nil
}

// Program - A struct to define an installed application/binary
type Program struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// FindProgram attempts to find the given program on the system.FindProgram
// Returns a struct with the name and path to the program
func (p *ProgramHelper) FindProgram(programName string) *Program {
	path, err := exec.LookPath(programName)
	if err != nil {
		return nil
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return nil
	}
	return &Program{
		Name: programName,
		Path: path,
	}
}

// GetFullPathToBinary returns the full path the the current binary
func (p *Program) GetFullPathToBinary() (string, error) {
	return filepath.Abs(p.Path)
}

// Run will execute the program with the given parameters
// Returns stdout + stderr as strings and an error if one occured
func (p *Program) Run(vars ...string) (stdout, stderr string, exitCode int, err error) {
	command, err := p.GetFullPathToBinary()
	if err != nil {
		return "", "", 1, err
	}
	cmd := exec.Command(command, vars...)
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	err = cmd.Run()
	stdout = string(stdo.Bytes())
	stderr = string(stde.Bytes())

	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			exitCode = 1
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return
}

// InstallGoPackage installs the given Go package
func (p *ProgramHelper) InstallGoPackage(packageName string) error {
	args := strings.Split("get -u "+packageName, " ")
	_, stderr, err := p.shell.Run("go", args...)
	if err != nil {
		fmt.Println(stderr)
	}
	return err
}

// RunCommand runs the given command
func (p *ProgramHelper) RunCommand(command string) error {
	args := strings.Split(command, " ")
	program := args[0]
	// TODO: Run FindProgram here and get the full path to the exe
	program, err := exec.LookPath(program)
	if err != nil {
		fmt.Printf("ERROR: Looks like '%s' isn't installed. Please install and try again.", program)
		return err
	}
	args = args[1:]
	_, stderr, err := p.shell.Run(program, args...)
	if err != nil {
		fmt.Println(stderr)
	}
	return err
}
