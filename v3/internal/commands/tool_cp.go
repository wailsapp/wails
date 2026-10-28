package commands

import (
	"fmt"
	"io"
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
		sourceFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer sourceFile.Close()
		sourceInfo, err := sourceFile.Stat()
		if err != nil {
			return err
		}
		input, err := io.ReadAll(sourceFile)
		if err != nil {
			return err
		}

		err = os.WriteFile(dst, input, sourceInfo.Mode().Perm())
		if err != nil {
			return fmt.Errorf("error creating %s: %s", dst, err.Error())
		}
	}
	return nil
}
