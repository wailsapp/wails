package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v3/internal/version"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type EjectOptions struct {
	Force bool `name:"force" description:"Replace an existing wails.ejected.hcl atomically"`
}

func Eject(options *EjectOptions, args []string) error {
	return ejectWithOperations(options, args, ejectOperations{
		getwd:   os.Getwd,
		write:   manifest.Eject,
		version: version.String(),
		output:  os.Stdout,
	})
}

type ejectOperations struct {
	getwd   func() (string, error)
	write   func(string, string, string, bool) error
	version string
	output  io.Writer
}

func ejectWithOperations(options *EjectOptions, args []string, operations ejectOperations) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: wails3 eject [--force]")
	}
	root, err := operations.getwd()
	if err != nil {
		return err
	}
	if err := operations.write(root, "", operations.version, options.Force); err != nil {
		return err
	}
	_, err = fmt.Fprintln(operations.output, "Wrote the complete resolved reference manifest to wails.ejected.hcl")
	return err
}
