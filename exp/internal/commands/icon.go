package commands

import "fmt"

type IconOptions struct {
	Input string

	WindowsFilename string
	MacFilename     string
}

func Icon(options *IconOptions) error {
	if options.Input == "" {
		return fmt.Errorf("input is required")
	}

	return nil
}
