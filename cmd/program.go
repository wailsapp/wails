package cmd

import (
	"os/exec"
	"path/filepath"
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
