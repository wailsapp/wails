package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

type CpOptions struct{}

func Cp(_ *CpOptions) error {
	DisableFooter = true

	// extract the source and destination from os.Args
	if len(os.Args) != 5 {
		return fmt.Errorf("cp requires a source and destination")
	}
	// Extract source
	source := os.Args[3]
	for _, destination := range os.Args[4:] {
		src, err := filepath.Abs(source)
		if err != nil {
			return err
		}
		dst, err := filepath.Abs(destination)
		if err != nil {
			return err
		}
		input, err := os.ReadFile(src)
		if err != nil {
			return err
		}

		err = os.WriteFile(dst, input, 0644)
		if err != nil {
			return fmt.Errorf("error creating %s: %s", dst, err.Error())
		}
	}
	return nil
}
