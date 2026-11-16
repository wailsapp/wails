package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

type CpOptions struct{}

func Cp(_ *CpOptions) error {
	DisableFooter = true

	if len(os.Args) < 5 {
		return fmt.Errorf("cp requires a source and destination")
	}
	src, err := filepath.Abs(os.Args[3])
	if err != nil {
		return err
	}
	for _, destination := range os.Args[4:] {
		dst, err := filepath.Abs(destination)
		if err != nil {
			return err
		}
		if info, err := os.Stat(dst); err == nil && info.IsDir() {
			dst = filepath.Join(dst, filepath.Base(src))
		}
		err = copyFile(src, dst, 0)
		if err != nil {
			return fmt.Errorf("error creating %s: %s", dst, err.Error())
		}
	}
	return nil
}
