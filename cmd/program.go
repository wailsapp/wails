package cmd

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"syscall"
)

// ProgramHelper - Utility functions around installed applications
type ProgramHelper struct{}

// NewProgramHelper - Creates a new ProgramHelper
func NewProgramHelper() *ProgramHelper {
	return &ProgramHelper{}
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

func (p *Program) GetFullPathToBinary() (string, error) {
	return filepath.Abs(p.Path)
}

// Run will execute the program with the given parameters
// Returns stdout + stderr as strings and an error if one occured
func (p *Program) Run(vars ...string) (stdout, stderr string, err error, exitCode int) {
	command, err := p.GetFullPathToBinary()
	if err != nil {
		return "", "", err, 1
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
